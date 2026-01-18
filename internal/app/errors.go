package app

import (
	"errors"

	"github.com/dannyvelas/conflux"
)

var (
	ErrInvalidArgs       = errors.New("invalid arguments")
	ErrInvalidFields     = conflux.ErrInvalidFields
	ErrHostAlreadyExists = errors.New("host already exists in ssh config file")
)
