package response

import "net/http"

var (
	StatusCodeUninitialized = -1
)

type ResponseWriter struct {
	http.ResponseWriter
	statuscode int
}

func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statuscode:     StatusCodeUninitialized,
	}
}

func (rw *ResponseWriter) WriteHeader(statuscode int) {
	rw.ResponseWriter.WriteHeader(statuscode)
	rw.statuscode = statuscode
}

func (rw *ResponseWriter) GetStatus() int {
	if rw.statuscode == StatusCodeUninitialized {
		panic("uninitialized status code")
	}
	return rw.statuscode
}

// SeenStatus reports whether WriteHeader was called on this wrapper and the status if so.
func (rw *ResponseWriter) SeenStatus() (code int, ok bool) {
	if rw.statuscode == StatusCodeUninitialized {
		return 0, false
	}
	return rw.statuscode, true
}
