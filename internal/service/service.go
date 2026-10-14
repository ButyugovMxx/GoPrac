package service

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/ButyugovMxx/GoPrac/internal/client"
)

type MovieService struct {
	client  client.MovieClient
	workers int
}

func NewMovieService(client client.MovieClient, workers int) *MovieService {
	return &MovieService{
		client:  client,
		workers: workers,
	}
}

type MovieError struct {
	Id  int
	Err error
}

type Loadres struct {
	Movies []client.Movieload
	Errors []MovieError
}

var ErrInvalidRange = errors.New("invalid movie range")
var ErrNotValueArgs = errors.New("uncorrect args: from > to")

func (m *MovieService) loadMovies(ctx context.Context, from int, to int) (Loadres, error) {
	if from == -1 || to == -1 {
		return Loadres{}, ErrInvalidRange
	}

	if to < from {
		return Loadres{}, ErrNotValueArgs
	}

	ctxT, cencel := context.WithTimeout(ctx, 5*time.Second)
	defer cencel()

	ch := make(chan int)
	defer close(ch)

	movies := make([]client.Movieload, 0)
	movieErrors := make([]MovieError, 0)
	var wg sync.WaitGroup
	var mu sync.Mutex
	for i := 0; i < m.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for id := range ch {
				movie, err := m.client.GetMovie(ctxT, id)
				mu.Lock()
				if err != nil {
					movieErrors = append(movieErrors, MovieError{
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

	wg.Wait()

	res := Loadres{
		Movies: movies,
		Errors: movieErrors,
	}
	return res, nil
}

func (m *MovieService) loadLastMovie(ctx context.Context) (client.Movieload, error) {
	var movie client.Movieload

	ctxT, cencel := context.WithTimeout(ctx, 5*time.Second)
	defer cencel()

	right := 99999
	left := 1
	for left <= right {
		middle := (right-left)/2 + left
		curMovie, err := m.client.GetMovie(ctxT, middle)
		if err == nil {
			left = middle + 1
		} else if errors.Is(err, client.ErrBadReq) {
			right = middle - 1
		} else {
			return client.Movieload{}, err
		}

		movie = curMovie
	}
	return movie, nil
}
