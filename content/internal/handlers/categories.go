package handlers

import (
	"content/internal/models"
	"content/internal/models/consts"
	"content/internal/models/errs"
	"content/internal/tools/server/servparams"
	"content/pkg/code"
	"errors"
	log "log/slog"
	"net/http"
	"strconv"
)

func GetMainCategories(db database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := db.GetNullParentCategories()
		if err != nil {
			log.Error("Ошибка получения основных категорий", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		code.WriteJSON(w, http.StatusOK, categories)
	}
}

func GetSubsidiariesCategories(db database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := servparams.ViewParam(r, consts.IdParam)
		if idRaw == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		id, err := strconv.Atoi(idRaw)
		if err != nil || id < 0 {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		categories, err := db.GetSubsidiariesCategories(uint32(id))
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				code.ErrorJSON(w, http.StatusNotFound, err)
				return
			}

			log.Error("Ошибка получения дочерних категорий",
				log.Int("id", id), log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		code.WriteJSON(w, http.StatusOK, categories)
	}
}

func GetCategory(db database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := servparams.ViewParam(r, consts.IdParam)
		if idRaw == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		id, err := strconv.Atoi(idRaw)
		if err != nil || id < 0 {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		categories, err := db.GetCategory(uint32(id))
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				code.ErrorJSON(w, http.StatusNotFound, err)
				return
			}

			log.Error("Ошибка получения категории",
				log.Int("categoryId", id), log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		code.WriteJSON(w, http.StatusOK, categories)
	}
}

func AddCategory(db database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var category models.Category

		err := code.ReadJSON(r.Body, &category)
		if err != nil {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		if category.Name == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		id, err := db.InsertCategory(&category)
		if err != nil {
			if errors.Is(err, errs.ErrBadRequest) {
				code.ErrorJSON(w, http.StatusBadRequest, err)
				return
			}

			log.Error("Ошибка добавления категории",
				log.Any("category", category), log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		code.WriteJSON(w, http.StatusCreated, id)
	}
}

func UpdateCategory(db database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var category models.Category

		err := code.ReadJSON(r.Body, &category)
		if err != nil {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		if category.Name == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		err = db.UpdateCategory(&category)
		if err != nil {
			if errors.Is(err, errs.ErrBadRequest) {
				code.ErrorJSON(w, http.StatusBadRequest, err)
				return
			}

			log.Error("Ошибка обновления категории",
				log.Any("category", category), log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteCategory(db database, s search) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idRaw := servparams.ViewParam(r, consts.IdParam)
		if idRaw == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		id, err := strconv.Atoi(idRaw)
		if err != nil || id < 0 {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		categoryId := uint32(id)
		materials, err := db.GetMaterials(&categoryId, nil, 0)
		if err != nil {
			log.Error("Ошибка получения материалов", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		operationTransaction, err := db.DeleteCategory(uint32(id))
		if err != nil {
			if errors.Is(err, errs.ErrBadRequest) {
				code.ErrorJSON(w, http.StatusBadRequest, err)
				return
			}

			log.Error("Ошибка удаления категории", log.Any("id", id), log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		for _, material := range materials {
			err = s.PutMaterial(
				models.MaterialSearch{
					Id:          strconv.Itoa(int(material.Id)),
					Name:        material.Name,
					Description: material.Description,
				})
			if err != nil {
				log.Error("Ошибка удаления материала из поискового движка", log.Any("error", err))
				if err = db.CancelOperation(operationTransaction); err != nil {
					log.Error("Ошибка отмены операции", log.Any("error", err))
				}
				code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
				return
			}
		}

		if err = db.ConfirmOperation(operationTransaction); err != nil {
			log.Error("Ошибка подтверждения операции", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
