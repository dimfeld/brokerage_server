package server

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/json-iterator/go"

	"github.com/dimfeld/brokerage_server/types"
	"github.com/dimfeld/httptreemux"
	log "github.com/inconshreveable/log15"
)

var (
	ErrBadMinStrike    = Error{code: http.StatusBadRequest, message: "min_strike must be a number"}
	ErrBadMaxStrike    = Error{code: http.StatusBadRequest, message: "max_strike must be a number"}
	ErrNoStrikesFound  = Error{code: http.StatusBadRequest, message: "no strikes found in requested range"}
	ErrNoStrikeArgs    = Error{code: http.StatusBadRequest, message: "Must specify one or more `strike` arguments, or `min_strike` and `max_strike`"}
	ErrNoExpiryArgs    = Error{code: http.StatusBadRequest, message: "Must specify one or more `expiry` arguments, or `min_expiry` and `max_expiry`"}
	ErrNoExpiriesFound = Error{code: http.StatusBadRequest, message: "no expiries found in requested range"}
)

func GetOptionsAttributes(logger log.Logger, engine types.BrokerageServerPluginV1, w *ResponseWriter, r *http.Request, params map[string]string) {

	data, err := engine.GetOptionsChain(r.Context(), params["symbol"])
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

func GetOptionsQuote(logger log.Logger, engine types.BrokerageServerPluginV1, w *ResponseWriter, r *http.Request, params map[string]string) {

	ctx := r.Context()
	symbol := params["symbol"]
	meta, err := engine.GetOptionsChain(ctx, symbol)
	if err != nil {
		errorResponse(w, err, nil)
		return
	}

	query := r.URL.Query()
	minStrike := query.Get("min_strike")
	maxStrike := query.Get("max_strike")
	strikeParams := query["strike"]

	var strikes []float64

	if minStrike != "" && maxStrike != "" {

		fMinStrike, err := strconv.ParseFloat(minStrike, 64)
		if err != nil {
			errorResponse(w, ErrBadMinStrike, nil)
			return
		}

		fMaxStrike, err := strconv.ParseFloat(maxStrike, 64)
		if err != nil {
			errorResponse(w, ErrBadMaxStrike, nil)
			return
		}

		// There aren't enough strikes to be worth doing a binary search.
		for _, value := range meta.Strikes {
			if value > fMaxStrike {
				break
			}

			if value >= fMinStrike {
				strikes = append(strikes, value)
			}
		}
	} else if len(strikeParams) != 0 {

		strikes = make([]float64, len(strikeParams))
		for i, strike := range strikeParams {
			fStrike, err := strconv.ParseFloat(strike, 64)
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				errorResponse(w, fmt.Errorf("Invalid strike %s", strike), nil)
				return
			}
			strikes[i] = fStrike
		}
	} else {
		errorResponse(w, ErrNoStrikeArgs, nil)
		return
	}

	if len(strikes) == 0 {
		errorResponse(w, ErrNoStrikesFound, nil)
		return
	}

	minExpiry := query.Get("min_expiry")
	maxExpiry := query.Get("max_expiry")
	expiryParams := query["expiry"]

	var expirations []string
	if minExpiry != "" && maxExpiry != "" {
		for _, value := range meta.Expirations {
			if value > maxExpiry {
				break
			}

			if value >= minExpiry {
				expirations = append(expirations, value)
			}
		}

	} else if len(expiryParams) > 0 {
		validExpiries := map[string]bool{}
		for _, val := range meta.Expirations {
			validExpiries[val] = true
		}

		for _, e := range expiryParams {
			if validExpiries[e] {
				expirations = append(expirations, e)
			}
		}

	} else {
		errorResponse(w, ErrNoExpiryArgs, nil)
		return
	}

	if len(expirations) == 0 {
		errorResponse(w, ErrNoExpiriesFound, nil)
		return
	}

	requestParams := types.OptionsQuoteParams{
		Underlying:  symbol,
		Expirations: expirations,
		Strikes:     strikes,
		// Always get both puts and calls.
		Puts:  true,
		Calls: true,
	}

	data, err := engine.GetOptionsQuotes(ctx, requestParams)
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
	router.GET("/options/:symbol/meta", Middleware(GetOptionsAttributes))
	router.GET("/options/:symbol/quotes", Middleware(GetOptionsQuote))
}
