package service

import (
	"context"
	"errors"
	"sync"

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

func (m *MovieService) LoadMovies(ctx context.Context, from int, to int) (Loadres, error) {
	if from == -1 || to == -1 {
		return Loadres{}, ErrInvalidRange
	}

	if to < from {
		return Loadres{}, ErrNotValueArgs
	}

	ch := make(chan int)

	movies := make([]client.Movieload, 0)
	movieErrors := make([]MovieError, 0)
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

	close(ch)

	wg.Wait()

	res := Loadres{
		Movies: movies,
		Errors: movieErrors,
	}
	return res, nil
}

func (m *MovieService) LoadLastMovie(ctx context.Context) (client.Movieload, error) {
	var movie client.Movieload
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
		} else if errors.Is(err, client.NotFound) {
			right = middle - 1
		} else {
			return client.Movieload{}, err
		}
	}

	if !found {
		return client.Movieload{}, client.NotFound
	}
	return movie, nil
}
