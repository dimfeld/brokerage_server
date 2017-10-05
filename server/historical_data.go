package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/dimfeld/httptreemux"
	log "github.com/inconshreveable/log15"
	"github.com/jinzhu/now"
	jsoniter "github.com/json-iterator/go"

	"github.com/dimfeld/brokerage_server/brokers"
	"github.com/dimfeld/brokerage_server/types"
)

type HistoricalDataResponse struct {
	Symbol string         `json:"symbol"`
	Type   string         `json:"type"`
	Bars   []*types.Quote `json:"bars"`
}

func GetHistoricalData(logger log.Logger, engines *brokers.EngineList, w *ResponseWriter, r *http.Request, params map[string]string) {

	symbol := params["symbol"]
	q := types.HistoricalDataParams{
		Symbol: symbol,
	}

	qs := r.URL.Query()
	if width, err := strconv.Atoi(qs.Get("barwidth")); err == nil {
		q.BarWidth = time.Duration(width) * time.Second
	} else {
		// Default to 1 day bar width.
		q.BarWidth = time.Duration(24) * time.Hour
	}

	if duration, err := strconv.Atoi(qs.Get("duration")); err == nil {
		q.Duration = time.Duration(duration) * time.Second
	} else {
		// Default to 1 day duration.
		q.Duration = time.Duration(24) * time.Hour
	}

	outputType := qs.Get("which")
	switch outputType {
	case "iv":
		q.Which = types.HistoricalDataTypeIv
	case "hv":
		q.Which = types.HistoricalDataTypeHv
	default:
		outputType = "price"
		q.Which = types.HistoricalDataTypePrice
	}

	if endTime, err := now.Parse(qs.Get("end")); err == nil {
		q.EndTime = endTime
	} else {
		q.EndTime = time.Now()
	}

	if includeAH, err := strconv.ParseBool(qs.Get("ah")); err == nil {
		q.IncludeAH = includeAH
	} else {
		q.IncludeAH = false
	}

	engine, err := engines.Get(brokers.PurposeEquityQuotes)
	if err != nil {
		errorResponse(w, err, nil)
		return
	}
	data, err := engine.GetHistoricalData(r.Context(), q)
	if err != nil {
		errorResponse(w, err, nil)
		return
	}

	w.WriteHeader(http.StatusOK)
	err = jsoniter.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data": HistoricalDataResponse{
			Symbol: symbol,
			Type:   outputType,
			Bars:   data,
		},
	})

	if err != nil {
		logger.Error("Response encoding error", "err", err)
	}
}

func addHistoricalDataHandlers(router *httptreemux.TreeMux, Middleware MiddlewareFunc) {
	router.GET("/quotes/:symbol/historical", Middleware(GetHistoricalData))
}
