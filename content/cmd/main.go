package main

import (
	"content/internal/tools/search"
	"context"
	log "log/slog"
	"os"
	"os/signal"
	"syscall"

	"content/internal/models/consts"
	"content/internal/tools/config"
	"content/internal/tools/database"
	"content/internal/tools/server"
	"content/pkg/logger"
)

func main() {
	// Получение конфига.
	cfg, err := config.NewConfig(consts.ContentConfigPath)
	if err != nil {
		log.Error("ошибка прочтения файла конфигруаций", log.Any("error", err))
		return
	}

	// Настройка логов.
	logger.SetLogger(cfg.Logger.Level)

	// Соединение с БД.
	db, err := database.NewDB(cfg)
	if err != nil {
		log.Error("ошибка подключения к базе данных", log.Any("error", err))
		return
	}

	// Соедение с поисковой системой.
	s, err := search.NewSearch(cfg)
	if err != nil {
		log.Error("ошибка подключения к поисковой системе", log.Any("error", err))
		return
	}

	log.Info("Запуск content service")

	// Настройка конфигурации сервера.
	srv, err := server.NewServer(cfg, db, s)
	if err != nil {
		log.Error("ошибка настройки конфигурации сервера", log.Any("error", err))
	}

	// Запуск сервера.
	go func() {
		//	err = srv.ListenAndServeTLS(consts.ContentCertPem, consts.ContentKeyPem)
		err = srv.ListenAndServe()
		if err != nil {
			log.Error("ошибка запуска сервера", log.Any("error", err))
			os.Exit(1)
		}
	}()

	// Завершение работы.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGTERM, syscall.SIGINT)

	<-shutdown
	log.Info("Завершение работы...")

	err = srv.Shutdown(context.Background())
	if err != nil {
		log.Error("ошибка при завершении работы сервера", log.Any("error", err))
	}
}
