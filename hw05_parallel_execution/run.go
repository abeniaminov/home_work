package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var (
	ErrErrorsLimitExceeded   = errors.New("errors limit exceeded")
	ErrWorkersCountLessThen1 = errors.New("workers count less than 1")
)

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	var errors int32
	var wg sync.WaitGroup

	if n < 1 {
		return ErrWorkersCountLessThen1
	}

	worker := func(ch <-chan Task) {
		wg.Add(1)
		go func() {
			for task := range ch {
				if task() != nil {
						atomic.AddInt32(&errors, 1)
				}
			}
			wg.Done()
		}()
	}

	ch := make(chan Task)

	defer func() {
		close(ch)
		wg.Wait()
	}()

	for i := 0; i < n; i++ {
		worker(ch)
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errors) >= int32(m) {
			return ErrErrorsLimitExceeded
		} 
		ch <- task	
	}
	return nil
}
