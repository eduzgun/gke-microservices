package errs

import "errors"

var (
	ErrInvalidCredentials = errors.New("invalid login credentials")
	ErrNotFound           = errors.New("not found")
)
