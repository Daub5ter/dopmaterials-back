package database

import (
	"context"
	"errors"
	"fmt"
	log "log/slog"
	"time"

	"content/internal/tools/config"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Database - это представление БД.
type Database struct {
	conn    *pgxpool.Pool
	timeout time.Duration
	ctx     context.Context
}

// NewDB создает новую структуру DB.
func NewDB(cfg *config.Config) (*Database, error) {
	connectDSN := fmt.Sprintf("user=%s password=%s host=%s port=%v dbname=%s sslmode=%s",
		cfg.Database.User, cfg.Database.Password, cfg.Database.Host, cfg.Database.Port, cfg.Database.DBName, cfg.Database.SSLMode)

	conn := ConnectToDB(connectDSN)

	if conn == nil {
		return nil, errors.New("нет подключения к Postgres")
	}

	dbase := Database{
		conn:    conn,
		timeout: cfg.Database.Timeout,
		ctx:     context.Background(),
	}

	if err := dbase.createTablesIfNotExists(); err != nil {
		return nil, err
	}

	return &dbase, nil
}

// createTablesIfNotExists создает таблицы, если они не созданы.
func (db Database) createTablesIfNotExists() error {
	if _, err := db.conn.Exec(db.ctx, createCategories); err != nil {
		return err
	}

	if _, err := db.conn.Exec(db.ctx, createMaterials); err != nil {
		return err
	}

	return nil
}

// openDB открывает соединение с pgsql.
func openDB(dsn string) (*pgxpool.Pool, error) {
	configPGX, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, err
	}

	configPGX.MaxConns = 100              // Максимальное количество соединений
	configPGX.MinConns = 10               // Минимальное количество соединений
	configPGX.MaxConnIdleTime = time.Hour // Максимальное время простоя соединения

	conn, err := pgxpool.ConnectConfig(context.Background(), configPGX)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

// ConnectToDB подключается к pgsql.
func ConnectToDB(dsn string) *pgxpool.Pool {
	var counts int

	for {
		connection, err := openDB(dsn)
		if err != nil {
			log.Debug("попытка подключение к Postgres...")
			counts++
		} else {
			log.Info("подключено к Postgres!")
			return connection
		}

		if counts > 10 {
			log.Error("ошибка подключения к postgres", log.Any("error", err))
			return nil
		}

		log.Debug("ожидание 2 секунды...")
		time.Sleep(2 * time.Second)
		continue
	}
}

// CloseConnection закрывает соединение с pgsql.
func (db Database) CloseConnection() error {
	db.conn.Close()

	return nil
}
