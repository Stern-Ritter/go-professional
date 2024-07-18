package hw05parallelexecution

import (
	"errors"
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		err := Run(tasks, workersCount, maxErrorsCount)

		require.Truef(t, errors.Is(err, ErrErrorsLimitExceeded), "actual err - %v", err)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		err := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.NoError(t, err)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("should complete all tasks when no errors received", func(t *testing.T) {
		counter := atomic.Int64{}
		successTask := func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			counter.Add(1)
			return nil
		}

		tasks := []Task{successTask, successTask, successTask}
		err := Run(tasks, 3, 1)

		expected := int64(len(tasks))
		got := counter.Load()
		require.NoError(t, err, "unexpected error when process tasks")
		assert.Equal(t, expected, got, "expect completed tasks count should be: %d, got: %d",
			expected, got)
	})

	t.Run("should complete all tasks when only one worker is used and no errors received", func(t *testing.T) {
		counter := atomic.Int64{}
		successTask := func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			counter.Add(1)
			return nil
		}

		tasks := []Task{successTask, successTask, successTask}
		err := Run(tasks, 1, 1)

		expected := int64(len(tasks))
		got := counter.Load()
		require.NoError(t, err, "unexpected error when process tasks")
		assert.Equal(t, expected, got, "expect completed tasks count should be: %d, got: %d",
			expected, got)
	})

	t.Run("should complete all tasks when error limit hasn`t been exceeded", func(t *testing.T) {
		counter := atomic.Int64{}
		successTask := func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			counter.Add(1)
			return nil
		}
		errorTask := func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			counter.Add(1)
			return errors.New("error completing task")
		}

		tasks := []Task{errorTask, errorTask, successTask}
		err := Run(tasks, 3, 2)

		expected := int64(len(tasks))
		got := counter.Load()
		require.NoError(t, err, "unexpected error when process tasks")
		assert.Equal(t, expected, got, "expect completed tasks count should be: %d, got: %d",
			expected, got)
	})

	t.Run("should return error when error limit has been exceeded", func(t *testing.T) {
		counter := atomic.Int64{}
		successTask := func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			counter.Add(1)
			return nil
		}
		errorTask := func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			counter.Add(1)
			return errors.New("error completing task")
		}

		tasks := []Task{errorTask, errorTask, successTask}
		err := Run(tasks, 3, 1)

		expected := int64(len(tasks))
		got := counter.Load()
		require.Error(t, err, "expect errors limit has been exceeded error")
		assert.ErrorIs(t, err, ErrErrorsLimitExceeded,
			"expect errors limit has been exceeded error: %v, got error: %v", ErrErrorsLimitExceeded, err)
		assert.LessOrEqual(t, got, expected, "expect completed tasks count should be less or equal to: %d, got: %d",
			expected, got)
	})

	t.Run("should complete all tasks when zero error limit has been used", func(t *testing.T) {
		counter := atomic.Int64{}
		successTask := func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			counter.Add(1)
			return nil
		}
		errorTask := func() error {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
			counter.Add(1)
			return errors.New("error completing task")
		}

		tasks := []Task{errorTask, errorTask, successTask}
		err := Run(tasks, 4, 0)

		expected := int64(len(tasks))
		got := counter.Load()
		require.NoError(t, err, "unexpected error when process tasks")
		assert.Equal(t, expected, got, "expect completed tasks count should be: %d, got: %d",
			expected, got)
	})

	t.Run("should correctly process empty tasks list", func(t *testing.T) {
		tasks := make([]Task, 0)
		err := Run(tasks, 3, 1)
		require.NoError(t, err, "unexpected error when process tasks: %v", err)
	})
}
