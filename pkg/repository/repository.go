package repository

import (
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"

	"database/sql"
	"effectiveMobileTest/pkg/api/utils"
	dbmodels "effectiveMobileTest/pkg/repository/models"
)

type Repository interface {
	Create(ctx utils.MyContext, song dbmodels.Song) error
	GetSongs(ctx utils.MyContext, filter dbmodels.SongFilter) ([]dbmodels.Song, error)
	GetLyrics(ctx utils.MyContext, id string) (string, error)
	Update(ctx utils.MyContext, input dbmodels.Song) error
	Delete(ctx utils.MyContext, id string) error
}

type Postgres struct {
	db *sqlx.DB
}

func NewPostgres(db *sqlx.DB) *Postgres {
	return &Postgres{db: db}
}

//go:embed sql/CreateSong.sql
var createSong string

func (r *Postgres) Create(ctx utils.MyContext, song dbmodels.Song) error {
	ctx.Logger.Debugf("executing CreateSong query: id=%s, group=%s, title=%s", song.Id, song.Group, song.Title)

	_, err := r.db.ExecContext(ctx.Ctx, createSong, song.Id, song.Group, song.Title, song.ReleaseDate, song.Text, song.Link)
	if err != nil {
		return fmt.Errorf("failed to insert song: %w", err)
	}

	ctx.Logger.Infof("song inserted successfully id=%s", song.Id)

	return nil
}

func (r *Postgres) GetSongs(ctx utils.MyContext, filter dbmodels.SongFilter) ([]dbmodels.Song, error) {
	ctx.Logger.Debugf("executing GetSongs query with filter: %+v", filter)

	var query strings.Builder
	query.WriteString("SELECT * FROM songs WHERE 1=1")

	if filter.Group != "" {
		query.WriteString(" AND group_name ILIKE '%" + filter.Group + "%'")
	}
	if filter.Title != "" {
		query.WriteString(" AND title ILIKE '%" + filter.Title + "%'")
	}
	if !filter.ReleaseDate.IsZero() {
		formattedDate := filter.ReleaseDate.Format("2006-01-02")
		query.WriteString(" AND release_date = '" + formattedDate + "'")
	}
	if filter.Text != "" {
		query.WriteString(" AND text ILIKE '%" + filter.Text + "%'")
	}
	if filter.Link != "" {
		query.WriteString(" AND link ILIKE '%" + filter.Link + "%'")
	}

	offset := (filter.Page - 1) * filter.Limit
	query.WriteString(fmt.Sprintf(" LIMIT %d OFFSET %d", filter.Limit, offset))

	var songs []dbmodels.Song
	err := r.db.SelectContext(ctx.Ctx, &songs, query.String())
	if err != nil {
		return nil, err
	}

	ctx.Logger.Debugf("retrieved %d songs", len(songs))

	return songs, nil
}

//go:embed sql/GetLyrics.sql
var getLyrics string

func (r *Postgres) GetLyrics(ctx utils.MyContext, id string) (string, error) {
	ctx.Logger.Debugf("executing GetLyrics query for song id=%s", id)

	var lyrics string
	err := r.db.QueryRowContext(ctx.Ctx, getLyrics, id).Scan(&lyrics)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("song with id %s not found", id)
		}
		return "", fmt.Errorf("failed to fetch lyrics: %w", err)
	}

	ctx.Logger.Debug("lyrics retrieved successfully")

	return lyrics, nil
}

func (r *Postgres) Update(ctx utils.MyContext, input dbmodels.Song) error {
	ctx.Logger.Debugf("executing Update query for song with args: %+v", input)

	var (
		queryBuilder strings.Builder
		args         []interface{}
		argIndex     int
	)

	queryBuilder.WriteString("UPDATE songs SET ")

	if input.Group != "" {
		argIndex++
		queryBuilder.WriteString(fmt.Sprintf("group_name = $%d, ", argIndex))
		args = append(args, input.Group)
	}
	if input.Title != "" {
		argIndex++
		queryBuilder.WriteString(fmt.Sprintf("title = $%d, ", argIndex))
		args = append(args, input.Title)
	}
	if !input.ReleaseDate.IsZero() {
		argIndex++
		queryBuilder.WriteString(fmt.Sprintf("release_date = $%d, ", argIndex))
		args = append(args, input.ReleaseDate)
	}
	if input.Text != "" {
		argIndex++
		queryBuilder.WriteString(fmt.Sprintf("text = $%d, ", argIndex))
		args = append(args, input.Text)
	}
	if input.Link != "" {
		argIndex++
		queryBuilder.WriteString(fmt.Sprintf("link = $%d, ", argIndex))
		args = append(args, input.Link)
	}

	queryStr := queryBuilder.String()
	queryStr = queryStr[:len(queryStr)-2]

	argIndex++
	queryStr += fmt.Sprintf(" WHERE id = $%d RETURNING id;", argIndex)
	args = append(args, input.Id)

	ctx.Logger.Debugf("update query: %s", queryStr)

	var id string
	err := r.db.QueryRowContext(ctx.Ctx, queryStr, args...).Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("song with id %s not found", input.Id)
		}
		return fmt.Errorf("failed to update song: %w", err)
	}

	ctx.Logger.Infof("song updated successfully id=%s", input.Id)

	return nil
}

//go:embed sql/DeleteSong.sql
var deleteSong string

func (r *Postgres) Delete(ctx utils.MyContext, id string) error {
	ctx.Logger.Debugf("checking if song exists id=%s", id)

	var exists bool
	err := r.db.QueryRow("SELECT EXISTS(SELECT 1 FROM songs WHERE id = $1)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check if song exists: %w", err)
	}

	if !exists {
		return fmt.Errorf("song with id %s not found", id)
	}

	ctx.Logger.Debugf("executing Delete query for song id=%s", id)

	_, err = r.db.ExecContext(ctx.Ctx, deleteSong, id)
	if err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	ctx.Logger.Infof("song deleted successfully id=%s", id)

	return nil
}
