package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/ButyugovMxx/GoPrac/internal/client"
	"github.com/ButyugovMxx/GoPrac/internal/service"
	"github.com/ButyugovMxx/GoPrac/internal/transport/httpserv"
)

func main() {
	c := client.MovieClient{}
	s := service.NewMovieService(c, 10)
	h := httpserv.NewHandler(s)

	httpServer := httpserv.NewServer(h)
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	httpServer.Run(ctx)
}
