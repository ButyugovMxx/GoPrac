package config

import "os"

type Config struct {
	MovieBaseUrl string
	Addr         string
	DbUrl        string
}

func Load() Config {
	config := Config{
		MovieBaseUrl: os.Getenv("MOVIE_BASE_URL"),
		Addr:         os.Getenv("HTTP_ADDR"),
		DbUrl:        os.Getenv("DATABASE_URL"),
	}

	if config.MovieBaseUrl == "" {
		config.MovieBaseUrl = "https://homeworksite.site"
	}

	if config.Addr == "" {
		config.Addr = ":8080"
	}

	return config
}
