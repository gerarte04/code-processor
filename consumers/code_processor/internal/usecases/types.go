package usecases

import "time"

type ProcessingServiceResponse struct {
    Output string `json:"output"`
    StatusCode int64 `json:"status_code"`
    ProcessingTime time.Duration
    Translator string
}
