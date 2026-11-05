package movieclient

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/ButyugovMxx/GoPrac/internal/domain"
)

type MovieClient struct {
	Client  *http.Client
	BaseUrl string
}

func NewMovieClient(url string) MovieClient {
	return MovieClient{
		Client:  http.DefaultClient,
		BaseUrl: url,
	}
}

func (c *MovieClient) GetMovie(ctx context.Context, id int) (domain.Movieload, error) {
	ctxT, cencel := context.WithTimeout(ctx, 5*time.Second)
	defer cencel()

	url := fmt.Sprintf("%s/%d/info.0.json", c.BaseUrl, id)

	req, err := http.NewRequestWithContext(ctxT, "GET", url, nil)
	if err != nil {
		return domain.Movieload{}, fmt.Errorf("cannot create request: %w", err)
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return domain.Movieload{}, fmt.Errorf("cannot get film %d error: %w", id, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return domain.Movieload{}, domain.NotFound
	}
	if resp.StatusCode != http.StatusOK {
		return domain.Movieload{}, fmt.Errorf(
			"upstream returned status %d",
			resp.StatusCode,
		)
	}

	defer resp.Body.Close()

	var movie domain.Movieload

	erro := json.NewDecoder(resp.Body).Decode(&movie)
	if erro != nil {
		return domain.Movieload{}, fmt.Errorf("error with parse json: %w", erro)
	}
	return movie, nil
}
