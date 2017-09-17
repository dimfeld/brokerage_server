package server

import (
	"encoding/json"
	"expvar"
	"net/http"
	"time"

	log "github.com/inconshreveable/log15"
	uuid "github.com/satori/go.uuid"

	"github.com/dimfeld/httptreemux"
)

type Handler func(logger log.Logger, engine interface{}, w *ResponseWriter, r *http.Request, params map[string]string)
type MiddlewareFactory func(handler Handler) httptreemux.HandlerFunc

var (
	ERR_NOT_IMPLEMENTED = Error{http.StatusNotImplemented, "Not implemented"}
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

func Start(ip string, port int, engine interface{}) { // TODO Proper type for broker engine
	log.Info("Starting server", "port", port)

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

	router.GET("/quote/:symbol", Middleware(GetQuote))

}

func HealthHandler(logger log.Logger, engine interface{}, w *ResponseWriter, r *http.Request, params map[string]string) {
	w.WriteHeader(http.StatusOK)
}