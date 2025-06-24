package usecase

import "errors"

var (
	ErrorGameNotStarted = errors.New("error game not started")
	ErrorWrongMoveCoordinates = errors.New("error wrong move coordinates")
)