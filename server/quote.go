package server

import (
	"net/http"

	log "github.com/inconshreveable/log15"
	jsoniter "github.com/json-iterator/go"

	"github.com/dimfeld/brokerage_server/types"
)

func GetQuote(logger log.Logger, engine types.BrokerageServerPluginV1, w *ResponseWriter, r *http.Request, params map[string]string) {
	data, err := engine.GetStockQuote(r.Context(), params["symbol"])
	if err != nil {
		errorResponse(w, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = jsoniter.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"symbol": params["symbol"],
		"data":   data,
	})

	if err != nil {
		logger.Error("Response encoding error", "err", err)
	}
}
