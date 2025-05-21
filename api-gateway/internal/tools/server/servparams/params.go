package servparams

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

// ViewParam получает парметр из uri запроса.
func ViewParam(r *http.Request, param string) string {
	return chi.URLParam(r, param)
}
