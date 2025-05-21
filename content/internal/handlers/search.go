package handlers

import (
	"content/internal/models"
	"content/internal/tools/server/servparams"
	"errors"
	log "log/slog"
	"net/http"
	"strconv"
	"strings"

	"content/internal/models/consts"
	"content/internal/models/errs"
	"content/pkg/code"
)

func SearchMaterials(s search) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := r.URL.Query()
		findPart := params.Get(consts.FindPartParam)

		offsetRaw := params.Get(consts.OffsetParam)
		var offset uint32

		if offsetRaw == "" {
			offset = 0
		} else {
			offsetNumRaw, err := strconv.Atoi(offsetRaw)
			if err != nil {
				code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
				return
			}

			offset = uint32(offsetNumRaw)
		}

		var categoryId *uint32
		categoryIdRaw := params.Get(consts.CategoryIdParam)
		if categoryIdRaw != "" {
			categoryIdNum, err := strconv.Atoi(categoryIdRaw)
			if err != nil {
				code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
				return
			}

			categoryIdNumUintRaw := uint32(categoryIdNum)
			categoryId = &categoryIdNumUintRaw
		}

		materialsIds, err := s.SearchMaterials(strings.ToLower(findPart), categoryId, offset)
		if err != nil {
			if errors.Is(err, errs.ErrBadRequest) {
				code.ErrorJSON(w, http.StatusBadRequest, err)
				return
			}
			if errors.Is(err, errs.ErrNotFound) {
				code.ErrorJSON(w, http.StatusNotFound, err)
				return
			}

			log.Error("Ошибка поиска материалов", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		code.WriteJSON(w, http.StatusOK, materialsIds)
	}
}

func PutMaterialSearch(s search) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := servparams.ViewParam(r, consts.IdParam)
		if id == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		var material models.MaterialSearch
		material.Id = id
		err := code.ReadJSON(r.Body, &material)
		if err != nil {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		err = s.PutMaterial(material)
		if err != nil {
			if errors.Is(err, errs.ErrBadRequest) {
				code.ErrorJSON(w, http.StatusBadRequest, err)
				return
			}
			if errors.Is(err, errs.ErrNotFound) {
				code.ErrorJSON(w, http.StatusNotFound, err)
				return
			}

			log.Error("Ошибка добавления материала для поиска", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteMaterialSearch(s search) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := servparams.ViewParam(r, consts.IdParam)
		if id == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		err := s.DeleteMaterial(id)
		if err != nil {
			if errors.Is(err, errs.ErrBadRequest) {
				code.ErrorJSON(w, http.StatusBadRequest, err)
				return
			}
			if errors.Is(err, errs.ErrNotFound) {
				code.ErrorJSON(w, http.StatusNotFound, err)
				return
			}

			log.Error("Ошибка добавления материала для поиска", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
