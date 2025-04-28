package fipcontextrepo

import (
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/skasyn/fip-context/fipservice"
)

type FipSongModel struct {
	ID           int64     `json:"id"`
	Title        string    `json:"title"`
	Interpreters []string  `json:"interpreters"`
	Genres       []string  `json:"genres"`
	CreatedAt    time.Time `json:"created_at"`
}

type FipSongRepository struct {
	DB *sql.DB
}

func NewFipSongRepository(db *sql.DB) *FipSongRepository {
	return &FipSongRepository{DB: db}
}

func (r *FipSongRepository) Create(song *FipSongModel) error {
	query := `
	INSERT INTO songs (title, interpreters, genres)
	VALUES ($1, $2, $3)
	RETURNING id, created_at`

	return r.DB.QueryRow(
		query,
		song.Title,
		pq.Array(song.Interpreters),
		pq.Array(song.Genres),
	).Scan(&song.ID, &song.CreatedAt)
}

func (r *FipSongRepository) FipSongToFipSongModel(song *fipservice.FipSong) *FipSongModel {
	var songAsModel FipSongModel

	songAsModel.Title = song.Name
	songAsModel.Interpreters = song.Interpreters
	songAsModel.Genres = song.Genres

	return &songAsModel
}
