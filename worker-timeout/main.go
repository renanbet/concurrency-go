package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Result struct {
	Value int
	Err   error
}

func handler(ctx context.Context, task int) Result {
	workTime := time.Duration(rand.Intn(400)+100) * time.Millisecond
	timer := time.NewTimer(workTime)
	defer timer.Stop()

	select {
	case <-timer.C:
		if task%7 == 0 {
			err := fmt.Errorf("task %d failed", task)
			return Result{
				Err: err,
			}
		}
		value := task * task
		return Result{
			Value: value,
		}
	case <-ctx.Done():
		return Result{
			Err: ctx.Err(),
		}
	}
}

func processTask(ctx context.Context, task int) Result {
	taskCtx, cancel := context.WithTimeout(ctx, 250*time.Millisecond)
	defer cancel()

	return handler(taskCtx, task)
}

func worker(ctx context.Context, tasks <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			res := processTask(ctx, task)
			select {
			case results <- res:
			case <-ctx.Done():
				return
			}
		case <-ctx.Done():
			return
		}
	}
}

func main() {
	tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 14, 15, 16, 21, 22, 23, 28, 29}
	numWorkers := 3

	taskChan := make(chan int)
	resultChan := make(chan Result)
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(ctx, taskChan, resultChan, &wg)
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	// Send tasks to the worker
	go func() {
		for _, task := range tasks {
			select {
			case taskChan <- task:
			case <-ctx.Done():
				close(taskChan)
				return
			}
		}
		close(taskChan)
	}()

	var results []int
	var errors []error

	// Collect results and errors
	for resultChan != nil {
		result, ok := <-resultChan
		if ok {
			if result.Err != nil {
				errors = append(errors, result.Err)
				if len(errors) >= 20 {
					fmt.Println("Too many errors, cancelling...")
					cancel()
				}
			} else {
				results = append(results, result.Value)
			}
		} else {
			resultChan = nil
		}
	}

	fmt.Println("Results:", len(results), results)
	fmt.Println("Errors:", len(errors), errors)
}
