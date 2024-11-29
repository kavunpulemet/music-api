package music

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"

	"effectiveMobileTest/models"
	"effectiveMobileTest/pkg/api/utils"
	"effectiveMobileTest/pkg/repository"
	"effectiveMobileTest/pkg/service/mappers"
)

type MusicService interface {
	Create(ctx utils.MyContext, req models.AddSongRequest) (string, error)
	GetSongs(ctx utils.MyContext, filter models.SongFilter) ([]models.Song, error)
	GetLyrics(ctx utils.MyContext, filter models.LyricsFilter) (string, error)
	Update(ctx utils.MyContext, id string, input models.UpdateSongRequest) error
	Delete(ctx utils.MyContext, id string) error
}

type ImplMusic struct {
	repo              repository.Repository
	client            *http.Client
	songDetailsAPIUrl string
}

func NewMusicService(repo repository.Repository, songDetailsAPIUrl string) *ImplMusic {
	return &ImplMusic{
		repo:              repo,
		client:            &http.Client{Timeout: 10 * time.Second},
		songDetailsAPIUrl: songDetailsAPIUrl,
	}
}

func (s *ImplMusic) Create(ctx utils.MyContext, req models.AddSongRequest) (string, error) {
	songId := uuid.New().String()

	ctx.Logger.Debugf("fetching song details for group=%s, title=%s", req.Group, req.Title)

	details, err := s.FetchSongDetails(req.Group, req.Title)
	if err != nil {
		return "", fmt.Errorf("failed to fetch song details: %w", err)
	}

	ctx.Logger.Debug("mapping song details to song model")

	song, err := mappers.MapDetailsToSong(req, details)
	if err != nil {
		return "", fmt.Errorf("failed to map song details to song: %w", err)
	}

	song.Id = songId

	ctx.Logger.Debugf("saving song with id=%s", songId)

	err = s.repo.Create(ctx, song)
	if err != nil {
		return "", fmt.Errorf("failed to save song: %w", err)
	}

	return songId, nil
}

func (s *ImplMusic) GetSongs(ctx utils.MyContext, filter models.SongFilter) ([]models.Song, error) {
	ctx.Logger.Debugf("retrieving songs with filter: %+v", filter)

	dbSongs, err := s.repo.GetSongs(ctx, mappers.MapToFilter(filter))
	if err != nil {
		return nil, fmt.Errorf("failed to get songs: %w", err)
	}

	ctx.Logger.Debug("mapping database songs to service songs")

	return mappers.MapFromSongs(dbSongs), nil
}

func (s *ImplMusic) GetLyrics(ctx utils.MyContext, filter models.LyricsFilter) (string, error) {
	ctx.Logger.Debugf("retrieving lyrics for song id=%s with pagination page=%d, limit=%d", filter.SongId, filter.Page, filter.Limit)

	lyrics, err := s.repo.GetLyrics(ctx, filter.SongId)
	if err != nil {
		return "", fmt.Errorf("failed to get lyrics: %w", err)
	}

	lyricsLines := strings.Split(lyrics, "\n\n")
	start := (filter.Page - 1) * filter.Limit
	end := start + filter.Limit

	if start > len(lyricsLines)-1 {
		return "", fmt.Errorf("no lyrics found for the specified page: invalid pagination parameters")
	}
	if end > len(lyricsLines) {
		end = len(lyricsLines)
	}

	paginatedLines := lyricsLines[start:end]
	result := strings.Join(paginatedLines, "\n\n")

	ctx.Logger.Debugf("returning paginated lyrics: %s", result)

	return result, nil
}

func (s *ImplMusic) Update(ctx utils.MyContext, id string, input models.UpdateSongRequest) error {
	ctx.Logger.Debugf("updating song with id=%s", id)

	if input.Group == "" && input.Title == "" && input.ReleaseDate.IsZero() &&
		input.Text == "" && input.Link == "" {
		return fmt.Errorf("invalid or missing fields in the request body")
	}

	err := s.repo.Update(ctx, mappers.MapUpdateToSong(id, input))
	if err != nil {
		return fmt.Errorf("failed to update song: %w", err)
	}

	return nil
}

func (s *ImplMusic) Delete(ctx utils.MyContext, id string) error {
	ctx.Logger.Debugf("deleting song with id=%s", id)

	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	return nil
}

func (s *ImplMusic) FetchSongDetails(group, song string) (models.SongDetails, error) {
	encodedGroup := url.QueryEscape(group)
	encodedSong := url.QueryEscape(song)
	url := fmt.Sprintf("%s/info?group=%s&song=%s", s.songDetailsAPIUrl, encodedGroup, encodedSong)

	log.Printf("fetching song details: URL=%s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return models.SongDetails{}, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := s.client.Do(req)
	if err != nil {
		return models.SongDetails{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.SongDetails{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var details models.SongDetails
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return models.SongDetails{}, fmt.Errorf("failed to decode response: %w", err)
	}

	return details, nil
}
