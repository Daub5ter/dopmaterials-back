package errs

import (
	"errors"
)

// ошибки, которые обрабатываются в функциях.
var (
	ErrNotFound   = errors.New("ничего не найдено")
	ErrBadRequest = errors.New("плохо сформулирован запрос")
	ErrInternal   = errors.New("неизвестная ошибка")
)
