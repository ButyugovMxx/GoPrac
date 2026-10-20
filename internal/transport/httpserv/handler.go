package httpserv

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/ButyugovMxx/GoPrac/internal/client"
	"github.com/ButyugovMxx/GoPrac/internal/service"
)

type Handler struct {
	service *service.MovieService
}

type partialMovieError struct {
	ID    int    `json:"id"`
	Error string `json:"error"`
}

func NewHandler(s *service.MovieService) *Handler {
	return &Handler{
		service: s,
	}
}

type ToFrom struct {
	to   int
	from int
}

// LoadFilms godoc
// @Summary Load movies by ID range
// @Description Loads movies from the external movie service.
// @Tags movies
// @Produce json
// @Param from query int true "First movie ID"
// @Param to query int true "Last movie ID"
// @Success 200 {array} client.Movieload
// @Header 200 {string} X-Partial-Errors "JSON array of per-movie errors, present when some movies could not be loaded"
// @Failure 400 {string} string "Invalid range"
// @Failure 502 {string} string "Upstream error"
// @Router /load-movies [get]
func (h *Handler) LoadFilms(w http.ResponseWriter, r *http.Request) {
	m := r.URL.Query()
	from, err := strconv.Atoi(m.Get("from"))
	if err != nil {
		http.Error(w, "invalid from", http.StatusBadRequest)
		return
	}

	to, err := strconv.Atoi(m.Get("to"))
	if err != nil {
		http.Error(w, "invalid to", http.StatusBadRequest)
		return
	}

	ctx := r.Context()

	movies, err := h.service.LoadMovies(ctx, from, to)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrInvalidRange):
			http.Error(w, "Service: invalid movie range", http.StatusBadRequest)
			return
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, "Request timeout", http.StatusGatewayTimeout)

		default:
			http.Error(w, "Movie stand is unavailable", http.StatusBadGateway)
		}

		return
	}

	if len(movies.Errors) > 0 {
		partialErrors := make([]partialMovieError, 0, len(movies.Errors))
		for _, movieErr := range movies.Errors {
			partialErrors = append(partialErrors, partialMovieError{
				ID:    movieErr.Id,
				Error: movieErr.Err.Error(),
			})
		}

		errorJSON, err := json.Marshal(partialErrors)
		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		w.Header().Set("X-Partial-Errors", string(errorJSON))
	}

	w.Header().Set("Content-Type", "application/json")
	erro := json.NewEncoder(w).Encode(movies.Movies)
	if erro != nil {
		log.Printf("failed to encode movies: %v", err)
	}
}

// LoadLastFilm godoc
// @Summary Get last available movie
// @Description Finds the movie with the highest available ID.
// @Tags movies
// @Produce json
// @Success 200 {object} client.Movieload
// @Failure 404 {string} string "Movie not found"
// @Failure 502 {string} string "Upstream error"
// @Router /last-movie [get]
func (h *Handler) LoadLastFilm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	film, err := h.service.LoadLastMovie(ctx)
	if err != nil {
		switch {
		case errors.Is(err, client.NotFound):
			http.Error(w, "Movie not found", http.StatusNotFound)

		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, "Request timeout", http.StatusGatewayTimeout)

		default:
			http.Error(w, "Movie stand is unavailable", http.StatusBadGateway)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	erro := json.NewEncoder(w).Encode(film)
	if erro != nil {
		log.Printf("failed to encode movies: %v", err)
	}
}
