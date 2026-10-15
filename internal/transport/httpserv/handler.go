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

func NewHandler(s *service.MovieService) *Handler {
	return &Handler{
		service: s,
	}
}

type ToFrom struct {
	to   int
	from int
}

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
		case errors.Is(err, service.ErrNotValueArgs):
			http.Error(w, "Servise: args is not currect", http.StatusBadRequest)
			return
		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, "Request timeout", http.StatusGatewayTimeout)

		default:
			http.Error(w, "Internal srver error", http.StatusInternalServerError)
		}

		return
	}

	w.Header().Set("Content-Type", "application/json")
	erro := json.NewEncoder(w).Encode(movies)
	if erro != nil {
		log.Printf("failed to encode movies: %v", err)
	}
}

func (h *Handler) LoadLastFilm(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	film, err := h.service.LoadLastMovie(ctx)
	if err != nil {
		switch {
		case errors.Is(err, client.ErrBadReq):
			http.Error(w, "Movie not found", http.StatusNotFound)

		case errors.Is(err, context.DeadlineExceeded):
			http.Error(w, "Request timeout", http.StatusGatewayTimeout)

		default:
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	erro := json.NewEncoder(w).Encode(film)
	if erro != nil {
		log.Printf("failed to encode movies: %v", err)
	}
}
