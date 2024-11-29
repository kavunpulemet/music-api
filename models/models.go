package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type Song struct {
	Id          string    `json:"id"`
	Group       string    `json:"group"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

type AddSongRequest struct {
	Group string `json:"group"`
	Title string `json:"song"`
}

type AddSongResponse struct {
	Id string `json:"id"`
}

type SongDetails struct {
	ReleaseDate string `json:"releaseDate"`
	Text        string `json:"text"`
	Link        string `json:"link"`
}

type UpdateSongRequest struct {
	Group       string    `json:"group"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"releaseDate" example:"yyyy-mm-dd"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
}

func (s *UpdateSongRequest) UnmarshalJSON(data []byte) error {
	var err error

	temp := struct {
		Group       string `json:"group"`
		Title       string `json:"title"`
		ReleaseDate string `json:"releaseDate"`
		Text        string `json:"text"`
		Link        string `json:"link"`
	}{}

	if err = json.Unmarshal(data, &temp); err != nil {
		return err
	}

	var parsedDate time.Time
	if temp.ReleaseDate != "" {
		parsedDate, err = time.Parse("2006-01-02", temp.ReleaseDate)
		if err != nil {
			return fmt.Errorf("invalid date format, try yyyy-mm-dd: %w", err)
		}
	}

	s.Group = temp.Group
	s.Title = temp.Title
	s.ReleaseDate = parsedDate
	s.Text = temp.Text
	s.Link = temp.Link

	return nil
}

type SongFilter struct {
	Group       string    `json:"group"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"releaseDate"`
	Text        string    `json:"text"`
	Link        string    `json:"link"`
	Page        int       `json:"page"`
	Limit       int       `json:"limit"`
}

type LyricsFilter struct {
	SongId string `json:"songId"`
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
}

type SongLyricsResponse struct {
	Lyrics string `json:"lyrics"`
}
