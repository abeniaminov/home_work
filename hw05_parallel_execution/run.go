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

	worker := func(ch chan Task, done <-chan struct{}) {
		go func() {
			for {
				select {
				case <-done:
					return
				case task := <-ch:
					errT := task()
					if errT != nil {
						atomic.AddInt32(&errors, 1)
					}
					wg.Done()
				}
			}
		}()
	}

	done := make(chan struct{})
	ch := make(chan Task)

	defer func() {
		wg.Wait()
		for i := 0; i < n; i++ {
			done <- struct{}{}
		}
	}()

	for i := 0; i < n; i++ {
		worker(ch, done)
	}

	for _, task := range tasks {
		if atomic.LoadInt32(&errors) < int32(m) {
			wg.Add(1)
			ch <- task
		} else {
			return ErrErrorsLimitExceeded
		}
	}
	return nil
}
