package server

import (
	"net/http"

	"github.com/json-iterator/go"

	"github.com/dimfeld/brokerage_server/types"
	"github.com/dimfeld/httptreemux"
	log "github.com/inconshreveable/log15"
)

func GetOptionAttributes(logger log.Logger, engine types.BrokerageServerPluginV1, w *ResponseWriter, r *http.Request, params map[string]string) {

	data, err := engine.GetOptions(r.Context(), params["symbol"])
	if err != nil {
		errorResponse(w, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = jsoniter.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data":   data,
	})

	if err != nil {
		logger.Error("Response encoding error", "err", err)
	}
}

func addOptionHandlers(router *httptreemux.TreeMux, Middleware MiddlewareFunc) {
	router.GET("/options/:symbol", Middleware(GetOptionAttributes))
}
