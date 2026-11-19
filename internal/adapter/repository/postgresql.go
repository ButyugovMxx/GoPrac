package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/ButyugovMxx/GoPrac/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)


type MovieDbConnect struct{
	BaseUrl string
	pool *pgxpool.Pool
}

func NewMovieDB(url string) (*MovieDbConnect, error){
	var db MovieDbConnect

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, url)
	if err != nil {
		return nil, fmt.Errorf("Db connection error: %w", err)
	}

	if err != pool.Ping(context.Background()){
		return nil, fmt.Errorf("Db connection: %w", err)
	}

	db.BaseUrl = url
	db.pool = pool
	return &db, nil
}

func (db *MovieDbConnect) GetMovie(ctx context.Context, id int) (domain.Movieload, error){
	var movie domain.Movieload

	query := `SELECT id, title, year, director FROM movies WHERE id = $1`
	err := db.pool.QueryRow(ctx, query, id,
	).Scan(&movie.Id, &movie.Title, &movie.Year, &movie.Director)
	
	if err != nil{
		return domain.Movieload{}, err
	}
	return movie, nil
}