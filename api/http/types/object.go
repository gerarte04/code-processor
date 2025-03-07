package types

import (
	"encoding/json"
	"http_server/repository"
	"http_server/repository/models"
	"http_server/usecases"
	"io"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
)

const (
    sessionIdHeaderName = "Authorization"
)

func GetSessionId(r *http.Request) (string, error) {
    value := strings.Split(r.Header.Get(sessionIdHeaderName), " ")
    
    if len(value) != 2 || value[0] != "Bearer" {
        return "", ErrorUnauthorized
    }

    return value[1], nil
}

func ProcessCreateError(w http.ResponseWriter, err error) error {
    if err == ErrorUnauthorized {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return err
    } else if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return err
    }

    return nil
}

type GetObjectHandlerRequest struct {
    Key uuid.UUID
    SessionId string
}

func CreateGetResultObjectHandlerRequest(r *http.Request) (*GetObjectHandlerRequest, error) {
    str := path.Base(r.URL.Path)
    key, err := uuid.Parse(str)

    if err != nil {
        return nil, ErrorInvalidKey
    }

    sessionId, err := GetSessionId(r)

    if err != nil {
        return nil, err
    }

    return &GetObjectHandlerRequest{Key: key, SessionId: sessionId}, nil
}

func CreateGetStatusObjectHandlerRequest(r *http.Request) (*GetObjectHandlerRequest, error) {
    str := path.Base(r.URL.Path)
    key, err := uuid.Parse(str)

    if err != nil {
        return nil, ErrorInvalidKey
    }

    sessionId, err := GetSessionId(r)

    if err != nil {
        return nil, err
    }

    return &GetObjectHandlerRequest{Key: key, SessionId: sessionId}, nil
}

type PostTaskObjectHandlerRequest struct {
    Dur time.Duration
    SessionId string
}

func CreatePostTaskObjectHandlerRequest(r *http.Request) (*PostTaskObjectHandlerRequest, error) {
	sessionId, err := GetSessionId(r)

	if err != nil {
		return nil, err
	}

    str, err := io.ReadAll(r.Body)

    if err != nil {
        return nil, err
    } else if len(str) == 0 {
        return &PostTaskObjectHandlerRequest{Dur: time.Second, SessionId: sessionId}, nil
    }

    mp := make(map[string]string)

    if err := json.Unmarshal(str, &mp); err != nil {
        return nil, err
    }

    d, ok := mp["duration"]

    if !ok || len(d) == 0 {
        return &PostTaskObjectHandlerRequest{Dur: time.Second, SessionId: sessionId}, nil
    } else if dur, err := time.ParseDuration(d); err != nil {
        return nil, err
    } else {
        return &PostTaskObjectHandlerRequest{Dur: dur, SessionId: sessionId}, nil
    }
}

type PostUserObjectHandlerRequest struct {
    Login string `json:"username"`
    Password string `json:"password"`
}

func CreatePostUserObjectHandlerRequest(r *http.Request) (*PostUserObjectHandlerRequest, error) {
    str, err := io.ReadAll(r.Body)

    if err != nil {
        return nil, err
    }

    var req PostUserObjectHandlerRequest

    if err = json.Unmarshal([]byte(str), &req); err != nil {
        return nil, err
    }

    return &req, nil
}

type GetResultObjectHandlerResponse struct {
    Result int `json:"result"`
}

type GetStatusObjectHandlerResponse struct {
    Status string `json:"status"`
}

type PostTaskObjectHandlerResponse struct {
    TaskId string `json:"task_id"`
}

type PostUserRegisterObjectHandlerResponse struct {}

type PostUserLoginObjectHandlerResponse struct {
    SessionId string `json:"token"`
}

func CreateGetResultObjectHandlerResponse(value *models.Task, err error) (*GetResultObjectHandlerResponse, error) {
    if err != nil {
        return nil, err
    } else if !value.Finished {
        return nil, usecases.ErrorTaskProcessing
    }

    return &GetResultObjectHandlerResponse{Result: value.Result}, nil
}

func CreateGetStatusObjectHandlerResponse(value *models.Task, err error) (*GetStatusObjectHandlerResponse, error) {
    if err != nil {
        return nil, err
    }

    if value.Finished {
        return &GetStatusObjectHandlerResponse{Status: "ready"}, nil
    } else {
        return &GetStatusObjectHandlerResponse{Status: "in_progress"}, nil
    }
}

func CreatePostTaskObjectHandlerResponse(value *uuid.UUID, err error) (*PostTaskObjectHandlerResponse, error) {
    if err != nil {
        return nil, err
    }

    return &PostTaskObjectHandlerResponse{TaskId: value.String()}, nil
}

func ProcessError(w http.ResponseWriter, err error, resp any) {
    if err == usecases.ErrorTaskProcessing {
        w.WriteHeader(http.StatusProcessing)
        w.Write([]byte(err.Error()))
        return
    } else if err == repository.ErrorTaskNotFound || err == ErrorNotFoundPath {
        http.Error(w, "Not Found", http.StatusNotFound)
        w.Write([]byte(err.Error()))
        return
    } else if err == ErrorUnauthorized || err == usecases.ErrorNoSessionExists || err == usecases.ErrorSessionExpired {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        w.Write([]byte(err.Error()))
        return
    } else if err == ErrorInvalidKey {
        http.Error(w, "Bad request", http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    } else if err != nil {
        http.Error(w, "Internal Error", http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    switch resp.(type) {
    case *PostTaskObjectHandlerResponse:
        w.WriteHeader(http.StatusCreated)
    }

    if resp != nil {
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            http.Error(w, "Internal Error", http.StatusInternalServerError)
        }
    }
}

func ProcessErrorPostUser(w http.ResponseWriter, err error, resp any) {
    if err == repository.ErrorUserAlreadyExists || err == repository.ErrorUserNotFound {
        http.Error(w, "Bad request", http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return
    } else if err == repository.ErrorWrongPassword || err == usecases.ErrorUserSessionExists {
        http.Error(w, "Forbidden", http.StatusForbidden)
        w.Write([]byte(err.Error()))
        return
    } else if err != nil {
        http.Error(w, "Internal Error", http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    switch resp.(type) {
    case *PostUserRegisterObjectHandlerResponse:
        w.WriteHeader(http.StatusCreated)
    case *PostUserLoginObjectHandlerResponse:
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            http.Error(w, "Internal Error", http.StatusInternalServerError)
        }
    }
}
