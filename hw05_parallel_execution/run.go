package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	errsLimit := int64(m)
	useErrsLimit := errsLimit > 0
	tasksCh := make(chan Task)
	errsCh := make(chan error)
	doneCh := make(chan struct{})
	var wg sync.WaitGroup
	var errsCount atomic.Int64

	go addTasks(tasksCh, doneCh, tasks)
	startWorkerPool(tasksCh, errsCh, doneCh, &wg, n)

	go func() {
		for range errsCh {
			errsCount.Add(1)
			if useErrsLimit && errsCount.Load() >= errsLimit {
				close(doneCh)
				return
			}
		}
	}()

	wg.Wait()
	close(errsCh)

	if useErrsLimit && errsCount.Load() >= errsLimit {
		return ErrErrorsLimitExceeded
	}
	return nil
}

func addTasks(tasksCh chan<- Task, doneCh <-chan struct{}, tasks []Task) {
	defer close(tasksCh)

	for _, task := range tasks {
		select {
		case tasksCh <- task:
		case <-doneCh:
			return
		}
	}
}

func startWorkerPool(tasksCh <-chan Task, errsCh chan<- error, doneCh <-chan struct{}, wg *sync.WaitGroup,
	workerCount int,
) {
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go worker(tasksCh, errsCh, doneCh, wg)
	}
}

func worker(tasksCh <-chan Task, errsCh chan<- error, doneCh <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-tasksCh:
			if !ok {
				return
			}
			if err := task(); err != nil {
				select {
				case errsCh <- err:
				case <-doneCh:
					return
				}
			}
		case <-doneCh:
			return
		}
	}
}
