package http

import (
	"http_server/middlewares"

	"github.com/go-chi/chi/v5"
)

const (
    rootPath = "/"
    healthCheckPath = "/health"

    getResultPath = "/result/*"
    getStatusPath = "/status/*"
    postTaskPath = "/task"

    postRegisterPath = "/register"
    postLoginPath = "/login"
)

// WithObjectHandlers registers object-related HTTP handlers.
func (s *Object) WithFreeUserHandlers(r chi.Router) func(r chi.Router) {
    return func(r chi.Router) {
        r.Post(postRegisterPath, s.postRegisterHandler)
        r.Post(postLoginPath, s.postLoginHandler)
        r.Get(healthCheckPath, s.healthCheckHandler)
    }
}

// Return function that sets group of handlers with authentication required
func (s *Object) WithSecuredUserHandlers(r chi.Router, authService middlewares.AuthMiddleware) func(r chi.Router) {
    return func(r chi.Router) {
        r.Group(func(r chi.Router) {
            r.Use(authService.Authenticate)
            r.Get(getResultPath, s.getResultHandler)
            r.Get(getStatusPath, s.getStatusHandler)
            r.Post(postTaskPath, s.postTaskHandler)
        })
    }
}

func (s *Object) RouteHandlers(r chi.Router, opts ...func(r chi.Router)) {
    r.Route(rootPath, func(r chi.Router) {
        for _, opt := range opts {
            opt(r)
        }
    })
}
