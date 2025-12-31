package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"concurrency.go/circuit-breaker/circuitbreaker"
	"concurrency.go/circuit-breaker/domain"
	"concurrency.go/circuit-breaker/ratelimit"
	"concurrency.go/circuit-breaker/worker"
	"github.com/cenkalti/backoff"
)

func handler(ctx context.Context, task int) domain.Result {
	workTime := time.Duration(7 * time.Second)
	timer := time.NewTimer(workTime)
	defer timer.Stop()

	if task%5 == 0 {
		log.Println(task, "is cancelled immediately")
		<-ctx.Done()
		return domain.Result{Err: ctx.Err()}
	}

	select {
	case <-timer.C:
		if task%7 == 0 {
			err := fmt.Errorf("task %d failed", task)
			return domain.Result{
				Err: err,
			}
		}
		value := task * 2
		return domain.Result{
			Value: value,
		}
	case <-ctx.Done():
		return domain.Result{
			Err: ctx.Err(),
		}
	}
}

func processTask(ctx context.Context, task int) domain.Result {
	expBackoff := backoff.NewExponentialBackOff()
	expBackoff.InitialInterval = 100 * time.Millisecond
	expBackoff.MaxInterval = 300 * time.Millisecond
	maxRetry := backoff.WithMaxRetries(expBackoff, uint64(3))
	result := domain.Result{}

	bo := backoff.WithContext(maxRetry, ctx)

	err := backoff.Retry(func() error {
		taskCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		result = handler(taskCtx, task)
		if result.Err != nil {
			if errors.Is(result.Err, context.DeadlineExceeded) {
				return result.Err
			}
			return backoff.Permanent(result.Err)
		}
		return nil
	}, bo)
	if err != nil {
		return domain.Result{Err: fmt.Errorf("task %d failed after retries: %w", task, err)}
	}

	return result
}

func main() {
	tasks := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 14, 15, 16, 21, 22, 23, 28, 29, 30, 35}
	log.Println("Starting processing tasks...", len(tasks))
	numWorkers := 5
	rateLimit := 3

	taskChan := make(chan int)
	resultChan := make(chan domain.Result)
	var wg sync.WaitGroup

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	rl := ratelimit.NewRateLimiter(rateLimit, time.Second)
	cb := circuitbreaker.NewCircuitBreaker(2, 5*time.Second)

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker.Process(ctx, rl, taskChan, resultChan, &wg, processTask, cb)
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
				log.Println(result.Err)
				if len(errors) >= 15 {
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
	fmt.Println("Errors:", len(errors))
	for _, err := range errors {
		fmt.Println(err)
	}
}
