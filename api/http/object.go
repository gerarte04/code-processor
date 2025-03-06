package http

import (
	"http_server/api/http/types"
	"http_server/usecases"
	"net/http"

	"github.com/go-chi/chi/v5"
)

// Object represents an HTTP handler for managing objects.
type Object struct {
    service usecases.Object
}

// NewHandler creates a new instance of Object.
func NewHandler(service usecases.Object) *Object {
    return &Object{service: service}
}

func (s *Object) getResultHandler(w http.ResponseWriter, r *http.Request) {
    req, err := types.CreateGetResultObjectHandlerRequest(r)

    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    value, err := s.service.GetTask(req.Key)
    resp, err := types.CreateGetResultObjectHandlerResponse(value, err)

    types.ProcessError(w, err, resp)
}

func (s *Object) getStatusHandler(w http.ResponseWriter, r *http.Request) {
    req, err := types.CreateGetStatusObjectHandlerRequest(r)

    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }

    value, err := s.service.GetTask(req.Key)
    resp, err := types.CreateGetStatusObjectHandlerResponse(value, err)

    types.ProcessError(w, err, resp)
}

func (s *Object) postHandler(w http.ResponseWriter, r *http.Request) {
    req, err := types.CreatePostObjectHandlerRequest(r)
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    }
    value, err := s.service.PostTask(req.Dur)
    resp, err := types.CreatePostObjectHandlerResponse(value, err)

    types.ProcessError(w, err, resp)
}

// WithObjectHandlers registers object-related HTTP handlers.
func (s *Object) WithObjectHandlers(r chi.Router) {
    r.Route("/", func(r chi.Router) {
        r.Get("/result/*", s.getResultHandler)
        r.Get("/status/*", s.getStatusHandler)
        r.Post("/task", s.postHandler)
    })
}
