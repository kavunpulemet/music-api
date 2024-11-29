package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"

	"effectiveMobileTest/models"
	"effectiveMobileTest/pkg/service/music"
	"effectiveMobileTest/utils"
)

// AddSong godoc
// @Summary Create a new song
// @Description Create a new song in the database with the provided details
// @Tags songs
// @Accept json
// @Produce json
// @Param song body models.AddSongRequest true "Song Data"
// @Success 200 {object} models.AddSongResponse "id of the created song"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /songs [post]
func AddSong(ctx utils.MyContext, service music.MusicService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.Logger.Debugf("AddSong handler invoked")

		var req models.AddSongRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		if req.Group == "" || req.Title == "" {
			utils.NewErrorResponse(ctx, w, "missing required fields: group or song")
			return
		}

		ctx.Logger.Debugf("creating song: %+v", req)
		id, err := service.Create(ctx, req)
		if err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		ctx.Logger.Infof("song created successfully with id=%s", id)

		response := models.AddSongResponse{Id: id}

		if err = utils.WriteResponse(w, http.StatusOK, response); err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		ctx.Logger.Debugf("response sent successfully for AddSong")
	}
}

// GetSongs godoc
// @Summary Get a list of songs
// @Description Retrieve songs from the database based on the provided filters
// @Tags songs
// @Accept  json
// @Produce  json
// @Param group query string false "Group name"
// @Param title query string false "Song title"
// @Param releaseDate query string false "Release date (YYYY-MM-DD)"
// @Param text query string false "Text"
// @Param link query string false "Link"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of songs per page" default(10)
// @Success 200 {array} models.Song "List of songs"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /songs [get]
func GetSongs(ctx utils.MyContext, service music.MusicService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.Logger.Debugf("GetSongs handler invoked")

		var (
			err         error
			releaseDate time.Time
		)

		releaseDateStr := r.URL.Query().Get("releaseDate")
		if releaseDateStr != "" {
			ctx.Logger.Debugf("parsing releaseDate: %s", releaseDateStr)
			releaseDate, err = time.Parse("2006-01-02", releaseDateStr)
			if err != nil {
				utils.NewErrorResponse(ctx, w, "invalid releaseDate format")
				return
			}
		}

		filter := models.SongFilter{
			Group:       r.URL.Query().Get("group"),
			Title:       r.URL.Query().Get("title"),
			ReleaseDate: releaseDate,
			Text:        r.URL.Query().Get("text"),
			Link:        r.URL.Query().Get("link"),
			Page:        getPageFromQuery(r),
			Limit:       getLimitFromQuery(r),
		}

		ctx.Logger.Debugf("filters applied: %+v", filter)

		songs, err := service.GetSongs(ctx, filter)
		if err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		if len(songs) == 0 {
			ctx.Logger.Warnf("no songs found for filters: %+v", filter)
		} else {
			ctx.Logger.Infof("songs retrieved successfully, count: %d", len(songs))
		}

		if err = utils.WriteResponse(w, http.StatusOK, songs); err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		ctx.Logger.Debugf("response sent successfully for GetSongs")
	}
}

// GetLyrics godoc
// @Summary Get lyrics for a specific song
// @Description Retrieve lyrics for a song by its ID based on the provided filters
// @Tags songs
// @Accept  json
// @Produce  json
// @Param songId path string true "Song ID"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Number of couplets per page" default(10)
// @Success 200 {object} models.SongLyricsResponse "Song lyrics"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /songs/{songId}/lyrics [get]
func GetLyrics(ctx utils.MyContext, service music.MusicService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.Logger.Debugf("GetLyrics handler invoked")

		songId := mux.Vars(r)["songId"]

		ctx.Logger.Debugf("fetching lyrics for songId: %s", songId)

		filter := models.LyricsFilter{
			SongId: songId,
			Page:   getPageFromQuery(r),
			Limit:  getLimitFromQuery(r),
		}

		ctx.Logger.Debugf("filters applied: %+v", filter)

		lyrics, err := service.GetLyrics(ctx, filter)
		if err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		ctx.Logger.Infof("lyrics retrieved successfully for songId: %s", songId)

		response := models.SongLyricsResponse{Lyrics: lyrics}

		if err = utils.WriteResponse(w, http.StatusOK, response); err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		ctx.Logger.Debugf("response sent successfully for GetLyrics")
	}
}

// UpdateSong godoc
// @Summary Update an existing song
// @Description Update the details of an existing song
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path string true "Song ID"
// @Param song body models.UpdateSongRequest true "Song Data"
// @Success 200 {object} utils.StatusResponse "Response indicating the status of the operation"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /songs/{id} [put]
func UpdateSong(ctx utils.MyContext, service music.MusicService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.Logger.Debugf("UpdateSong handler invoked")

		var input models.UpdateSongRequest

		id := mux.Vars(r)["id"]

		ctx.Logger.Debugf("decoding update data for songId=%s", id)

		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		ctx.Logger.Debugf("updating song with data: %+v", input)

		if err := service.Update(ctx, id, input); err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		ctx.Logger.Infof("song updated successfully with id=%s", id)

		if err := utils.WriteResponse(w, http.StatusOK, utils.StatusResponse{Status: "ok"}); err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		ctx.Logger.Debugf("response sent successfully for UpdateSong")
	}
}

// DeleteSong godoc
// @Summary Delete a song
// @Description Delete a song from the system by its ID
// @Tags songs
// @Accept  json
// @Produce  json
// @Param id path string true "Song ID"
// @Success 200 {object} utils.StatusResponse "Response indicating the status of the operation"
// @Failure 400 {object} utils.ErrorResponse "Bad Request"
// @Failure 404 {object} utils.ErrorResponse "Not Found"
// @Failure 500 {object} utils.ErrorResponse "Internal Server Error"
// @Router /songs/{id} [delete]
func DeleteSong(ctx utils.MyContext, service music.MusicService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx.Logger.Debugf("DeleteSong handler invoked")

		id := mux.Vars(r)["id"]

		ctx.Logger.Debugf("deleting song with id=%s", id)

		if err := service.Delete(ctx, id); err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		ctx.Logger.Infof("song deleted successfully with id=%s", id)

		if err := utils.WriteResponse(w, http.StatusOK, utils.StatusResponse{Status: "ok"}); err != nil {
			utils.NewErrorResponse(ctx, w, err.Error())
			return
		}

		ctx.Logger.Debugf("response sent successfully for DeleteSong")
	}
}

func getPageFromQuery(r *http.Request) int {
	page := r.URL.Query().Get("page")
	if page == "" {
		return 1
	}
	p, err := strconv.Atoi(page)
	if err != nil {
		return 1
	}
	return p
}

func getLimitFromQuery(r *http.Request) int {
	limit := r.URL.Query().Get("limit")
	if limit == "" {
		return 10
	}
	l, err := strconv.Atoi(limit)
	if err != nil {
		return 10
	}
	return l
}
