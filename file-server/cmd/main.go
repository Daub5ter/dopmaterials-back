package main

import (
	"context"
	log "log/slog"
	"os"
	"os/signal"
	"syscall"

	"file-server/internal/models/consts"
	"file-server/internal/tools/config"
	"file-server/internal/tools/server"
	"file-server/pkg/logger"
)

func main() {
	// Получение конфига.
	cfg, err := config.NewConfig(consts.FileServerConfigPath)
	if err != nil {
		log.Error("ошибка прочтения файла конфигруаций", log.Any("error", err))
		return
	}

	// Настройка логов.
	logger.SetLogger(cfg.Logger.Level)

	log.Info("Запуск file service")

	// Настройка конфигурации сервера.
	srv, err := server.NewServer(cfg)
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
