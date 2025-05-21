package search

import (
	"bytes"
	"content/internal/models/errs"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"content/internal/models"
	"content/internal/models/consts"

	"github.com/opensearch-project/opensearch-go/opensearchapi"
)

func (s Search) PutMaterial(material models.MaterialSearch) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()

	document, err := json.Marshal(material)
	if err != nil {
		return fmt.Errorf("ошибка преобразования в json %w", err)
	}

	request := opensearchapi.IndexRequest{
		Index:      s.index,
		DocumentID: material.Id,
		Body:       bytes.NewReader(document),
		Refresh:    "true",
	}

	response, err := request.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("ошибка запроса добавления документа: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusCreated && response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusBadRequest {
			return errs.ErrBadRequest
		}
		if response.StatusCode == http.StatusNotFound {
			return errs.ErrNotFound
		}

		return fmt.Errorf("ошибка создания документа: %s", response)
	}

	return nil
}

func (s Search) DeleteMaterial(materialId string) error {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()

	deleteRequest := opensearchapi.DeleteRequest{
		Index:      s.index,
		DocumentID: materialId,
	}

	response, err := deleteRequest.Do(ctx, s.client)
	if err != nil {
		return fmt.Errorf("ошибка удаления документа: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusBadRequest {
			return errs.ErrBadRequest
		}
		if response.StatusCode == http.StatusNotFound {
			return errs.ErrNotFound
		}

		return fmt.Errorf("ошибка удаления документа: %s", response)
	}

	return nil
}

func (s Search) SearchMaterials(findPart string, categoryId *uint32, offset uint32) (ids []string, err error) {
	ctx, cancel := context.WithTimeout(s.ctx, s.timeout)
	defer cancel()

	matchQuery :=
		`{
			"size": %v,
			"from": %v,
			"query": {
				"multi_match": {
					"query": "%s",
					"fuzziness": "AUTO",
					"fields": ["name^2", "description"]
				}
			},`

	categoryIdQuery := `	
			"post_filter": {
				"term": {
					"category_id": %v
				}
			},`

	sourceQuery := `
		"_source": false
	}`

	var query string
	if categoryId == nil {
		query = fmt.Sprintf(matchQuery+sourceQuery,
			consts.LimitMaterials,
			offset,
			findPart,
		)
	} else {
		query = fmt.Sprintf(matchQuery+categoryIdQuery+sourceQuery,
			consts.LimitMaterials,
			offset,
			findPart,
			*categoryId,
		)
	}

	search := opensearchapi.SearchRequest{
		Index: []string{s.index},
		Body:  strings.NewReader(query),
	}

	response, err := search.Do(ctx, s.client)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса поиска документов: %w", err)
	}
	defer func() { _ = response.Body.Close() }()

	if response.StatusCode != http.StatusOK {
		if response.StatusCode == http.StatusBadRequest {
			return nil, errs.ErrBadRequest
		}
		if response.StatusCode == http.StatusNotFound {
			return nil, errs.ErrNotFound
		}

		return nil, fmt.Errorf("ошибка поиска документов: %s", response)
	}

	var searchResponse struct {
		Hits struct {
			Hits []struct {
				Id string `json:"_id"`
			} `json:"hits"`
		} `json:"hits"`
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка прочтения тела ответа %w", err)
	}

	err = json.Unmarshal(body, &searchResponse)
	if err != nil {
		return nil, fmt.Errorf("ошибка преобразования данных из json %w", err)
	}

	ids = make([]string, len(searchResponse.Hits.Hits))
	for i := 0; i < len(searchResponse.Hits.Hits); i++ {
		ids[i] = searchResponse.Hits.Hits[i].Id
	}

	return ids, nil
}
