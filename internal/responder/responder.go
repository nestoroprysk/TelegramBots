package responder

import (
	"encoding/json"
	"net/http"
)

// Responder is a JSend responder.
//
// Source: https://github.com/omniti-labs/jsend.
type Responder interface {
	// Succeed sets status to Success and Data to result.
	Succeed(interface{}) error
	// Fail sets status to Fail and Message to result (e.g., invalid input or precondition failed).
	Fail(error) error
	// Error sets status to Error and Message to result (e.g., coding or infra issue).
	Error(error) error
}

// Repsonse is a structure of every response.
//
// Succeed sets status to Success and Data to result.
// Fail sets status to Fail and Message to result (e.g., invalid input or precondition failed).
// Error sets status to Error and Message to result (e.g., coding or infra issue).
type Response struct {
	// Status is a response status.
	// May be success, fail, or error.
	Status Status `json:"status"`
	// Data is the success result.
	Data interface{} `json:"data"`
	// Message is the result for either fail or error.
	Message string `json:"message"`
	// StatusCode is HTTP status code.
	StatusCode int `json:"status_code"`
}

// Status is a response status.
type Status string

const (
	// Success indicates that the API is successfully executed.
	Success Status = "success"
	// Fail indicates the the API call has invalid input.
	Fail = "fail"
	// Error indicates that the API errored while executing.
	Error = "error"
)

var _ Responder = &responder{}

// responder implements the Responder interface.
type responder struct {
	http.ResponseWriter
}

// New creates a new responder.
func New(w http.ResponseWriter) Responder {
	// TODO: replace the responder with some nice library
	return &responder{
		ResponseWriter: w,
	}
}

// Succeed sets status to Success and Data to result.
// It returns http.StatusOK.
func (r responder) Succeed(b interface{}) error {
	return r.respond(Response{
		Status:     Success,
		Data:       b,
		StatusCode: http.StatusOK,
	})
}

// Fail sets status to Fail and Message to result (e.g., invalid input or precondition failed).
// It returns http.StatusBadRequest.
func (r responder) Fail(err error) error {
	return r.respond(Response{
		Status:     Fail,
		Message:    err.Error(),
		StatusCode: http.StatusBadRequest,
	})
}

// Error sets status to Error and Message to result (e.g., coding or infra issue).
// It returns http.StatusInternalServerError.
func (r responder) Error(err error) error {
	return r.respond(Response{
		Status:     Error,
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	})
}

// respond writes to the response writer using the JSON encoder.
func (r responder) respond(i Response) error {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(i.StatusCode)
	return json.NewEncoder(r).Encode(i)
}
