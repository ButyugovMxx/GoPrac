package httpserv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HttpServer struct {
	server *http.Server
}

func NewServer(h *Handler) *HttpServer {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /load-movies", h.LoadFilms)
	mux.HandleFunc("GET /load-last", h.LoadLastFilm)

	s := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	return &HttpServer{
		server: s,
	}
}

func (s *HttpServer) Run(ctx context.Context) error {
	errCh := make(chan error, 1)

	go func() {
		log.Println("HTTP server starting")
		errCh <- s.server.ListenAndServe()
	}()

	select {
	case err := <-errCh:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return fmt.Errorf("cannot start server: %w", err)

	case <-ctx.Done():
		log.Println("graceful shutdown server")
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
