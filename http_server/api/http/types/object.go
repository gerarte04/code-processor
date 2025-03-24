package types

import (
	"encoding/json"
	"http_server/repository"
	"http_server/repository/models"
	"http_server/usecases"
	"io"
	"net/http"
	"path"

	"github.com/google/uuid"
)

func ProcessCreateError(w http.ResponseWriter, err error) error {
    if err != nil {
        http.Error(w, "Bad request", http.StatusBadRequest)
        w.Write([]byte(err.Error()))
        return err
    }

    return nil
}

type GetObjectHandlerRequest struct {
    Key uuid.UUID
}

func CreateGetResultObjectHandlerRequest(r *http.Request) (*GetObjectHandlerRequest, error) {
    str := path.Base(r.URL.Path)
    key, err := uuid.Parse(str)

    if err != nil {
        return nil, ErrorInvalidKey
    }

    return &GetObjectHandlerRequest{Key: key}, nil
}

func CreateGetStatusObjectHandlerRequest(r *http.Request) (*GetObjectHandlerRequest, error) {
    str := path.Base(r.URL.Path)
    key, err := uuid.Parse(str)

    if err != nil {
        return nil, ErrorInvalidKey
    }

    return &GetObjectHandlerRequest{Key: key}, nil
}

type PostTaskObjectHandlerRequest struct {
    Translator string `json:"translator"`
    Code string `json:"code"`
}

type PutCommitObjectHandlerRequest struct {
    TaskId string `json:"task_id"`
    Output string `json:"output"`
    StatusCode int64 `json:"status_code"`
}

func CreatePostTaskObjectHandlerRequest(r *http.Request) (*PostTaskObjectHandlerRequest, error) {
    str, err := io.ReadAll(r.Body)

    if err != nil {
        return nil, err
    }

    var req PostTaskObjectHandlerRequest

    if err = json.Unmarshal([]byte(str), &req); err != nil {
        return nil, err
    }

    return &req, nil
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

func CreatePutCommitObjectHandlerRequest(r *http.Request) (*PutCommitObjectHandlerRequest, error) {
    str, err := io.ReadAll(r.Body)

    if err != nil {
        return nil, err
    }

    var req PutCommitObjectHandlerRequest

    if err = json.Unmarshal([]byte(str), &req); err != nil {
        return nil, err
    }

    return &req, nil
}

type GetResultObjectHandlerResponse struct {
    Output string `json:"result"`
    // StatusCode int64 `json:"status_code"`
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

    return &GetResultObjectHandlerResponse{
        Output: value.Output,
        //StatusCode: value.Result.StatusCode,
    }, nil
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
    } else if err == repository.ErrorWrongUserCreds || err == usecases.ErrorUserSessionExists {
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
