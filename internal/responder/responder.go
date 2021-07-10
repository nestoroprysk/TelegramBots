package responder

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
	// Close closes the responder.
	Close() error
}
