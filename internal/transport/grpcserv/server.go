package grpcserv

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	moviepb "github.com/ButyugovMxx/GoPrac/api/proto"
	"google.golang.org/grpc"
)

type Gserver struct {
	server *grpc.Server
	addr   string
}

func NewGserver(h moviepb.MovieServiceServer) *Gserver {
	s := grpc.NewServer()

	moviepb.RegisterMovieServiceServer(s, h)

	return &Gserver{
		server: s,
		addr:   ":50051",
	}
}

func (g *Gserver) Run(ctx context.Context) error {
	listener, err := net.Listen("tcp", g.addr)
	if err != nil {
		return fmt.Errorf("grpc listen: %w", err)
	}
	defer listener.Close()

	ErrCh := make(chan error)

	go func() {
		log.Println("grpc-server started")
		ErrCh <- g.server.Serve(listener)
		close(ErrCh)
	}()

	select {
	case err := <-ErrCh:
		if errors.Is(err, grpc.ErrServerStopped) {
			return nil
		}
		return fmt.Errorf("grpc serv err: %w", err)

	case <-ctx.Done():
		log.Println("grpc shutdown")
	}

	helpCh := make(chan struct{})

	go func() {
		g.server.GracefulStop()
		helpCh <- struct{}{}
	}()

	select {
	case <-helpCh:
		log.Println("grpc is good stopped")
		return nil
	case <-time.After(time.Second * 10):
		g.server.Stop()
		return fmt.Errorf("grpc sutdown timeout")
	}
}
