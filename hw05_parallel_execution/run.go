package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	tasksCh := make(chan Task)
	errsCh := make(chan error)
	doneCh := make(chan struct{})
	var wg sync.WaitGroup
	errs := make([]error, 0, m)

	go addTasks(tasksCh, doneCh, tasks)
	startWorkerPool(tasksCh, errsCh, doneCh, &wg, n)

	useErrorsLimit := m > 0
	go func() {
		for err := range errsCh {
			errs = append(errs, err)
			if useErrorsLimit && len(errs) >= m {
				close(doneCh)
				return
			}
		}
	}()

	wg.Wait()
	close(errsCh)

	if useErrorsLimit && len(errs) >= m {
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
