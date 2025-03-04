package types

import (
	"encoding/json"
	"fmt"
	"main/domain"
	"main/repository"
	"net/http"
)

type GetObjectHandlerRequest struct {
	Key string `json:"key"`
}

func CreateGetObjectHandlerRequest(r *http.Request) (*GetObjectHandlerRequest, error) {
	key := r.URL.Query().Get("key")
	if key == "" {
		return nil, fmt.Errorf("missing key")
	}
	return &GetObjectHandlerRequest{Key: key}, nil
}

type PutObjectHandlerRequest struct {
	domain.Object
}

func CreatePutObjectHandlerRequest(r *http.Request) (*PutObjectHandlerRequest, error) {
	var req PutObjectHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

type DeleteObjectHandlerRequest struct {
	Key string `json:"key"`
}

func CreateDeleteObjectHandlerRequest(r *http.Request) (*DeleteObjectHandlerRequest, error) {
	key := r.URL.Query().Get("key")
	if key == "" {
		return nil, fmt.Errorf("missing key")
	}
	return &DeleteObjectHandlerRequest{Key: key}, nil
}

type PostObjectHandlerRequest struct {
	domain.Object
}

func CreatePostObjectHandlerRequest(r *http.Request) (*PostObjectHandlerRequest, error) {
	var req PostObjectHandlerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, fmt.Errorf("error while decoding json: %v", err)
	}
	return &req, nil
}

type GetObjectHandlerResponse struct {
	Value *string `json:"value"`
}

func ProcessError(w http.ResponseWriter, err error, resp any) {
	if err == repository.NotFound {
		http.Error(w, "Key not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Internal Error", http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	if resp != nil {
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			http.Error(w, "Internal Error", http.StatusInternalServerError)
		}
	}
}