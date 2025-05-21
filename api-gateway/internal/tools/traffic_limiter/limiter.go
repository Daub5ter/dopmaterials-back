package traffic_limiter

import (
	"context"
	"errors"
	log "log/slog"
	"strconv"
	"time"

	"api-gateway/internal/tools/config"

	"github.com/redis/go-redis/v9"
)

type TrafficLimiter struct {
	maxRequestsTimeLivingIpAddress int
	conn                           *redis.Client
	timeout                        time.Duration
	timeLivingIpAddress            time.Duration
}

func NewTrafficLimiter(cfg *config.Config) (*TrafficLimiter, error) {
	conn := connectToRedis(cfg.Limiter.Addr, cfg.Limiter.Password)

	if conn == nil {
		return nil, errors.New("нет подключения к Redis")
	}

	return &TrafficLimiter{
		conn:                           conn,
		timeout:                        cfg.Limiter.Timeout,
		maxRequestsTimeLivingIpAddress: cfg.Limiter.MaxRequestsTimeLivingIpAddress,
		timeLivingIpAddress:            cfg.Limiter.TimeLivingIpAddress,
	}, nil
}

func (tl TrafficLimiter) HasOpportunityToRequest(ipAddress string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), tl.timeout)
	defer cancel()

	// Получаем значение ключа
	valueRow, err := tl.conn.Get(ctx, ipAddress).Result()
	if err == nil {
		value, errV := strconv.Atoi(valueRow)
		if errV != nil {
			return false, err
		}

		if value >= tl.maxRequestsTimeLivingIpAddress {
			return false, nil
		} else {
			// Увеличиваем значение ключа
			err = tl.conn.Incr(ctx, ipAddress).Err()
			if err != nil {
				return false, err
			}
		}

		return true, nil
	}

	// Если ключ не существует, создаем его с начальным значением 1 и временем жизни 1 минута
	if errors.Is(err, redis.Nil) {
		err = tl.conn.Set(ctx, ipAddress, 1, tl.timeLivingIpAddress).Err()
		if err != nil {
			return false, err
		}
	} else {
		return false, err
	}

	return true, nil
}

// openRedis открывает соединение с redis.
func openRedis(ctx context.Context, addr, password string) (*redis.Client, error) {
	rConn := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       0,
	})

	err := rConn.Ping(ctx).Err()
	if err != nil {
		return nil, err
	}

	return rConn, nil
}

// connectToRedis подключается к redis.
func connectToRedis(addr, password string) *redis.Client {
	var counts int

	for {
		connection, err := openRedis(context.Background(), addr, password)
		if err != nil {
			log.Warn("попытка подключение к Redis...")
			counts++
		} else {
			log.Info("подключено к Redis")
			return connection
		}

		if counts > 10 {
			log.Error("ошибка подключения Redis:", log.Any("error", err))
			return nil
		}

		log.Debug("ожидание 2 секунды...")
		time.Sleep(2 * time.Second)
		continue
	}
}
