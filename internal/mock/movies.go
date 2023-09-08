package mock

import (
	"time"

	"greenlight.jattueyi.com/internal/data"
)

var mockMovie = data.Movie{
	ID:        1,
	CreatedAt: time.Now(),
	Title:     "Terminator",
	Year:      1993,
	Runtime:   100,
	Genres:    []string{"Action", "Sci-Fi"},
	Version:   0,
}

var mockMetadata = data.Metadata{
	CurrentPage:  0,
	PageSize:     0,
	FirstPage:    0,
	LastPage:     0,
	TotalRecords: 0,
}

type MovieModel struct{}

func (m MovieModel) Insert(movie *data.Movie) error {
	return nil
}

func (m MovieModel) Get(id int64) (*data.Movie, error) {
	return &mockMovie, nil
}

func (m MovieModel) Update(movie *data.Movie) error {
	return nil
}

func (m MovieModel) Delete(id int64) error {
	return nil
}

func (m MovieModel) GetAll(title string, genres []string, filters data.Filters) ([]*data.Movie, data.Metadata, error) {
	return []*data.Movie{&mockMovie}, mockMetadata, nil
}
