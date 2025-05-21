// Маршрутизатор.

package server

import (
	"content/internal/handlers"
	"content/internal/models/consts"
	"content/internal/tools/config"
	"content/internal/tools/database"
	"content/internal/tools/search"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// server - структура сервера.
type serv struct {
	db     *database.Database
	search *search.Search
}

func NewServer(cfg *config.Config, db *database.Database, search *search.Search) (*http.Server, error) {
	s := serv{
		db:     db,
		search: search,
	}

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

	r.Get("/materials", handlers.GetMaterials(s.db, s.search))
	r.Get("/materials/{"+consts.IdParam+"}", handlers.GetMaterial(s.db))
	r.Post("/materials", handlers.AddMaterial(s.db, s.search))
	r.Put("/materials", handlers.UpdateMaterial(s.db, s.search))
	r.Delete("/materials/{"+consts.IdParam+"}", handlers.DeleteMaterial(s.db, s.search))
	r.Patch("/materials/{"+consts.IdParam+"}/restore", handlers.RestoreMaterial(s.db, s.search))

	r.Get("/materials/search", handlers.SearchMaterials(s.search))
	r.Post("/materials/search/{"+consts.IdParam+"}", handlers.PutMaterialSearch(s.search))
	r.Delete("/materials/search/{"+consts.IdParam+"}", handlers.DeleteMaterialSearch(s.search))

	r.Get("/categories", handlers.GetMainCategories(s.db))
	r.Get("/categories/{"+consts.IdParam+"}/subsidiaries", handlers.GetSubsidiariesCategories(s.db))
	r.Get("/categories/{"+consts.IdParam+"}", handlers.GetCategory(s.db))
	r.Post("/categories", handlers.AddCategory(s.db))
	r.Put("/categories", handlers.UpdateCategory(s.db))
	r.Delete("/categories/{"+consts.IdParam+"}", handlers.DeleteCategory(s.db, s.search))

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Domain, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}, nil
}
