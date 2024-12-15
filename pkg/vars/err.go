package vars

import "errors"

var (
	ErrAlreadyExist = errors.New("already exist")
	ErrUnauthorized = errors.New("unauthorized")
)
