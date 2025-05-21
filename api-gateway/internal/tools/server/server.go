// Маршрутизатор.

package server

import (
	"fmt"
	log "log/slog"
	"net/http"

	"api-gateway/internal/handlers"
	"api-gateway/internal/models/consts"
	"api-gateway/internal/models/errs"
	"api-gateway/internal/tools/config"
	"api-gateway/internal/tools/traffic_limiter"
	"api-gateway/pkg/code"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

// serv - структура сервера.
type serv struct {
	tl *traffic_limiter.TrafficLimiter
}

func NewServer(cfg *config.Config) (*http.Server, error) {
	s := serv{}

	var err error
	s.tl, err = traffic_limiter.NewTrafficLimiter(cfg)
	if err != nil {
		return nil, err
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

	r.Use(s.middleware)

	r.Get("/materials", handlers.GetMaterials())
	r.Get("/materials/{"+consts.IdParam+"}", handlers.GetMaterial())
	r.Post("/materials", handlers.AddMaterial())
	r.Put("/materials", handlers.UpdateMaterial())
	r.Delete("/materials/{"+consts.IdParam+"}", handlers.DeleteMaterial())
	r.Patch("/materials/{"+consts.IdParam+"}/restore", handlers.RestoreMaterial())

	r.Get("/categories", handlers.GetMainCategories())
	r.Get("/categories/{"+consts.IdParam+"}/subsidiaries", handlers.GetSubsidiariesCategories())
	r.Get("/categories/{"+consts.IdParam+"}", handlers.GetCategory())
	r.Post("/categories", handlers.AddCategory())
	r.Put("/categories", handlers.UpdateCategory())
	r.Delete("/categories/{"+consts.IdParam+"}", handlers.DeleteCategory())

	return &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Domain, cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.Timeout,
		WriteTimeout: cfg.Server.Timeout,
	}, nil
}

// middleware - middleware для получения IP-адреса клиента
// и ограничения доступа к приложению, если лимит запросов исчерпан
func (s serv) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Debug("запрос на адрес", log.String("URI", r.RequestURI))

		ipAddress := r.RemoteAddr

		// Если используется прокси-сервер, можно попробовать получить IP-адрес из заголовка X-Forwarded-For
		if xForwardedFor := r.Header.Get("X-Forwarded-For"); xForwardedFor != "" {
			ipAddress = xForwardedFor
		}

		// Проверяем возможность запроса по IP-адресу
		hasOpportunity, err := s.tl.HasOpportunityToRequest(ipAddress)
		if err != nil {
			log.Error("ошибка middleware", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		if !hasOpportunity {
			code.ErrorJSON(w, http.StatusTooManyRequests, errs.ErrTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}
