package errors

import "errors"

var (
	ErrInvalidDateTime = errors.New("DateTime is invalid")
	InvalidSignature   = errors.New("invalid signature")
	TokenExpired       = errors.New("token is expired")
	InvalidToken       = errors.New("invalid token")
)
