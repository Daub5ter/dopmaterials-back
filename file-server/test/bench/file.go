package bench

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

var file = "test.mp4"
var fileName = "test"
var dir = "testdir"

func fileAddFromFile() error {
	fileBody, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = os.Mkdir(
		dir,
		os.FileMode(0755),
	)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return err
		}

		return fmt.Errorf("не получается создать директорию %w", err)
	}

	filePath := dir + "/" + fileName + ".mp4"
	err = os.WriteFile(
		filePath,
		fileBody,
		os.FileMode(0644),
	)
	if err != nil {
		err = fmt.Errorf("не получается записать файл %w", err)

		errRemoveDir := os.RemoveAll(dir)
		if errRemoveDir != nil {
			err = fmt.Errorf("не получается удалить директорию, %v", errRemoveDir)
		}

		return err
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
		err = fmt.Errorf("не получается обработать файл,%w", err)

		errRemoveDir := os.RemoveAll(dir)
		if errRemoveDir != nil {
			err = fmt.Errorf("не получается удалить директорию, %w", errRemoveDir)
		}

		return err
	}

	err = os.Remove(filePath)
	if err != nil {
		err = fmt.Errorf("не получается удалить изначальный файл %w", err)

		errRemoveDir := os.RemoveAll(dir)
		if errRemoveDir != nil {
			err = fmt.Errorf("не получается удалить директорию, %v", errRemoveDir)
		}

		return err
	}

	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("не получается удалить директорию, %w", err)
	}

	return nil
}

func fileAddFromPipe() error {
	fileBody, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	err = os.Mkdir(dir, os.FileMode(0755))
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return err
		}
		return fmt.Errorf("не получается создать директорию %w", err)
	}

	// Создаем pipe
	pr, pw := io.Pipe()
	defer pw.Close()

	cmd := exec.Command(
		"ffmpeg",
		"-i", "pipe:0", // Входные данные из stdin
		"-codec:", "copy",
		"-start_number", "0",
		"-hls_time", "10",
		"-hls_list_size", "0",
		"-f", "hls",
		dir+"/"+fileName+".m3u8", // Путь к выходному файлу .m3u8
	)

	cmd.Stdin = pr // Устанавливаем stdin для команды

	go func() {
		defer pw.Close()
		pw.Write(fileBody) // Записываем содержимое файла в pipe
	}()

	err = cmd.Run()
	if err != nil {
		err = fmt.Errorf("не получается обработать файл: %w", err)
		errRemoveDir := os.RemoveAll(dir)
		if errRemoveDir != nil {
			err = fmt.Errorf("не получается удалить директорию: %w", errRemoveDir)
		}
		return err
	}

	err = os.RemoveAll(dir)
	if err != nil {
		return fmt.Errorf("не получается удалить директорию: %w", err)
	}

	return nil
}
