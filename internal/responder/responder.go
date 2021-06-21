package responder

import (
	"encoding/json"
	"net/http"
)

// Responder is a JSend responder.
// Source: https://github.com/omniti-labs/jsend.
type Responder interface {
	// Respond sends the JSend response.
	Respond(Response) error
}

type Response struct {
	Status  Status `json:"status"`
	Data    []byte `json:"data"`
	Message string `json:"message"`
}

type Status string

const (
	// Success indicates that the API is successfully executed.
	Success Status = "success"
	// Fail indicates the the API call has invalid input.
	Fail = "fail"
	// Error indicates that the API errored while executing.
	Error = "error"
)

type responder struct {
	http.ResponseWriter
}

func New(w http.ResponseWriter) Responder {
	// TODO: replace the responder with some nice library
	return &responder{
		ResponseWriter: w,
	}
}

func (r responder) Respond(i Response) error {
	r.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(r).Encode(i)
}
