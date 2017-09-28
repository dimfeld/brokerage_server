package server

import (
	"context"
	"encoding/json"
	"expvar"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/inconshreveable/log15"
	uuid "github.com/satori/go.uuid"

	"github.com/dimfeld/httptreemux"

	"github.com/dimfeld/brokerage_server/brokers"
)

type Handler func(logger log.Logger, engines *brokers.EngineList, w *ResponseWriter, r *http.Request, params map[string]string)
type MiddlewareFactory func(handler Handler) httptreemux.HandlerFunc

var (
	ErrNotImplemented = Error{http.StatusNotImplemented, "Not implemented"}
)

type HttpError interface {
	Code() int
	Error() string
}

type Error struct {
	code    int
	message string
}

func (e Error) Error() string { return e.message }
func (e Error) Code() int     { return e.code }

// ResponseWriter is a simple wrapper around the http package's interface of the same name
// that allows tracking of whether or not a status code was written.
type ResponseWriter struct {
	http.ResponseWriter
	WroteStatus bool
	StatusCode  int
	Error       error
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		WroteStatus:    false,
		StatusCode:     200,
	}
}

func (w *ResponseWriter) GetStatus() int {
	return w.StatusCode
}

func (w *ResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.StatusCode = code
	w.WroteStatus = true
}

func errorResponse(w *ResponseWriter, err error, meta map[string]interface{}) {
	if meta == nil {
		meta = make(map[string]interface{})
	}

	w.Error = err

	if !w.WroteStatus {
		if httpError, ok := err.(HttpError); ok {
			w.WriteHeader(httpError.Code())
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	}

	meta["status"] = "error"
	meta["message"] = err.Error()
	json.NewEncoder(w).Encode(meta)
}

func expvarWrapper(w http.ResponseWriter, r *http.Request, params map[string]string) {
	expvar.Handler().ServeHTTP(w, r)
}

type MiddlewareFunc func(handler Handler) httptreemux.HandlerFunc

func Run(bind string, engine *brokers.EngineList) {
	log.Info("Starting server", "port", bind)

	var Middleware = func(handler Handler) httptreemux.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request, params map[string]string) {
			reqID := r.Header.Get("x-request-id")
			if reqID == "" {
				reqID = uuid.NewV4().String()
			}

			logger := log.New("req_id", reqID)
			writer := NewResponseWriter(w)
			start := time.Now()

			w.Header().Set("content-type", "application/json")

			handler(logger, engine, writer, r, params)
			duration := time.Now().Sub(start) / time.Millisecond
			status := writer.GetStatus()
			if status == 500 {
				logger.Error("internal error", "err", writer.Error)
			}
			logger.Info(r.URL.String(), "method", r.Method, "responseTime", int64(duration), "statusCode", status)

		}
	}

	router := httptreemux.New()
	router.PanicHandler = httptreemux.ShowErrorsJsonPanicHandler

	router.GET("/healthz", Middleware(HealthHandler))
	router.GET("/debug/vars", expvarWrapper)

	addQuoteHandlers(router, Middleware)
	addOptionHandlers(router, Middleware)

	server := &http.Server{
		Addr:    bind,
		Handler: router,
	}

	go func() {
		server.ListenAndServe()
	}()

	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)
	<-stopChan

	log.Warn("Shutting down...")
	shutdownCtx, cancelFunc := context.WithTimeout(context.Background(), time.Duration(10*time.Second))
	defer cancelFunc()
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		log.Error("Error shutting down!", "err", err.Error())
	}
}

func HealthHandler(logger log.Logger, engine *brokers.EngineList, w *ResponseWriter, r *http.Request, params map[string]string) {
	w.WriteHeader(http.StatusOK)
}
