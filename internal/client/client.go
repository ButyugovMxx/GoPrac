package client

import (
	"context"
	"errors"
	"net/http"
)

type Movieload struct {
	Id       int32  `json:"id"`
	Title    string `json:"title"`
	Year     int32  `json:"year"`
	Director string `json:"director"`
}

type MovieClient struct {
	Client http.Client
}

var ErrBadReq = errors.New("404")

func (c *MovieClient) GetMovie(ctx context.Context, id int) (Movieload, error) {
	return Movieload{}, errors.New("GetMovie not implemented")
}
