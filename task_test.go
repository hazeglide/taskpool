package taskpool

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTask_BasicUseCase(t *testing.T) {
	testString := "This is a Test"
	task := NewRunnableTask(func(outChan chan interface{}, errChan chan error) {
		outChan <- testString
	})
	result, err := task.Get()
	assert.Equal(t, ErrInProgress, err)
	task.Start()
	result, err = task.Get()
	assert.Equal(t, nil, err)
	assert.Equal(t, testString, result)
}

func TestTask_ErrorCase(t *testing.T) {
	testError := errors.New("TestError")
	task := NewRunnableTask(func(outChan chan interface{}, errChan chan error) {
		errChan <- testError
	})
	result, err := task.Get()
	assert.Equal(t, ErrInProgress, err)
	task.Start()
	result, err = task.Get()
	assert.Equal(t, testError, err)
	assert.Equal(t, nil, result)
}

func TestTask_ErrorChannelClosed(t *testing.T) {
	task := NewRunnableTask(func(outChan chan interface{}, errChan chan error) {
		close(errChan)
	})
	result, err := task.Get()
	assert.Equal(t, ErrInProgress, err)
	task.Start()
	result, err = task.Get()
	assert.Equal(t, ErrChanClosed, err)
	assert.Equal(t, nil, result)
}

func TestTask_OutChannelClosed(t *testing.T) {
	task := NewRunnableTask(func(outChan chan interface{}, errChan chan error) {
		close(outChan)
	})
	result, err := task.Get()
	assert.Equal(t, ErrInProgress, err)
	task.Start()
	result, err = task.Get()
	assert.Equal(t, ErrChanClosed, err)
	assert.Equal(t, nil, result)
}

func TestTask_ResultAndError(t *testing.T) {
	testError := errors.New("TestError")
	task := NewRunnableTask(func(outChan chan interface{}, errChan chan error) {
		outChan <- "This is a Test"
		errChan <- testError
	})
	result, err := task.Get()
	assert.Equal(t, ErrInProgress, err)
	task.Start()
	result, err = task.Get()
	assert.Equal(t, testError, err)
	assert.Equal(t, nil, result)
}

func TestTask_GetWaitWithZeroDuration(t *testing.T) {
	testError := errors.New("TestError")
	task := NewRunnableTask(func(outChan chan interface{}, errChan chan error) {
		outChan <- "This is a Test"
		errChan <- testError
	})
	result, err := task.Get()
	assert.Equal(t, ErrInProgress, err)
	task.Start()
	result, err = task.GetWait(0 * time.Second)
	assert.Equal(t, testError, err)
	assert.Equal(t, nil, result)
}

func TestTask_GetWait(t *testing.T) {
	testError := errors.New("TestError")
	task := NewRunnableTask(func(outChan chan interface{}, errChan chan error) {
		outChan <- "This is a Test"
		errChan <- testError
	})
	result, err := task.Get()
	assert.Equal(t, ErrInProgress, err)
	task.Start()
	result, err = task.GetWait(1 * time.Second)
	assert.Equal(t, testError, err)
	assert.Equal(t, nil, result)
}
