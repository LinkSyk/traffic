package internal

import "errors"

var (
	ErrInvalidParam = errors.New("invalid param")

	ErrRunning    = errors.New("running error")
	ErrStopServer = errors.New("stop server")

	ErrNoAvailableNode = errors.New("no available node")
	ErrSocketRead      = errors.New("read socket error")
	ErrSocketWrite     = errors.New("write socket error")
)
