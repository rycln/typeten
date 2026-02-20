package domain

import "errors"

// Domain errors for invalid state or invalid input.
// Callers can use errors.Is to branch on these.
var (
	ErrInvalidID        = errors.New("domain: invalid id")
	ErrInvalidUser      = errors.New("domain: invalid user")
	ErrInvalidTextInfo  = errors.New("domain: invalid text info")
	ErrInvalidFragment  = errors.New("domain: invalid text fragment")
	ErrInvalidSession   = errors.New("domain: invalid session")
	ErrInvalidSessionOp = errors.New("domain: invalid session operation")
)
