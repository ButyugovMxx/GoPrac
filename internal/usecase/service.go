package usecase

import (
	"context"
	"errors"
	"sync"

	"github.com/ButyugovMxx/GoPrac/internal/domain"
)

type MovieService struct {
	client  MovieProvider
	workers int
}

func NewMovieService(client MovieProvider, workers int) *MovieService {
	return &MovieService{
		client:  client,
		workers: workers,
	}
}

func (m *MovieService) LoadMovies(ctx context.Context, from int, to int) (domain.Loadres, error) {

	if to < from {
		return domain.Loadres{}, domain.ErrInvalidRange
	}

	if from < 1 {
		return domain.Loadres{}, domain.ErrInvalidRange
	}

	ch := make(chan int)

	movies := make([]domain.Movieload, 0)
	movieErrors := make([]domain.MovieError, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < m.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range ch {
				movie, err := m.client.GetMovie(ctx, id)
				mu.Lock()
				if err != nil {
					movieErrors = append(movieErrors, domain.MovieError{
						Id:  id,
						Err: err,
					})
				} else {
					movies = append(movies, movie)
				}
				mu.Unlock()
			}

		}()
	}

	for i := from; i < to+1; i++ {
		ch <- i
	}

	close(ch)

	wg.Wait()

	res := domain.Loadres{
		Movies: movies,
		Errors: movieErrors,
	}
	return res, nil
}

func (m *MovieService) LoadLastMovie(ctx context.Context) (domain.Movieload, error) {
	var movie domain.Movieload
	found := false
	right := 99999
	left := 1
	for left <= right {
		middle := (right-left)/2 + left
		curMovie, err := m.client.GetMovie(ctx, middle)
		if err == nil {
			movie = curMovie
			found = true
			left = middle + 1
		} else if errors.Is(err, domain.NotFound) {
			right = middle - 1
		} else {
			return domain.Movieload{}, err
		}
	}

	if !found {
		return domain.Movieload{}, domain.NotFound
	}
	return movie, nil
}
