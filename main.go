package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type Description struct {
	Id       int    `json:"id"`
	Title    string `json:"title"`
	Year     int    `json:"year"`
	Director string `json:"director"`
}

func worker(ctx context.Context, tasks <-chan int, t time.Duration) {
	for {
		select {
		case <-ctx.Done():

			return
		case id, ok := <-tasks:
			if !ok {
				return
			}
			func() {
				newCtx, cancel := context.WithTimeout(ctx, t)
				defer cancel()
				var des Description

				url := fmt.Sprintf("https://homeworksite.site/%d/info.0.json", id)

				req, err := http.NewRequestWithContext(newCtx, http.MethodGet, url, nil)
				if err != nil {
					fmt.Printf("Error: cannot create request\n")
					return
				}

				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					switch {
					case errors.Is(err, context.DeadlineExceeded):
						fmt.Printf("id %d: timeout\n", id)

					case errors.Is(err, context.Canceled):
						fmt.Printf("id: %d, canceled\n", id)

					default:
						fmt.Printf("id: %d, request error", id)
					}

					return
				}
				defer resp.Body.Close()

				if resp.StatusCode != http.StatusOK {
					fmt.Printf("id %d: http error %d\n", id, resp.StatusCode)
				}

				erro := json.NewDecoder(resp.Body).Decode(&des)

				if erro != nil {
					fmt.Printf("Error: crash json!")
					return
				}
				cancel()

				res := fmt.Sprintf("%d — %s — %d — %s", des.Id, des.Title, des.Year, des.Director)

				fmt.Println(res)
			}()
		}

	}
}

func main() {

	from := flag.Int("from", -1, "int")
	to := flag.Int("to", -1, "int")
	workers := flag.Int("workers", 10, "int")
	timeout := flag.Duration("timeout", 5*time.Second, "request timeout")

	flag.Parse()

	if *from == -1 || *to == -1 {
		fmt.Printf("Error: flag from, to is required")
		os.Exit(1)
	}

	if *to < *from {
		fmt.Printf("Error: argument from > to")
		os.Exit(1)
	}

	if *workers <= 0 {
		fmt.Printf("Error: incorrect number of workers")
		os.Exit(1)
	}

	if *timeout <= 0 {
		fmt.Printf("Error: timeout must be greater than 0")
		os.Exit(1)
	}

	tasks := make(chan int)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var wg sync.WaitGroup
	for iw := 0; iw < *workers; iw++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(ctx, tasks, *timeout)
		}()
	}

	for id := *from; id <= *to; id++ {
		select {
		case <-ctx.Done():
			close(tasks)
			wg.Wait()
			return
		case tasks <- id:
		}
	}
	close(tasks)
	wg.Wait()

}
