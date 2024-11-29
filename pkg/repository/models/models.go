package dbmodels

import "time"

type Song struct {
	Id          string    `db:"id"`
	Group       string    `db:"group_name"`
	Title       string    `db:"title"`
	ReleaseDate time.Time `db:"release_date"`
	Text        string    `db:"text"`
	Link        string    `db:"link"`
}

type SongFilter struct {
	Group       string    `db:"group_name"`
	Title       string    `db:"title"`
	ReleaseDate time.Time `db:"release_date"`
	Text        string    `db:"text"`
	Link        string    `db:"link"`
	Page        int       `db:"page"`
	Limit       int       `db:"limit"`
}
