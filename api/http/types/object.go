package types

import (
	"encoding/json"
	"http_server/repository"
	"http_server/repository/models"
	"http_server/usecases"
	"io"
	"net/http"
	"path"
	"time"

	"github.com/google/uuid"
)

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

type PostObjectHandlerRequest struct {
    Dur time.Duration
}

func CreatePostObjectHandlerRequest(r *http.Request) (*PostObjectHandlerRequest, error) {
    str, _ := io.ReadAll(r.Body)

    if len(str) == 0 {
        return &PostObjectHandlerRequest{Dur: time.Second}, nil
    }

    dur, err := time.ParseDuration(string(str))

    if err != nil {
        return nil, err
    }

    return &PostObjectHandlerRequest{Dur: dur}, nil
}

type GetResultObjectHandlerResponse struct {
    Result int `json:"result"`
}

type GetStatusObjectHandlerResponse struct {
    Status string `json:"status"`
}

type PostObjectHandlerResponse struct {
    TaskId string `json:"task_id"`
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

func CreatePostObjectHandlerResponse(value *uuid.UUID, err error) (*PostObjectHandlerResponse, error) {
    if err != nil {
        return nil, err
    }

    return &PostObjectHandlerResponse{TaskId: value.String()}, nil
}

func ProcessError(w http.ResponseWriter, err error, resp any) {
    if err == usecases.ErrorTaskProcessing {
        http.Error(w, err.Error(), http.StatusProcessing)
        return
    } else if err == repository.ErrorTaskNotFound {
        http.Error(w, "Key not found", http.StatusNotFound)
    } else if err != nil {
        http.Error(w, "Internal Error", http.StatusInternalServerError)
        w.Write([]byte(err.Error()))
        return
    }

    switch resp.(type) {
    case *PostObjectHandlerResponse:
        w.WriteHeader(http.StatusCreated)
    }

    if resp != nil {
        if err := json.NewEncoder(w).Encode(resp); err != nil {
            http.Error(w, "Internal Error", http.StatusInternalServerError)
        }
    }
}