package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/lib/pq"

	"effectiveMobileTest/config"
	_ "effectiveMobileTest/docs"
	"effectiveMobileTest/utils"
)

// @title Music API
// @version 1.0
// @description Music API that allows you to add, get, update and delete songs
// @host localhost:81
// @BasePath /api
func main() {
	mainCtx := context.Background()
	ctx, cancel := context.WithCancel(mainCtx)
	defer cancel()

	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("failed to read config: %v", err)
	}

	logger, err := utils.NewLogger(cfg.LoggerLevel)
	if err != nil {
		log.Fatalf("failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	app := NewApp(ctx, logger, cfg)

	if err = app.InitDatabase(); err != nil {
		logger.Fatalf("failed to initialize db: %s", err.Error())
	}

	if err = app.RunMigrations(); err != nil {
		logger.Fatalf("failed to run migrations: %s", err.Error())
	}

	app.InitService()

	if err = app.Run(); err != nil {
		logger.Errorf(err.Error())
		return
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)

	logger.Info("app is running")
	<-sigChan

	logger.Info("stopping app")

	if err = app.Shutdown(ctx); err != nil {
		logger.Errorf(err.Error())
		return
	}

	logger.Info("app stopped successfully")
}
