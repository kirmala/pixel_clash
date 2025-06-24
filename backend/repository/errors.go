package repository

import "errors"

var (
	ErrorWaitingNotFound = errors.New("error waiting game not found")
	ErrorKeyNotFound = errors.New("error key not found")
	ErrorAlreadyExists = errors.New("error key already exists")
	ErrorGameAlreadyStarted = errors.New("error game already started")
)