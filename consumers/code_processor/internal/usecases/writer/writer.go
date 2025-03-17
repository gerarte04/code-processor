package writer

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

type ResponseWriter struct {
    WriteAddr string 
    cli *http.Client
}

func NewResponseWriter() *ResponseWriter {
    return &ResponseWriter{
        WriteAddr: "http://http_server:8080/commit",
        cli: &http.Client{},
    }
}

func (w *ResponseWriter) WriteResponse(resp any) error {
    m, err := json.Marshal(resp)
    if err != nil {
        return err
    }

    req, err := http.NewRequest("PUT", w.WriteAddr, strings.NewReader(string(m)))
    if err != nil {
        return err
    }

    reqResp, err := w.cli.Do(req)
    if err != nil {
        return err
    }

    message, err := io.ReadAll(reqResp.Body)
    if err != nil {
        return err
    }

    log.Printf("http response code: %d, message: %s", reqResp.StatusCode, string(message))

    return nil
}
