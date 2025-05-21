// Маршрутизатор.

package server

import (
	"file-server/internal/handlers"
	"file-server/internal/models/consts"
	"file-server/internal/tools/config"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func NewServer(cfg *config.Config) (*http.Server, error) {
	// Создание и настройка маршрутизатора.
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	r.Get("/photos/{"+consts.FileNameWithExtensionParam+"}", handlers.GetPhoto())
	r.Post("/photos/{"+consts.FileNameParam+"}", handlers.AddPhoto())
	r.Delete("/photos/{"+consts.FileNameWithExtensionParam+"}", handlers.DeletePhoto())

	r.Get("/videos/{"+consts.FileNameParam+"}/{"+consts.FileNamePartParam+"}", handlers.GetVideo())
	r.Post("/videos/{"+consts.FileNameParam+"}", handlers.AddVideo())
	r.Delete("/videos/{"+consts.FileNameParam+"}", handlers.DeleteVideo())

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Domain, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}, nil
}
