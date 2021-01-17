package internal

import "errors"

var (
	ErrInvalidParam = errors.New("invalid param")

	ErrRunning    = errors.New("running error")
	ErrStopServer = errors.New("stop server")
)
