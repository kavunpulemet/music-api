package mappers

import (
	"fmt"
	"time"

	"effectiveMobileTest/models"
	dbmodels "effectiveMobileTest/pkg/repository/models"
)

func MapUpdateToSong(id string, updateSong models.UpdateSongRequest) dbmodels.Song {
	return dbmodels.Song{
		Id:          id,
		Group:       updateSong.Group,
		Title:       updateSong.Title,
		ReleaseDate: updateSong.ReleaseDate,
		Text:        updateSong.Text,
		Link:        updateSong.Link,
	}
}

func MapDetailsToSong(req models.AddSongRequest, details models.SongDetails) (dbmodels.Song, error) {
	parsedDate, err := time.Parse("02.01.2006", details.ReleaseDate)
	if err != nil {
		return dbmodels.Song{}, fmt.Errorf("invalid release date format: %w", err)
	}

	return dbmodels.Song{
		Id:          "",
		Group:       req.Group,
		Title:       req.Title,
		ReleaseDate: parsedDate,
		Text:        details.Text,
		Link:        details.Link,
	}, nil
}

func MapToFilter(serviceFilter models.SongFilter) dbmodels.SongFilter {
	return dbmodels.SongFilter{
		Group:       serviceFilter.Group,
		Title:       serviceFilter.Title,
		ReleaseDate: serviceFilter.ReleaseDate,
		Text:        serviceFilter.Text,
		Link:        serviceFilter.Link,
		Page:        serviceFilter.Page,
		Limit:       serviceFilter.Limit,
	}
}
