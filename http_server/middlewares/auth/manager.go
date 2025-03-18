package auth

import (
	"http_server/usecases"
	"net/http"
	"strings"
)

const (
    sessionIdHeaderName = "Authorization"
)

type AuthMiddleware struct {
    sessMgr usecases.SessionManager
}

func NewObject(sessMgr usecases.SessionManager) *AuthMiddleware {
    return &AuthMiddleware{
        sessMgr: sessMgr,
    }
}

func GetSessionId(r *http.Request) (string, error) {
    value := strings.Split(r.Header.Get(sessionIdHeaderName), " ")
    
    if len(value) != 2 || value[0] != "Bearer" {
        return "", ErrorUnauthorized
    }

    return value[1], nil
}

func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if sessId, err := GetSessionId(r); err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        } else if _, err := am.sessMgr.GetSession(sessId); err != nil {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }

        next.ServeHTTP(w, r)
    })
}
