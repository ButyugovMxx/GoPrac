package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type Movieload struct {
	Id       int32  `json:"id"`
	Title    string `json:"title"`
	Year     int32  `json:"year"`
	Director string `json:"director"`
}

type MovieClient struct {
	Client  *http.Client
	BaseUrl string
}

var NotFound = errors.New("film not found")

func NewMovieClient(url string) MovieClient {
	return MovieClient{
		Client:  http.DefaultClient,
		BaseUrl: url,
	}
}

func (c *MovieClient) GetMovie(ctx context.Context, id int) (Movieload, error) {
	ctxT, cencel := context.WithTimeout(ctx, 5*time.Second)
	defer cencel()

	url := fmt.Sprintf("%s/%d/info.0.json", c.BaseUrl, id)

	req, err := http.NewRequestWithContext(ctxT, "GET", url, nil)
	if err != nil {
		return Movieload{}, fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return Movieload{}, fmt.Errorf("cannot get film %d error: %w", id, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return Movieload{}, NotFound
	}
	if resp.StatusCode != http.StatusOK {
		return Movieload{}, fmt.Errorf(
			"upstream returned status %d",
			resp.StatusCode,
		)
	}

	defer resp.Body.Close()

	var movie Movieload

	erro := json.NewDecoder(resp.Body).Decode(&movie)
	if erro != nil {
		return Movieload{}, fmt.Errorf("error with parse json: %w", erro)
	}
	return movie, nil
}
