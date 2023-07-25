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

	worker := func(cc chan chan Task, done <-chan struct{}) chan Task {
		workerChan := make(chan Task)
		go func() {
			defer close(workerChan)
			for {
				select {
				case <-done:
					return
				case task := <-workerChan:
					errT := task()
					if errT != nil {
						atomic.AddInt32(&errors, 1)
					}
					cc <- workerChan
					wg.Done()
				}
			}
		}()
		return workerChan
	}

	done := make(chan struct{})
	chChans := make(chan chan Task, n)

	defer func() {
		wg.Wait()
		for i := 0; i < n; i++ {
			done <- struct{}{}
		}
		close(chChans)
	}()

	for i := 0; i < n; i++ {
		chChans <- worker(chChans, done)
	}

	for _, task := range tasks {
		if n < 1 {
			return ErrWorkersCountLessThen1
		}
		if atomic.LoadInt32(&errors) < int32(m) {
			wg.Add(1)
			ch := <-chChans
			ch <- task
		} else {
			return ErrErrorsLimitExceeded
		}
	}
	return nil
}
