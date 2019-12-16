package taskpool

import (
	"errors"
	"time"
)

var (
	ErrChanClosed = errors.New("channel closed")
	ErrInProgress = errors.New("task still in progress")
	ErrBufferFull = errors.New("could not put task, channel buffer is full")
)

type TaskPool interface {
	Start()
	Stop()
	Put(Task) error
}

type Task interface {
	Start()
	Get() (interface{}, error)
	GetWait(duration time.Duration) (interface{}, error)
}