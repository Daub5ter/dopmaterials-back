package main

import (
	"context"
	log "log/slog"
	"os"
	"os/signal"
	"syscall"

	"api-gateway/internal/models/consts"
	"api-gateway/internal/tools/config"
	"api-gateway/internal/tools/server"
	"api-gateway/pkg/logger"
)

func main() {
	// Получение конфига.
	cfg, err := config.NewConfig(consts.ApiGatewayConfigPath)
	if err != nil {
		log.Error("ошибка прочтения файла конфигруаций", log.Any("error", err))
		return
	}

	// Настройка логов.
	logger.SetLogger(cfg.Logger.Level)

	log.Info("Запуск api-gateway")

	// Настройка конфигурации сервера.
	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Error("ошибка настройки конфигурации сервера", log.Any("error", err))
	}

	// Запуск сервера.
	go func() {
		//err = srv.ListenAndServeTLS(consts.ApiGatewayCertPem, consts.ApiGatewayKeyPem)
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
