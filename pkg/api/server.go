package api

import (
	"context"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"

	"effectiveMobileTest/config"
	"effectiveMobileTest/pkg/api/handler"
	"effectiveMobileTest/pkg/api/middlewares"
	"effectiveMobileTest/pkg/service/music"
	"effectiveMobileTest/utils"
)

const (
	maxHeaderBytes = 1 << 20 // 1 MB
	readTimeout    = 10 * time.Second
	writeTimeout   = 10 * time.Second
)

type Server struct {
	httpServer *http.Server
	router     *mux.Router
}

func NewServer(ctx utils.MyContext, config config.Config) *Server {
	router := mux.NewRouter()

	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	wrappedRouter := middlewares.RecoveryMiddleware(ctx, router)

	return &Server{
		httpServer: &http.Server{
			Addr:           config.ServerPort,
			MaxHeaderBytes: maxHeaderBytes,
			ReadTimeout:    readTimeout,
			WriteTimeout:   writeTimeout,
			Handler:        wrappedRouter,
		},
		router: router,
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}

func (s *Server) HandleMusic(ctx utils.MyContext, service music.MusicService) {
	s.router.HandleFunc("/api/songs", handler.AddSong(ctx, service)).Methods(http.MethodPost)
	s.router.HandleFunc("/api/songs", handler.GetSongs(ctx, service)).Methods(http.MethodGet)
	s.router.HandleFunc("/api/songs/{songId}/lyrics", handler.GetLyrics(ctx, service)).Methods(http.MethodGet)
	s.router.HandleFunc("/api/songs/{id}", handler.UpdateSong(ctx, service)).Methods(http.MethodPut)
	s.router.HandleFunc("/api/songs/{id}", handler.DeleteSong(ctx, service)).Methods(http.MethodDelete)
}
