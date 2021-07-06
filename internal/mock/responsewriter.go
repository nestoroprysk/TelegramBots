package mock

import "net/http"

// ResponseWriter is a mock for testing responses.
type ResponseWriter struct {
	// Written lists write calls.
	Written [][]byte
	// StatusCode is HTTP status code.
	StatusCode int
	header     http.Header
}

var _ http.ResponseWriter = &ResponseWriter{}

// NewResponseWriter creates a mock response writer.
func NewResponseWriter() ResponseWriter {
	return ResponseWriter{
		header: http.Header(map[string][]string{}),
	}
}

// Header returns a  header that may be written to.
func (rw ResponseWriter) Header() http.Header {
	return rw.header
}

// Write adds a line to Written and returns the line len.
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.Written = append(rw.Written, b)
	return len(b), nil
}

// WriteHeader sets a status code
func (rw *ResponseWriter) WriteHeader(statusCode int) {
	rw.StatusCode = statusCode
}
