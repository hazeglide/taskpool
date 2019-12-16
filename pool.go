package taskpool

import (
	"sync"
)

// InMemoryPool implements a TaskPool backed by a buffered channel
type InMemoryPool struct {
	size        int
	taskChan    chan Task
	controlChan chan Control
	waitGroup   sync.WaitGroup
}

func NewInMemoryPool(poolSize int, bufferSize int) TaskPool {
	taskChan := make(chan Task, bufferSize)
	controlChan := make(chan Control, poolSize)
	return InMemoryPool{size: poolSize, taskChan: taskChan, controlChan: controlChan}
}

func (pool InMemoryPool) Start() {
	for i := 0; i < pool.size; i++ {
		startWorker(pool)
	}
}

func (pool InMemoryPool) Stop() {
	for i := 0; i < pool.size; i++ {
		pool.controlChan <- SignalShutdown
	}
	pool.waitGroup.Wait()
}

func (pool InMemoryPool) Put(task Task) error {
	select {
	case pool.taskChan <- task:
		return nil
	default:
		return ErrBufferFull
	}
}

func startWorker(pool InMemoryPool) {
	pool.waitGroup.Add(1)
	go func(waitGroup sync.WaitGroup, taskChan chan Task, controlChan chan Control) {
		select {
		case task, ok := <-taskChan:
			if ok {
				task.Start()
			}
		case control, ok := <-controlChan:
			if ok {
				if control == SignalShutdown {
					waitGroup.Done()
					return
				}
			}
		}
	}(pool.waitGroup, pool.taskChan, pool.controlChan)
}
