package mock

import "bytes"

// ReadCloser is a mock io.ReadCloser.
type ReadCloser struct {
	// IsClosed indicates if the ReadCloser is closed.
	IsClosed bool
	*bytes.Buffer
}

// TODO: rename New to Make
// TODO: add interface assertion

// NewReadCloser creates a mock io.ReadCloser.
func NewReadCloser(b []byte) ReadCloser {
	return ReadCloser{
		Buffer: bytes.NewBuffer(b),
	}
}

// Close sets IsClosed to true.
func (rc *ReadCloser) Close() error {
	rc.IsClosed = true
	return nil
}
