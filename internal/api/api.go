package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"

	"github.com/alexadhy/wfreq/internal/logging"
	"github.com/alexadhy/wfreq/internal/store"
)

// API creates new API instance
type API struct {
	http.Handler
	l logging.Logger
	s *store.Store
}

// getLogger returns logger from HTTP context.
func getLogger(r *http.Request) logrus.FieldLogger {
	if log, ok := r.Context().Value(logging.LoggerKey).(logrus.FieldLogger); ok && log != nil {
		return log
	}

	return logging.New()
}

// New constructs a new API instance
func New(log logging.Logger) *API {
	a := &API{l: log, s: store.New(0)}
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Default().Handler)
	r.Use(logging.SetLoggerMiddleware(log))
	r.Post("/upload", a.handleUpload)
	r.Post("/", a.handleWordFrequencies)
	a.Handler = r
	return a
}

func (a *API) log(r *http.Request) logrus.FieldLogger {
	return getLogger(r)
}

// Error is the object returned to the client when there's an error.
type Error struct {
	Error string `json:"error"`
}

func (a *API) renderError(w http.ResponseWriter, r *http.Request, err error, status int) {
	if err == context.DeadlineExceeded {
		status = http.StatusRequestTimeout
	}

	// we fallback to 500
	if status == 0 {
		status = http.StatusInternalServerError
	}

	if status != http.StatusNotFound {
		a.log(r).Warnf("%d: %s", status, err)
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(&Error{Error: err.Error()}); err != nil {
		a.log(r).WithError(err).Warn("Failed to encode error")
	}
}
