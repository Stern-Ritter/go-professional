package hw05parallelexecution

import (
	"errors"
	"fmt"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if len(tasks) == 0 {
		return nil
	}

	taskCh := make(chan Task)
	resCh := make(chan struct{}, len(tasks))
	errCh := make(chan error, m)
	interruptCh := make(chan struct{})

	go startTasksProcess(taskCh, interruptCh, tasks)

	numberOfWorkers := getMin(len(tasks), n)
	startTasksProcessWorkerPool(resCh, errCh, taskCh, interruptCh, numberOfWorkers)

	err := awaitTasksComplete(interruptCh, resCh, errCh, len(tasks), m)

	return err
}

func startTasksProcess(taskCh chan<- Task, interruptCh <-chan struct{}, tasks []Task) {
	for _, task := range tasks {
		select {
		case taskCh <- task:
		case <-interruptCh:
			close(taskCh)
			return
		}
	}
}

func startTasksProcessWorkerPool(resCh chan<- struct{}, errorCh chan<- error, taskCh <-chan Task,
	interruptCh <-chan struct{}, count int,
) {
	for i := 0; i < count; i++ {
		go func() {
			for {
				select {
				case task, ok := <-taskCh:
					if !ok {
						return
					}

					err := task()
					if err == nil {
						resCh <- struct{}{}
					} else {
						errorCh <- err
					}
				case <-interruptCh:
					return
				}
			}
		}()
	}
}

func awaitTasksComplete(interruptCh chan<- struct{}, resCh <-chan struct{}, errCh <-chan error, totalTasks int, limitErr int) error {
	defer close(interruptCh)

	needCheckErrLimit := limitErr > 0
	completedTasks := 0
	caughtErr := 0

	for {
		select {
		case <-resCh:
			completedTasks++

			if completedTasks >= totalTasks {
				return nil
			}
		case err := <-errCh:
			completedTasks++
			caughtErr++

			fmt.Printf("Error #%d: %s\n", caughtErr, err)

			if needCheckErrLimit && caughtErr > limitErr {
				return ErrErrorsLimitExceeded
			}
			if completedTasks >= totalTasks {
				return nil
			}
		}
	}
}

func getMin(f, s int) int {
	if f < s {
		return f
	}
	return s
}
