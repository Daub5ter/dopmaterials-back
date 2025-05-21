package handlers

import (
	"io"
	log "log/slog"
	"net/http"
	"os"

	"file-server/internal/models/consts"
	"file-server/internal/models/errs"
	"file-server/internal/tools/server/servparams"
	"file-server/pkg/code"
)

func GetPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := servparams.ViewParam(r, consts.FileNameWithExtensionParam)
		if fileName == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		http.ServeFile(w, r, consts.StoragePhotosPath+fileName)
	}
}

func AddPhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := servparams.ViewParam(r, consts.FileNameParam)
		if fileName == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		// Ограничение в MB (n << 20) - максимальный размер n мегабайт.
		r.Body = http.MaxBytesReader(w, r.Body, consts.MaxPhotosSizeMB<<20)

		err := r.ParseMultipartForm(consts.MaxPhotosSizeMB << 20)
		if err != nil {
			if err.Error() == consts.HttpRequestBodyTooLarge {
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				return
			}

			log.Error("не получается прочитать форму", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		file, _, err := r.FormFile(consts.MultiFormPhotoKey)
		if err != nil {
			log.Error("не получается получить файл", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer func() { _ = file.Close() }()

		fileBody, err := io.ReadAll(file)
		if err != nil {
			log.Error("не получается прочитать файл", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		// Определение типа файла.
		var contentTypeBody []byte
		if len(fileBody) < 512 {
			contentTypeBody = fileBody
		} else {
			contentTypeBody = fileBody[:512]
		}

		contentType := http.DetectContentType(contentTypeBody)
		fileExtension, ok := consts.AllowedPhotoContentType[contentType]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		err = os.WriteFile(
			consts.StoragePhotosPath+fileName+fileExtension,
			fileBody,
			os.FileMode(0644),
		)
		if err != nil {
			log.Error("не получается записать файл", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func DeletePhoto() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := servparams.ViewParam(r, consts.FileNameWithExtensionParam)
		if fileName == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		err := os.Remove(consts.StoragePhotosPath + fileName)
		if err != nil {
			log.Error("не получается удалить файл", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
