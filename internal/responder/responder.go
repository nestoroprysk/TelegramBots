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
	Succeed(string) error
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
	Data string `json:"data"`
	// Message is the result for either fail or error.
	Message string `json:"message"`
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
func (r responder) Succeed(b string) error {
	return r.respond(Response{
		Status: Success,
		Data:   b,
	})
}

// Fail sets status to Fail and Message to result (e.g., invalid input or precondition failed).
func (r responder) Fail(err error) error {
	return r.respond(Response{
		Status:  Fail,
		Message: err.Error(),
	})
}

// Error sets status to Error and Message to result (e.g., coding or infra issue).
func (r responder) Error(err error) error {
	return r.respond(Response{
		Status:  Error,
		Message: err.Error(),
	})
}

// respond writes to the response writer using the JSON encoder.
func (r responder) respond(i Response) error {
	r.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(r).Encode(i)
}
