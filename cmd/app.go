package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
	"go.uber.org/zap"

	"effectiveMobileTest/config"
	"effectiveMobileTest/pkg/api"
	"effectiveMobileTest/pkg/repository"
	"effectiveMobileTest/pkg/service/music"
	"effectiveMobileTest/utils"
)

type App struct {
	ctx    utils.MyContext
	server *api.Server
	db     *sqlx.DB
	config config.Config
}

func NewApp(ctx context.Context, logger *zap.SugaredLogger, config config.Config) *App {
	return &App{
		ctx:    utils.NewMyContext(ctx, logger),
		config: config,
	}
}

func (a *App) InitDatabase() error {
	a.ctx.Logger.Infof("connecting to DB with: %s", a.config.DBConnectionString)

	var err error

	for i := 0; i < 10; i++ {
		a.db, err = sqlx.Open("postgres", a.config.DBConnectionString)
		if err != nil {
			a.ctx.Logger.Warnf("failed to connect to DB, attempt %d: %v", i+1, err)
			time.Sleep(5 * time.Second)
			continue
		}

		if err = a.db.Ping(); err == nil {
			a.ctx.Logger.Infof("successfully connected to DB")
			break
		}

		a.ctx.Logger.Warnf("DB not ready, attempt %d: %v\n", i+1, err)
		time.Sleep(5 * time.Second)
	}

	return nil
}

func (a *App) RunMigrations() error {
	a.ctx.Logger.Info("running database migrations")

	if err := goose.Up(a.db.DB, "./migrations"); err != nil {
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	a.ctx.Logger.Info("migrations applied successfully")
	return nil
}

func (a *App) InitService() {
	a.ctx.Logger.Info("initializing services")

	s := music.NewMusicService(repository.NewPostgres(a.db), a.config.SongDetailsAPIUrl)
	a.server = api.NewServer(a.ctx, a.config)
	a.server.HandleMusic(a.ctx, s)

	a.ctx.Logger.Info("services initialized successfully")
}

func (a *App) Run() error {
	a.ctx.Logger.Info("starting HTTP server")

	go func() {
		if err := a.server.Run(); err != nil {
			a.ctx.Logger.Fatalf("error occured while running HTTP server: %s", err.Error())
		}
	}()

	a.ctx.Logger.Info("HTTP server is running")
	return nil
}

func (a *App) Shutdown(ctx context.Context) error {
	a.ctx.Logger.Info("shutting down HTTP server")

	err := a.server.Shutdown(ctx)
	if err != nil {
		a.ctx.Logger.Errorf("failed to shut down HTTP server: %v", err)
		return err
	}

	err = a.db.Close()
	if err != nil {
		a.ctx.Logger.Errorf("failed to disconnect from BD: %v", err)
	}

	a.ctx.Logger.Info("server shut down successfully")
	return nil
}
