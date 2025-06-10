package repository

import "errors"

var (
	ErrorWaitingNotFound = errors.New("waiting game not found")
	ErrorKeyNotFound = errors.New("key not found")
	ErrorAlreadyExists = errors.New("key already exists")
)