package worker

import (
	"context"
	"sync"

	"concurrency.go/rate-limit/domain"
	"concurrency.go/rate-limit/ratelimit"
)

func Process(
	ctx context.Context,
	rl *ratelimit.RateLimiter,
	tasks <-chan int,
	results chan<- domain.Result,
	wg *sync.WaitGroup,
	handler func(context.Context, int) domain.Result) {
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

			res := handler(ctx, task)
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
