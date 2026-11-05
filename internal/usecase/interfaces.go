package usecase

import (
	"context"

	"github.com/ButyugovMxx/GoPrac/internal/domain"
)

type MovieProvider interface {
	GetMovie(ctx context.Context, id int) (domain.Movieload, error)
}

type MovieUseCase interface {
	LoadLastMovie(ctx context.Context) (domain.Movieload, error)
	LoadMovies(ctx context.Context, from int, to int) (domain.Loadres, error)
}
