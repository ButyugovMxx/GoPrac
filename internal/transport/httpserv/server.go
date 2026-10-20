package httpserv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

type HttpServer struct {
	server *http.Server
}

func NewServer(h *Handler, addr string) *HttpServer {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /load-movies", h.LoadFilms)
	mux.HandleFunc("GET /last-movie", h.LoadLastFilm)
	mux.Handle("GET /swagger/", httpSwagger.Handler())

	s := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	return &HttpServer{
		server: s,
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		log.Printf("HTTP server starting on %s", s.server.Addr)
		errCh <- s.server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("cannot start server: %w", err)

	case <-ctx.Done():
		log.Println("shutdown http server")
	}

	ctxShutdown, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	if err := s.server.Shutdown(ctxShutdown); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}

	log.Println("HTTP server stopped")
	return nil
}
