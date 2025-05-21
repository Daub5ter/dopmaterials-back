package search

import (
	"context"
	"crypto/tls"
	"fmt"
	log "log/slog"
	"net/http"
	"time"

	"content/internal/tools/config"

	"github.com/opensearch-project/opensearch-go"
	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

type Search struct {
	index   string
	client  *opensearch.Client
	timeout time.Duration
	ctx     context.Context
}

func NewSearch(cfg *config.Config) (*Search, error) {
	client, err := connect(cfg.Search.Username, cfg.Search.Password, cfg.Search.Addresses)
	if err != nil {
		return nil, err
	}

	search := Search{
		index:   cfg.Search.Index,
		client:  client,
		timeout: cfg.Search.Timeout,
		ctx:     context.Background(),
	}

	err = search.createIndexIfNotExists()
	if err != nil {
		return nil, fmt.Errorf("не создается индекс %w", err)
	}

	return &search, nil
}

func connect(username, password string, addresses []string) (*opensearch.Client, error) {
	client, err := opensearch.NewClient(opensearch.Config{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Addresses: addresses,
		Username:  username,
		Password:  password,
	})
	if err != nil {
		return nil, fmt.Errorf("ошибка создания клиента: %w", err)
	}

	var counts int
	for {
		_, err = client.API.Ping()
		if err != nil {
			log.Debug("попытка подключение к OpenSearch")
			counts++
		} else {
			log.Info("подключено к OpenSearch")
			return client, nil
		}

		if counts > 10 {
			log.Error("ошибка подключения к OpenSearch", log.Any("error", err))
			return nil, fmt.Errorf("ошибка подключения к OpenSearch %w", err)
		}

		log.Debug("ожидание 2 секунды...")
		time.Sleep(2 * time.Second)
		continue
	}
}

func (s Search) createIndexIfNotExists() error {
	request := opensearchapi.IndicesExistsRequest{
		Index: []string{s.index},
	}

	response, err := request.Do(context.Background(), s.client)
	if err != nil {
		return fmt.Errorf("ошибка запроса проверки существования индекса: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode == http.StatusNotFound {
		err = s.createIndex()
		if err != nil {
			return err
		}
	}

	return nil
}

func (s Search) createIndex() error {
	request := opensearchapi.IndicesCreateRequest{
		Index: s.index,
	}

	response, err := request.Do(context.Background(), s.client)
	if err != nil {
		return fmt.Errorf("ошибка запроса создания индекса: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		return fmt.Errorf("ошибка создания индекса: %s", response)
	}

	return nil
}
