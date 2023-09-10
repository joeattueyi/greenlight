package data

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/testcontainers/testcontainers-go"

	"database/sql"

	"github.com/testcontainers/testcontainers-go/wait"
)

func TestMovieModel(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	req := testcontainers.ContainerRequest{
		Image:        "postgres:12",
		ExposedPorts: []string{"5432/tcp"},
		AutoRemove:   true,
		Env: map[string]string{
			"POSTGRES_USER":     "greenlight",
			"POSTGRES_PASSWORD": "pa55word",
			"POSTGRES_DB":       "greenlight",
		},
		WaitingFor: wait.ForListeningPort("5432/tcp"),
	}

	postgres, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		t.Fatal(err)
	}
	defer postgres.Terminate(ctx)

	db, err := sql.Open("postgres", "postgres://greenlight:pa55word@localhost/greenlight?sslmode=disable")
	if err != nil {
		t.Fatal(err)
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(25)
	duration, err := time.ParseDuration("15m")
	if err != nil {
		t.Fatal(err)
	}
	db.SetConnMaxIdleTime(duration)

	if err = db.PingContext(ctx); err != nil {
		t.Fatal(err)
	}

	var dbMovieModel = DBMovieModel{DB: db}
	_, err = dbMovieModel.DB.ExecContext(ctx, "DELETE FROM movies")
	if err != nil {
		t.Fatal(err)
	}

	t.Run("Insert", func(t *testing.T) {
		movie := Movie{
			Title:   "Titanic",
			Year:    1997,
			Runtime: 100,
			Genres:  []string{"Romance"},
		}

		err := dbMovieModel.Insert(&movie)
		if err != nil {
			t.Fatal(err)
		}

		_, err = dbMovieModel.DB.ExecContext(ctx, "DELETE FROM movies")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Insert & Get", func(t *testing.T) {
		movie := Movie{
			Title:   "Titanic",
			Year:    1997,
			Runtime: 100,
			Genres:  []string{"Romance"},
		}

		err := dbMovieModel.Insert(&movie)
		if err != nil {
			t.Fatal(err)
		}

		gotMovie, err := dbMovieModel.Get(movie.ID)
		if err != nil {
			t.Fatal(err)
		}

		if gotMovie.Title != movie.Title {
			t.Fatal(err)
		}

		_, err = dbMovieModel.DB.ExecContext(ctx, "DELETE FROM movies")
		if err != nil {
			t.Fatal(err)
		}

	})

	t.Run("Insert & Update", func(t *testing.T) {
		movie := Movie{
			Title:   "Titanic",
			Year:    1997,
			Runtime: 100,
			Genres:  []string{"Romance"},
		}

		err := dbMovieModel.Insert(&movie)
		if err != nil {
			t.Fatal(err)
		}

		movie.Title = "Titanic II"
		err = dbMovieModel.Update(&movie)
		if err != nil {
			t.Fatal(err)
		}

		_, err = dbMovieModel.DB.ExecContext(ctx, "DELETE FROM movies")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("Insert & Delete", func(t *testing.T) {
		movie := Movie{
			Title:   "Titanic",
			Year:    1997,
			Runtime: 100,
			Genres:  []string{"Romance"},
		}

		err := dbMovieModel.Insert(&movie)
		if err != nil {
			t.Fatal(err)
		}

		err = dbMovieModel.Delete(movie.ID)
		if err != nil {
			dbMovieModel.DB.ExecContext(ctx, "DELETE FROM movies")
			t.Fatal(err)
		}
	})

	t.Run("Insert Many & Get All", func(t *testing.T) {
		movies := []Movie{
			{
				Title:   "Avatar",
				Year:    2009,
				Runtime: 200,
				Genres:  []string{"Sci-Fi", "Animation"},
			},
			{
				Title:   "Avengers: Endgame",
				Year:    2019,
				Runtime: 180,
				Genres:  []string{"Action"},
			},
			{
				Title:   "Avatar: Way of Water",
				Year:    2022,
				Runtime: 230,
				Genres:  []string{"Sci-Fi", "Animation"},
			},
			{
				Title:   "Titanic",
				Year:    1997,
				Runtime: 100,
				Genres:  []string{"Romance"},
			},
			{
				Title:   "Star Wars: The Force Awakens",
				Year:    2015,
				Runtime: 190,
				Genres:  []string{"Sci-Fi"},
			},
		}

		for _, movie := range movies {
			err := dbMovieModel.Insert(&movie)
			if err != nil {
				t.Fatal(err)
			}
		}

		var input struct {
			Title  string
			Genres []string
			Filters
		}

		input.Title = ""
		input.Genres = []string{}
		input.Filters.Page = 1
		input.Filters.PageSize = 20
		input.Filters.Sort = "-id"
		input.Filters.SortSafelist = []string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}

		gotMovies, _, err := dbMovieModel.GetAll(input.Title, input.Genres, input.Filters)
		if err != nil {
			t.Fatal(err)
		}

		for _, m := range gotMovies {
			fmt.Printf("%+v\n", m)
		}

		if gotMovies[0].ID < gotMovies[len(gotMovies)-1].ID {
			t.Fatal(err)
		}
	})
}
