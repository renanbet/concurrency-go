package worker

import (
	"context"
	"log"
	"sync"

	"concurrency.go/circuit-breaker/circuitbreaker"
	"concurrency.go/circuit-breaker/domain"
	"concurrency.go/circuit-breaker/ratelimit"
)

func Process(
	ctx context.Context,
	rl *ratelimit.RateLimiter,
	tasks <-chan int,
	results chan<- domain.Result,
	wg *sync.WaitGroup,
	processTask func(context.Context, int) domain.Result,
	cb *circuitbreaker.CircuitBreaker) {
	defer wg.Done()

	for {
		select {
		case task, ok := <-tasks:
			if !ok {
				return
			}

			if err := rl.Acquire(ctx); err != nil {
				return
			}

			res := domain.Result{}
			log.Println("Processing task:", task)
			if err := cb.Allow(); err != nil {
				res = domain.Result{Err: err}
			} else {
				res = processTask(ctx, task)
				if res.Err != nil {
					cb.Failure(res.Err)
				} else {
					cb.Success()
				}
			}

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
