package handlers

import (
	"errors"
	"io"
	log "log/slog"
	"net/http"
	"os"
	"os/exec"

	"file-server/internal/models/consts"
	"file-server/internal/models/errs"
	"file-server/internal/tools/server/servparams"
	"file-server/pkg/code"
)

func GetVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := servparams.ViewParam(r, consts.FileNameParam)
		if fileName == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		filePart := servparams.ViewParam(r, consts.FileNamePartParam)
		if filePart == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		http.ServeFile(w, r, consts.StorageVideosPath+fileName+"/"+filePart)
	}
}

func AddVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := servparams.ViewParam(r, consts.FileNameParam)
		if fileName == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		// Ограничение в MB (n << 20) - максимальный размер n мегабайт.
		r.Body = http.MaxBytesReader(w, r.Body, consts.MaxVideoSizeMB<<20)

		err := r.ParseMultipartForm(consts.MaxVideoSizeMB << 20)
		if err != nil {
			if err.Error() == consts.HttpRequestBodyTooLarge {
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				return
			}

			log.Error("не получается прочитать форму", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		file, _, err := r.FormFile(consts.MultiFormVideoKey)
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
		fileExtension, ok := consts.AllowedVideoContentType[contentType]
		if !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		dir := consts.StorageVideosPath + fileName
		err = os.Mkdir(
			dir,
			os.FileMode(0755),
		)
		if err != nil {
			if errors.Is(err, os.ErrExist) {
				code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			}

			log.Error("не получается создать директорию", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		filePath := dir + "/" + fileName + fileExtension
		err = os.WriteFile(
			filePath,
			fileBody,
			os.FileMode(0644),
		)
		if err != nil {
			log.Error("не получается записать файл", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)

			errRemoveDir := os.RemoveAll(dir)
			if errRemoveDir != nil {
				log.Error("не получается удалить директорию", log.Any("error", errRemoveDir))
			}

			return
		}

		cmd := exec.Command(
			"ffmpeg",
			"-i", filePath, // Входной файл
			"-codec:", "copy", // Копируем кодеки
			"-start_number", "0", // Начальный номер сегмента
			"-hls_time", "10", // Длительность каждого сегмента
			"-hls_list_size", "0", // Полный плейлист
			"-f", "hls", // Формат - HLS
			dir+"/"+fileName+".m3u8", // Путь к выходному файлу .m3u8
		)

		err = cmd.Run()
		if err != nil {
			log.Error("не получается обработать файл", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)

			errRemoveDir := os.RemoveAll(dir)
			if errRemoveDir != nil {
				log.Error("не получается удалить директорию", log.Any("error", errRemoveDir))
			}

			return
		}

		err = os.Remove(filePath)
		if err != nil {
			log.Error("не получается удалить изначальный файл", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)

			errRemoveDir := os.RemoveAll(dir)
			if errRemoveDir != nil {
				log.Error("не получается удалить директорию", log.Any("error", errRemoveDir))
			}

			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func DeleteVideo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fileName := servparams.ViewParam(r, consts.FileNameParam)
		if fileName == "" {
			code.ErrorJSON(w, http.StatusBadRequest, errs.ErrBadRequest)
			return
		}

		err := os.RemoveAll(consts.StorageVideosPath + "/" + fileName)
		if err != nil {
			log.Error("не получается удалить файл", log.Any("error", err))
			w.WriteHeader(http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
