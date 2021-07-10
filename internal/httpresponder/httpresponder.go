package httpresponder

import (
	"encoding/json"
	"net/http"

	"github.com/nestoroprysk/TelegramBots/internal/responder"
)

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

// httpresponder implements the Responder interface.
type httpresponder struct {
	responder.Responder
	http.ResponseWriter
}

var _ responder.Responder = &httpresponder{}

// Wrap covers a responder with the HTTPresponder.
func Wrap(r responder.Responder, w http.ResponseWriter) responder.Responder {
	// TODO: replace the responder with some nice library
	return &httpresponder{
		Responder:      r,
		ResponseWriter: w,
	}
}

// Succeed sets status to Success and Data to result.
// It returns http.StatusOK.
func (r httpresponder) Succeed(b interface{}) error {
	if err := r.respond(Response{
		Status:     Success,
		Data:       b,
		StatusCode: http.StatusOK,
	}); err != nil {
		return r.Responder.Error(err)
	}

	return r.Responder.Succeed(b)
}

// Fail sets status to Fail and Message to result (e.g., invalid input or precondition failed).
// It returns http.StatusBadRequest.
func (r httpresponder) Fail(err error) error {
	if err := r.respond(Response{
		Status:     Fail,
		Message:    err.Error(),
		StatusCode: http.StatusBadRequest,
	}); err != nil {
		return r.Responder.Error(err)
	}

	return r.Responder.Fail(err)
}

// Error sets status to Error and Message to result (e.g., coding or infra issue).
// It returns http.StatusInternalServerError.
func (r httpresponder) Error(err error) error {
	if err := r.respond(Response{
		Status:     Error,
		Message:    err.Error(),
		StatusCode: http.StatusInternalServerError,
	}); err != nil {
		return r.Responder.Error(err)
	}

	return r.Responder.Error(err)
}

// respond writes to the response writer using the JSON encoder.
func (r httpresponder) respond(i Response) error {
	r.Header().Set("Content-Type", "application/json")
	r.WriteHeader(i.StatusCode)
	return json.NewEncoder(r).Encode(i)
}
