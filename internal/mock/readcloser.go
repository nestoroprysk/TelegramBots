package mock

import (
	"bytes"
	"io"
)

// ReadCloser is a mock io.ReadCloser.
type ReadCloser struct {
	// IsClosed indicates if the ReadCloser is closed.
	IsClosed bool
	*bytes.Buffer
}

var _ io.ReadCloser = &ReadCloser{}

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
