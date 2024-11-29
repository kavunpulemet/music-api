package mappers

import (
	"effectiveMobileTest/models"
	dbmodels "effectiveMobileTest/pkg/repository/models"
)

func MapFromSongs(repositorySongs []dbmodels.Song) []models.Song {
	serviceSongs := make([]models.Song, len(repositorySongs))
	for i, repositorySong := range repositorySongs {
		serviceSongs[i] = models.Song{
			Id:          repositorySong.Id,
			Group:       repositorySong.Group,
			Title:       repositorySong.Title,
			ReleaseDate: repositorySong.ReleaseDate,
			Text:        repositorySong.Text,
			Link:        repositorySong.Link,
		}
	}

	return serviceSongs
}
