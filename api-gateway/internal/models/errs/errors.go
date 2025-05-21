package errs

import (
	"errors"
)

// ошибки, которые обрабатываются в функциях.
var (
	ErrForbidden       = errors.New("нет доступа")
	ErrTooManyRequests = errors.New("слишком много запросов")
	ErrNotFound        = errors.New("ничего не найдено")
	ErrBadRequest      = errors.New("плохо сформулирован запрос")
	ErrCreds           = errors.New("ошибка сертификатов")
	ErrConnection      = errors.New("ошибка соединения с сервером")
	ErrUnauthorized    = errors.New("для выполнения этого действия нужно авторизоваться")
	ErrInternal        = errors.New("неизвестная ошибка")
)
