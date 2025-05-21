package handlers

import (
	"errors"
	log "log/slog"
	"net/http"
	"strconv"
	"strings"

	"content/internal/models"
	"content/internal/models/consts"
	"content/internal/models/errs"
	"content/internal/tools/server/servparams"
	"content/pkg/code"
)

func GetMaterials(db database, s search) http.HandlerFunc {
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

		if findPart == "" {
			limit := uint32(consts.LimitMaterials)
			materials, err := db.GetMaterials(categoryId, &limit, offset)
			if err != nil {
				log.Error("Ошибка получения материалов", log.Any("error", err))
				code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
				return
			}

			code.WriteJSON(w, http.StatusOK, materials)
			return
		}

		materialsIds, err := s.SearchMaterials(strings.ToLower(findPart), categoryId, offset)
		if err != nil {
			log.Error("Ошибка поиска материалов", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		if len(materialsIds) == 0 {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		materials, err := db.GetSearchedMaterials(materialsIds)
		if err != nil {
			log.Error("Ошибка получения материалов после поиска", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		code.WriteJSON(w, http.StatusOK, materials)
		return
	}
}

func GetMaterial(db database) http.HandlerFunc {
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

		material, err := db.GetMaterial(uint64(id))
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				code.ErrorJSON(w, http.StatusNotFound, err)
				return
			}

			log.Error("Ошибка получения материала", log.Int("id", id), log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		code.WriteJSON(w, http.StatusOK, material)
	}
}

func AddMaterial(db database, s search) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var material models.Material

		err := code.ReadJSON(r.Body, &material)
		if err != nil {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		if material.Name == "" || material.Description == "" ||
			material.PreviewMeta == "" || material.VideoMeta == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		operationTransaction, id, err := db.InsertMaterial(&material)
		if err != nil {
			if errors.Is(err, errs.ErrBadRequest) {
				code.ErrorJSON(w, http.StatusBadRequest, err)
				return
			}

			log.Error("Ошибка добавления материала",
				log.Any("material", material), log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		err = s.PutMaterial(
			models.MaterialSearch{
				Id:          strconv.Itoa(int(id)),
				CategoryId:  material.CategoryId,
				Name:        material.Name,
				Description: material.Description,
			})
		if err != nil {
			log.Error("Ошибка добавления материала в поисковой движок", log.Any("error", err))
			if err = db.CancelOperation(operationTransaction); err != nil {
				log.Error("Ошибка отмены операции", log.Any("error", err))
			}
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		if err = db.ConfirmOperation(operationTransaction); err != nil {
			log.Error("Ошибка подтверждения операции", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		code.WriteJSON(w, http.StatusCreated, id)
	}
}

func UpdateMaterial(db database, s search) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var material models.Material

		err := code.ReadJSON(r.Body, &material)
		if err != nil {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		if material.Name == "" || material.Description == "" ||
			material.PreviewMeta == "" || material.VideoMeta == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		operationTransaction, err := db.UpdateMaterial(&material)
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				code.ErrorJSON(w, http.StatusNotFound, err)
				return
			}
			if errors.Is(err, errs.ErrBadRequest) {
				code.ErrorJSON(w, http.StatusBadRequest, err)
				return
			}

			log.Error("Ошибка изменения материала",
				log.Any("material", material), log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		err = s.PutMaterial(
			models.MaterialSearch{
				Id:          strconv.Itoa(int(material.Id)),
				CategoryId:  material.CategoryId,
				Name:        material.Name,
				Description: material.Description,
			})
		if err != nil {
			log.Error("Ошибка обновления материала в поисковом движке", log.Any("error", err))
			if err = db.CancelOperation(operationTransaction); err != nil {
				log.Error("Ошибка отмены операции", log.Any("error", err))
			}
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		if err = db.ConfirmOperation(operationTransaction); err != nil {
			log.Error("Ошибка подтверждения операции", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteMaterial(db database, s search) http.HandlerFunc {
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

		operationTransaction, err := db.DeleteMaterial(uint64(id))
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				code.ErrorJSON(w, http.StatusNotFound, err)
				return
			}

			log.Error("Ошибка удаления материала",
				log.Int("id", id), log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		err = s.DeleteMaterial(idRaw)
		if err != nil {
			log.Error("Ошибка удаления материала из поискового движка", log.Any("error", err))
			if err = db.CancelOperation(operationTransaction); err != nil {
				log.Error("Ошибка отмены операции", log.Any("error", err))
			}
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		if err = db.ConfirmOperation(operationTransaction); err != nil {
			log.Error("Ошибка подтверждения операции", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func RestoreMaterial(db database, s search) http.HandlerFunc {
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

		operationTransaction, err := db.RestoreMaterial(uint64(id))
		if err != nil {
			if errors.Is(err, errs.ErrNotFound) {
				code.ErrorJSON(w, http.StatusNotFound, err)
				return
			}

			log.Error("Ошибка восстановления материала",
				log.Int("id", id), log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		material, err := db.GetMaterial(uint64(id))
		if err != nil {
			log.Error("Ошибка получения материала", log.Int("id", id), log.Any("error", err))
			if err = db.CancelOperation(operationTransaction); err != nil {
				log.Error("Ошибка отмены операции", log.Any("error", err))
			}
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		err = s.PutMaterial(
			models.MaterialSearch{
				Id:          strconv.Itoa(id),
				CategoryId:  material.CategoryId,
				Name:        material.Name,
				Description: material.Description,
			})
		if err != nil {
			log.Error("Ошибка добавления материала в поисковой движок", log.Any("error", err))
			if err = db.CancelOperation(operationTransaction); err != nil {
				log.Error("Ошибка отмены операции", log.Any("error", err))
			}
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		if err = db.ConfirmOperation(operationTransaction); err != nil {
			log.Error("Ошибка подтверждения операции", log.Any("error", err))
			code.ErrorJSON(w, http.StatusInternalServerError, errs.ErrInternal)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
