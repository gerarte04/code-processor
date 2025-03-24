package usecases

type ProcessingServiceResponse struct {
    Output string `json:"output"`
    StatusCode int64 `json:"status_code"`
}
