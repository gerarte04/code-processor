package writer

type ResponseWriter struct {
    WriteAddr string 
}

func NewResponseWriter() *ResponseWriter {
    return &ResponseWriter{
        WriteAddr: "http://localhost:8080/commit",
    }
}

func (w *ResponseWriter) WriteResponse(resp any) error {
    return nil
}
