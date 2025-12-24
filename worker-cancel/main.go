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

func worker(ctx context.Context, tasks <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}
			time.Sleep(time.Duration(rand.Intn(400)+100) * time.Millisecond) // Simulate processing time
			if task%7 == 0 {
				err := fmt.Errorf("task %d failed", task)
				select {
				case results <- Result{Err: err}:
				case <-ctx.Done():
					return
				}
			} else {
				value := task * task
				select {
				case results <- Result{Value: value}:
				case <-ctx.Done():
					return
				}
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
				if len(errors) >= 10 {
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

	fmt.Println("Results:", results)
	fmt.Println("Errors:", errors)
}
