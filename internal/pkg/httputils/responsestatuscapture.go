package httputils

import "net/http"

type HTTPResponseStatusWriter struct {
	http.ResponseWriter
	Status int
}

func (h *HTTPResponseStatusWriter) WriteHeader(status int) {
	h.Status = status
	h.ResponseWriter.WriteHeader(status)
}
