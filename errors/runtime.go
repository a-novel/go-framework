package errors

import "fmt"

var (
	// ErrTimeout is thrown on any time out.
	ErrTimeout = fmt.Errorf("connection timed out")
)
