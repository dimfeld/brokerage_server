package server

import (
	"net/http"
	"sync"

	log "github.com/inconshreveable/log15"
	jsoniter "github.com/json-iterator/go"

	"github.com/dimfeld/brokerage_server/types"
	"github.com/dimfeld/httptreemux"
)

type QuoteAndSymbol struct {
	*types.Quote
	Symbol string `json:"symbol"`
}

func GetQuote(logger log.Logger, engine types.BrokerageServerPluginV1, w *ResponseWriter, r *http.Request, params map[string]string) {
	data, err := engine.GetStockQuote(r.Context(), params["symbol"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse(w, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = jsoniter.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data": QuoteAndSymbol{
			data, params["symbol"],
		},
	})

	if err != nil {
		logger.Error("Response encoding error", "err", err)
	}
}

func GetQuotes(logger log.Logger, engine types.BrokerageServerPluginV1, w *ResponseWriter, r *http.Request, _ map[string]string) {
	body := struct {
		Symbols []string `json:"symbols"`
	}{}

	if err := jsoniter.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		errorResponse(w, err, nil)
	}

	results := make([]QuoteAndSymbol, len(body.Symbols))
	wg := sync.WaitGroup{}
	var symbolErr error
	ctx := r.Context()
	for i, symbol := range body.Symbols {
		wg.Add(1)
		go func(i int, symbol string) {
			quote, err := engine.GetStockQuote(ctx, symbol)
			if err != nil {
				symbolErr = err
			}
			results[i] = QuoteAndSymbol{
				quote,
				symbol,
			}
			wg.Done()
		}(i, symbol)
	}

	wg.Wait()

	if symbolErr != nil {
		errorResponse(w, symbolErr, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	err := jsoniter.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data":   results,
	})

	if err != nil {
		logger.Error("Response encoding error", "err", err)
	}
}

func addQuoteHandlers(router *httptreemux.TreeMux, Middleware MiddlewareFunc) {
	router.GET("/quotes/:symbol", Middleware(GetQuote))
	router.POST("/quotes", Middleware(GetQuotes))
}
