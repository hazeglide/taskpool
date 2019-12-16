package taskpool

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPool_BasicPoolUseCase(t *testing.T) {
	pool := NewInMemoryPool(2, 2)
	pool.Start()
	defer pool.Stop()

	testString := "This is a Test"
	task := NewRunnableTask(func(c chan interface{}, errors chan error) {
		c <- testString
	})
	err := pool.Put(task)
	assert.Equal(t, nil, err)
	result, err := task.GetWait(1 * time.Second)

	assert.Equal(t, nil, err)
	assert.Equal(t, testString, result)
}

func TestPool_PutOnStoppedPool(t *testing.T) {
	pool := NewInMemoryPool(2, 2)
	pool.Start()
	pool.Stop()
	task := NewRunnableTask(func(c chan interface{}, errors chan error) {
		c <- "Test"
	})
	err := pool.Put(task)
	assert.Equal(t, nil, err)
	pool.Start()
	result, err := task.GetWait(1 * time.Second)

	assert.Equal(t, nil, err)
	assert.Equal(t, "Test", result)
}

func TestPool_PutOnFullBuffer(t *testing.T) {
	testString := "This is a Test"
	pool := NewInMemoryPool(2, 2)

	task := NewRunnableTask(func(c chan interface{}, errors chan error) {
		c <- testString
	})
	err := pool.Put(task)
	assert.Equal(t, nil, err)
	err = pool.Put(task)
	assert.Equal(t, nil, err)
	err = pool.Put(task)
	assert.Equal(t, ErrBufferFull, err)
	pool.Start()
	result, err := task.GetWait(1 * time.Second)

	assert.Equal(t, nil, err)
	assert.Equal(t, testString, result)
}

func BenchmarkInMemoryPool(b *testing.B) {
	pool := NewInMemoryPool(1, 1)
	pool.Start()
	defer pool.Stop()
	task := NewRunnableTask(func(outChan chan interface{}, errors chan error) {
		outChan <- true
	})
	for i := 0; i < b.N; i++ {
		_ = pool.Put(task)
	}
}
