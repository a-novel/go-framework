package goframework

import "fmt"

var (
	// ErrInvalidEntity is thrown when the data passed to a model is not valid.
	ErrInvalidEntity      = fmt.Errorf("entity is not valid")
	ErrInvalidCredentials = fmt.Errorf("invalid credentials")
)
