package http

import (
	"main/api/http/types"
	"main/usecases"
	"net/http"
	"path"

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

// @Summary Get an object
// @Description Get an object by its key
// @Tags object
// @Accept  json
// @Produce json
// @Param key query string true "Key of the object"
// @Success 200 {object} types.GetObjectHandlerResponse
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Object not found"
// @Router /object [get]
func (s *Object) getHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateGetObjectHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	value, err := s.service.Get(req.Key, path.Split(r.URL.Path))
	types.ProcessError(w, err, &types.GetObjectHandlerResponse{Value: value})
}

// @Summary Create or update an object
// @Description Create or update an object with the specified key and value
// @Tags object
// @Accept  json
// @Produce json
// @Param request body domain.Object true "Object data"
// @Success 200 {string} string "Object created or updated successfully"
// @Failure 400 {string} string "Bad request"
// @Router /object [put]
func (s *Object) putHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePutObjectHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = s.service.Put(req.Key, req.Value)
	types.ProcessError(w, err, nil)
}

// @Summary Create an object
// @Description Create a new object with the specified key and value
// @Tags object
// @Accept  json
// @Produce json
// @Param request body domain.Object true "Object data"
// @Success 200 {string} string "Object created successfully"
// @Failure 400 {string} string "Bad request"
// @Router /object [post]
func (s *Object) postHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreatePostObjectHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = s.service.Post(req.Key, req.Value)
	types.ProcessError(w, err, nil)
}

// @Summary Delete an object
// @Description Delete an object by its key
// @Tags object
// @Accept  json
// @Produce json
// @Param key query string true "Key of the object"
// @Success 200 {string} string "Object deleted successfully"
// @Failure 400 {string} string "Bad request"
// @Failure 404 {string} string "Object not found"
// @Router /object [delete]
func (s *Object) deleteHandler(w http.ResponseWriter, r *http.Request) {
	req, err := types.CreateDeleteObjectHandlerRequest(r)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	err = s.service.Delete(req.Key)
	types.ProcessError(w, err, nil)
}

// WithObjectHandlers registers object-related HTTP handlers.
func (s *Object) WithObjectHandlers(r chi.Router) {
	r.Route("/object", func(r chi.Router) {
		r.Get("/", s.getHandler)
		r.Post("/", s.postHandler)
		r.Put("/", s.putHandler)
		r.Delete("/", s.deleteHandler)
	})
}
