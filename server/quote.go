package server

import (
	"net/http"

	log "github.com/inconshreveable/log15"

	"github.com/dimfeld/brokerage_server/types"
)

func GetQuote(logger log.Logger, engine types.BrokerageServerPluginV1, w *ResponseWriter, r *http.Request, params map[string]string) {
	errorResponse(w, ERR_NOT_IMPLEMENTED, nil)
}
