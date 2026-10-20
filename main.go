package main

import (
	"context"
	"fmt"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ButyugovMxx/GoPrac/internal/client"
	"github.com/ButyugovMxx/GoPrac/internal/config"
	"github.com/ButyugovMxx/GoPrac/internal/service"
	"github.com/ButyugovMxx/GoPrac/internal/transport/grpcserv"
	"github.com/ButyugovMxx/GoPrac/internal/transport/httpserv"
)

func main() {
	cfg := config.Load()

	c := client.NewMovieClient(cfg.MovieBaseUrl)
	s := service.NewMovieService(c, 10)
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
