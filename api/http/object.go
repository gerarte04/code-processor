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

// @Description get result
// @Tags task
// @Produce json
// @Param task_id path string true "Task id"
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} types.GetResultObjectHandlerResponse
// @Success 102 {string} string "Processing"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Object not found"
// @Failure 500 {string} string "Internal error"
// @Router /result/{task_id} [get]
func (s *Object) getResultHandler(w http.ResponseWriter, r *http.Request) {
    req, err := types.CreateGetResultObjectHandlerRequest(r)
    if err = types.ProcessCreateError(w, err); err != nil {
        return
    }         

    value, err := s.service.GetTask(req.Key, req.SessionId)
    resp, err := types.CreateGetResultObjectHandlerResponse(value, err)

    types.ProcessError(w, err, resp)
}

// @Description get status
// @Tags task
// @Produce json
// @Param task_id path string true "Task id"
// @Param Authorization header string true "Authorization token"
// @Success 200 {object} types.GetStatusObjectHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Object not found"
// @Failure 500 {string} string "Internal error"
// @Router /status/{task_id} [get]
func (s *Object) getStatusHandler(w http.ResponseWriter, r *http.Request) {
    req, err := types.CreateGetStatusObjectHandlerRequest(r)
    if err = types.ProcessCreateError(w, err); err != nil {
        return
    }  

    value, err := s.service.GetTask(req.Key, req.SessionId)
    resp, err := types.CreateGetStatusObjectHandlerResponse(value, err)

    types.ProcessError(w, err, resp)
}

// @Description post task
// @Tags task
// @Accept  json
// @Produce json
// @Param duration body string true "Task duration"
// @Param Authorization header string true "Authorization token"
// @Success 201 {object} types.PostTaskObjectHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Object not found"
// @Failure 500 {string} string "Internal error"
// @Router /task [post]
func (s *Object) postTaskHandler(w http.ResponseWriter, r *http.Request) {
    req, err := types.CreatePostTaskObjectHandlerRequest(r)
    if err = types.ProcessCreateError(w, err); err != nil {
        return
    }  

    value, err := s.service.PostTask(req.Dur, req.SessionId)
    resp, err := types.CreatePostTaskObjectHandlerResponse(value, err)

    types.ProcessError(w, err, resp)
}

// @Description post register
// @Tags user
// @Accept  json
// @Produce json
// @Param credentials body types.PostUserObjectHandlerRequest true "login and password"
// @Success 201 {object} types.PostUserRegisterObjectHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Object not found"
// @Failure 500 {string} string "Internal error"
// @Router /register [post]
func (s *Object) postRegisterHandler(w http.ResponseWriter, r *http.Request) {
    req, err := types.CreatePostUserObjectHandlerRequest(r)
    if err = types.ProcessCreateError(w, err); err != nil {
        return
    }

    err = s.service.RegisterUser(req.Login, req.Password)
    types.ProcessErrorPostUser(w, err, &types.PostUserRegisterObjectHandlerResponse{})
}

// @Description post login
// @Tags user
// @Accept  json
// @Produce json
// @Param credentials body types.PostUserObjectHandlerRequest true "login and password"
// @Success 201 {object} types.PostUserLoginObjectHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Object not found"
// @Failure 500 {string} string "Internal error"
// @Router /login [post]
func (s *Object) postLoginHandler(w http.ResponseWriter, r *http.Request) {
    req, err := types.CreatePostUserObjectHandlerRequest(r)
    if err = types.ProcessCreateError(w, err); err != nil {
        return
    }

    value, err := s.service.LoginUser(req.Login, req.Password)
    types.ProcessErrorPostUser(w, err, &types.PostUserLoginObjectHandlerResponse{SessionId: value})
}

// WithObjectHandlers registers object-related HTTP handlers.
func (s *Object) WithObjectHandlers(r chi.Router) {
    r.Route("/", func(r chi.Router) {
        r.Get("/result/*", s.getResultHandler)
        r.Get("/status/*", s.getStatusHandler)
        r.Post("/task", s.postTaskHandler)
        r.Post("/register", s.postRegisterHandler)
        r.Post("/login", s.postLoginHandler)
    })
}
