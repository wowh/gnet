package gnet

import "github.com/pkg/errors"

var (
	ErrServerAlreadyStart         = errors.New("tcp server already started")
	ErrConnListenerNotRegister    = errors.New("conn listener not register")
	ErrTcpConnectionAlreadyClosed = errors.New("tcp connection already closed")
)
