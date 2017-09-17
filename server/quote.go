package server

import (
	"net/http"

	log "github.com/inconshreveable/log15"
)

func GetQuote(logger log.Logger, engine interface{}, w *ResponseWriter, r *http.Request, params map[string]string) {
	errorResponse(w, ERR_NOT_IMPLEMENTED, nil)
}
