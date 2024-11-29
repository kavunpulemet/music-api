package mappers

import (
	"testing"
	"time"

	"effectiveMobileTest/models"
	dbmodels "effectiveMobileTest/pkg/repository/models"
)

func TestMapUpdateToSong(t *testing.T) {
	id := "123"
	input := models.UpdateSongRequest{
		Group:       "Test Group",
		Title:       "Test Title",
		ReleaseDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Text:        "Some lyrics",
		Link:        "http://example.com",
	}

	expected := dbmodels.Song{
		Id:          "123",
		Group:       "Test Group",
		Title:       "Test Title",
		ReleaseDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Text:        "Some lyrics",
		Link:        "http://example.com",
	}

	result := MapUpdateToSong(id, input)

	if result.Id != expected.Id || result.Group != expected.Group || result.Title != expected.Title ||
		!result.ReleaseDate.Equal(expected.ReleaseDate) || result.Text != expected.Text || result.Link != expected.Link {
		t.Errorf("MapUpdateToSong() = %+v, want %+v", result, expected)
	}
}

func TestMapDetailsToSong(t *testing.T) {
	req := models.AddSongRequest{
		Group: "Test Group",
		Title: "Test Title",
	}
	details := models.SongDetails{
		ReleaseDate: "01.01.2024",
		Text:        "Some lyrics",
		Link:        "http://example.com",
	}

	expectedDate, _ := time.Parse("02.01.2006", "01.01.2024")
	expected := dbmodels.Song{
		Id:          "",
		Group:       "Test Group",
		Title:       "Test Title",
		ReleaseDate: expectedDate,
		Text:        "Some lyrics",
		Link:        "http://example.com",
	}

	result, err := MapDetailsToSong(req, details)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Id != expected.Id || result.Group != expected.Group || result.Title != expected.Title ||
		!result.ReleaseDate.Equal(expected.ReleaseDate) || result.Text != expected.Text || result.Link != expected.Link {
		t.Errorf("MapDetailsToSong() = %+v, want %+v", result, expected)
	}
}

func TestMapToFilter(t *testing.T) {
	input := models.SongFilter{
		Group:       "Test Group",
		Title:       "Test Title",
		ReleaseDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Text:        "Some lyrics",
		Link:        "http://example.com",
		Page:        1,
		Limit:       10,
	}

	expected := dbmodels.SongFilter{
		Group:       "Test Group",
		Title:       "Test Title",
		ReleaseDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		Text:        "Some lyrics",
		Link:        "http://example.com",
		Page:        1,
		Limit:       10,
	}

	result := MapToFilter(input)

	if result.Group != expected.Group || result.Title != expected.Title || !result.ReleaseDate.Equal(expected.ReleaseDate) ||
		result.Text != expected.Text || result.Link != expected.Link || result.Page != expected.Page || result.Limit != expected.Limit {
		t.Errorf("MapToFilter() = %+v, want %+v", result, expected)
	}
}

func TestMapFromSongs(t *testing.T) {
	repositorySongs := []dbmodels.Song{
		{
			Id:          "123",
			Group:       "Group 1",
			Title:       "Song 1",
			ReleaseDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Text:        "Lyrics 1",
			Link:        "http://example1.com",
		},
		{
			Id:          "456",
			Group:       "Group 2",
			Title:       "Song 2",
			ReleaseDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Text:        "Lyrics 2",
			Link:        "http://example2.com",
		},
	}

	expected := []models.Song{
		{
			Id:          "123",
			Group:       "Group 1",
			Title:       "Song 1",
			ReleaseDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Text:        "Lyrics 1",
			Link:        "http://example1.com",
		},
		{
			Id:          "456",
			Group:       "Group 2",
			Title:       "Song 2",
			ReleaseDate: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			Text:        "Lyrics 2",
			Link:        "http://example2.com",
		},
	}

	result := MapFromSongs(repositorySongs)

	if len(repositorySongs) != len(result) {
		t.Errorf("expected %d songs, got %d", len(repositorySongs), len(result))
	}

	for i, expectedSong := range expected {
		if result[i].Id != expectedSong.Id {
			t.Errorf("expected Id %s, got %s", expectedSong.Id, result[i].Id)
		}
		if result[i].Group != expectedSong.Group {
			t.Errorf("expected Group %s, got %s", expectedSong.Group, result[i].Group)
		}
		if result[i].Title != expectedSong.Title {
			t.Errorf("expected Title %s, got %s", expectedSong.Title, result[i].Title)
		}
		if !result[i].ReleaseDate.Equal(expectedSong.ReleaseDate) {
			t.Errorf("expected ReleaseDate %v, got %v", expectedSong.ReleaseDate, result[i].ReleaseDate)
		}
		if result[i].Text != expectedSong.Text {
			t.Errorf("expected Text %s, got %s", expectedSong.Text, result[i].Text)
		}
		if result[i].Link != expectedSong.Link {
			t.Errorf("expected Link %s, got %s", expectedSong.Link, result[i].Link)
		}
	}
}
