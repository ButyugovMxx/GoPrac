package config

import "os"

type Config struct {
	MovieBaseUrl string
	Addr         string
}

func Load() *Config {
	config := Config{
		MovieBaseUrl: os.Getenv("MOVIE_BASE_URL"),
		Addr:         os.Getenv("ADDR"),
	}

	if config.MovieBaseUrl == "" {
		config.MovieBaseUrl = "https://homeworksite.site"
	}

	if config.Addr == "" {
		config.Addr = ":8080"
	}

	return &config
}
