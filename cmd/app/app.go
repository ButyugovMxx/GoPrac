package app

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	movieclient "github.com/ButyugovMxx/GoPrac/internal/adapter/movieclient"
	"github.com/ButyugovMxx/GoPrac/internal/config"
	"github.com/ButyugovMxx/GoPrac/internal/handler/grpcserv"
	"github.com/ButyugovMxx/GoPrac/internal/handler/httpserv"
	"github.com/ButyugovMxx/GoPrac/internal/usecase"

	_ "github.com/ButyugovMxx/GoPrac/docs"
)

func Start() {
	cfg := config.Load()

	c := movieclient.NewMovieClient(cfg.MovieBaseUrl)
	s := usecase.NewMovieService(&c, 10)
	h := httpserv.NewHandler(s)

	gh := grpcserv.NewGhandler(s)
	grpcServer := grpcserv.NewGserver(gh)

	httpServer := httpserv.NewServer(h, cfg.Addr)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := httpServer.Run(ctx)
		if err != nil {
			fmt.Printf("cannot start http-server: %v", err)
		}
	}()
	wg.Add(1)
	go func() {
		defer wg.Done()
		erro := grpcServer.Run(ctx)
		if erro != nil {
			fmt.Printf("cannot start grpc-server: %v", erro)
		}
	}()

	wg.Wait()
}
