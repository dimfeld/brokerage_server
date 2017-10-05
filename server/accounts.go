package server

import (
	"net/http"
	"sync"

	log "github.com/inconshreveable/log15"
	jsoniter "github.com/json-iterator/go"

	"github.com/dimfeld/brokerage_server/brokers"
	"github.com/dimfeld/brokerage_server/types"
	"github.com/dimfeld/httptreemux"
)

func GetAccounts(logger log.Logger, engines *brokers.EngineList, w *ResponseWriter, r *http.Request, params map[string]string) {

	wg := &sync.WaitGroup{}

	brokerAccounts := make([][]*types.Account, len(engines.Engines))
	i := 0
	var outErr error
	for _, e := range engines.Engines {
		wg.Add(1)
		go func(e *brokers.BrokerEngine, index int) {
			if accounts, err := e.AccountList(r.Context()); err == nil {
				brokerAccounts[index] = accounts
			} else {
				outErr = err
			}

			wg.Done()
		}(e, i)
		i += 1
	}

	wg.Wait()

	accounts := make([]*types.Account, 0, len(brokerAccounts))
	for _, ac := range brokerAccounts {
		accounts = append(accounts, ac...)
	}

	if outErr != nil {
		errorResponse(w, outErr, nil)
		return
	}

	err := jsoniter.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data":   accounts,
	})

	if err != nil {
		logger.Error("Response encoding error", "err", err)
	}
}

func GetPositions(logger log.Logger, engines *brokers.EngineList, w *ResponseWriter, r *http.Request, params map[string]string) {

	wg := &sync.WaitGroup{}

	brokerPositions := make([][]*types.Position, len(engines.Engines))
	i := 0
	var outErr error
	for _, e := range engines.Engines {
		wg.Add(1)
		go func(e *brokers.BrokerEngine, index int) {
			if positions, err := e.GetPositions(r.Context()); err == nil {
				brokerPositions[index] = positions
			} else {
				outErr = err
			}

			wg.Done()
		}(e, i)
		i += 1
	}

	wg.Wait()

	totalLength := 0
	for _, p := range brokerPositions {
		totalLength += len(p)
	}

	positions := make([]*types.Position, 0, totalLength)
	for _, p := range brokerPositions {
		positions = append(positions, p...)
	}

	if outErr != nil {
		errorResponse(w, outErr, nil)
		return
	}

	err := jsoniter.NewEncoder(w).Encode(map[string]interface{}{
		"status": "ok",
		"data":   positions,
	})

	if err != nil {
		logger.Error("Response encoding error", "err", err)
	}
}

func addAccountHandlers(router *httptreemux.TreeMux, Middleware MiddlewareFunc) {
	router.GET("/accounts", Middleware(GetAccounts))
	router.GET("/positions", Middleware(GetPositions))
}
