package taskpool

import (
	"time"
)

type Runnable func(chan interface{}, chan error)

type RunnableTask struct {
	function Runnable
	retChan  chan interface{}
	errChan  chan error
}

func NewRunnableTask(function Runnable) RunnableTask {
	return RunnableTask{function: function, retChan: make(chan interface{}, 1), errChan: make(chan error, 1)}
}

func (g RunnableTask) Get() (interface{}, error) {
	// first check for errors
	err := g.error()
	if err != nil {
		return nil, err
	}

	// no errors? try fetching the value
	select {
	case value, ok := <-g.retChan:
		if ok {
			return value, nil
		} else {
			return nil, ErrChanClosed
		}
	default:
		return nil, ErrInProgress
	}
}

func (g RunnableTask) Start() {
	g.function(g.retChan, g.errChan)
}

func (g RunnableTask) error() error {
	select {
	case value, ok := <-g.errChan:
		if ok {
			return value
		} else {
			return ErrChanClosed
		}
	default:
		return nil
	}
}

func (g RunnableTask) GetWait(duration time.Duration) (interface{}, error) {
	if duration.Seconds() == 0 {
		return g.Get()
	}
	maxWaitTime := time.Now().Add(duration)
	result, err := g.Get()
	for err == ErrInProgress && time.Now().Before(maxWaitTime) {
		result, err = g.Get()
	}

	return result, err
}
