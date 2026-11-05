package grpcserv

import (
	"context"
	"errors"

	moviepb "github.com/ButyugovMxx/GoPrac/api/proto"
	"github.com/ButyugovMxx/GoPrac/internal/client"
	"github.com/ButyugovMxx/GoPrac/internal/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GHandler struct {
	service *service.MovieService

	moviepb.UnimplementedMovieServiceServer
}

func NewGhandler(s *service.MovieService) *GHandler {
	return &GHandler{
		service: s,
	}
}

func toProtoMovie(m client.Movieload) *moviepb.Movie {
	return &moviepb.Movie{
		Id:       m.Id,
		Title:    m.Title,
		Year:     m.Year,
		Director: m.Director,
	}
}

func (h *GHandler) LoadMovies(ctx context.Context, req *moviepb.LoadMoviesRequest) (*moviepb.LoadMoviesResponse, error) {

	result, err := h.service.LoadMovies(ctx, int(req.GetFrom()), int(req.GetTo()))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRange):
			return nil, status.Error(codes.InvalidArgument, err.Error())

		case errors.Is(err, service.ErrInvalidRange):
			return nil, status.Error(codes.InvalidArgument, err.Error())

		case errors.Is(err, context.Canceled):
			return nil, status.Error(codes.Canceled, err.Error())

		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())

		default:
			return nil, status.Error(codes.Unavailable, err.Error())
		}
	}

	resp := &moviepb.LoadMoviesResponse{
		Movies:  make([]*moviepb.Movie, 0),
		MErrors: make([]*moviepb.MovieError, 0),
	}

	for _, movie := range result.Movies {
		resp.Movies = append(resp.Movies, toProtoMovie(movie))
	}

	for _, err := range result.Errors {
		resp.MErrors = append(resp.MErrors, &moviepb.MovieError{
			Id:    int32(err.Id),
			Error: err.Err.Error(),
		})
	}

	return resp, nil
}

func (h *GHandler) GetLastMovie(ctx context.Context, req *moviepb.GetLastMovieRequest) (*moviepb.Movie, error) {

	movie, err := h.service.LoadLastMovie(ctx)
	if err != nil {
		switch {
		case errors.Is(err, client.NotFound):
			return nil, status.Error(codes.NotFound, err.Error())

		case errors.Is(err, context.DeadlineExceeded):
			return nil, status.Error(codes.DeadlineExceeded, err.Error())

		default:
			return nil, status.Error(codes.Unavailable, err.Error())
		}
	}

	return toProtoMovie(movie), nil
}
