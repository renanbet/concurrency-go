package circuitbreaker

import (
	"errors"
	"log"
	"sync"
	"time"
)

type Status int

const (
	OPEN                    Status = 0
	CLOSED                  Status = 1
	HALFOPEN                Status = 2
	ErrCircuitOpen                 = "circuit breaker is open"
	ErrHalfOpenLimitReached        = "circuit breaker HALFOPEN limit reached"
)

type CircuitBreaker struct {
	mu               sync.Mutex
	maxFailures      int
	failures         int
	status           Status
	resetTimeout     time.Duration
	openUntil        time.Time
	halfOpenMaxCalls int
	halfOpenCalls    int
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
	now := time.Now()
	return &CircuitBreaker{
		mu:               sync.Mutex{},
		maxFailures:      maxFailures,
		failures:         0,
		status:           CLOSED,
		resetTimeout:     resetTimeout,
		openUntil:        now,
		halfOpenMaxCalls: 2,
		halfOpenCalls:    0,
	}
}

func (cb *CircuitBreaker) Allow() error {
	cb.mu.Lock()
	defer cb.mu.Unlock()

	if cb.status == OPEN {
		now := time.Now()
		if now.After(cb.openUntil) {
			cb.status = HALFOPEN
			cb.halfOpenCalls = 0
			log.Println("OPEN -> HALFOPEN")
			return nil
		}
		log.Println(ErrCircuitOpen)
		return errors.New(ErrCircuitOpen)
	}
	return nil
}

func (cb *CircuitBreaker) Success() {
	log.Println("Success")
	cb.mu.Lock()
	defer cb.mu.Unlock()
	if cb.status == HALFOPEN {
		cb.halfOpenCalls++
		log.Println("halfOpenCalls:", cb.halfOpenCalls)
		if cb.halfOpenCalls >= cb.halfOpenMaxCalls {
			cb.status = CLOSED
			cb.halfOpenCalls = 0
			cb.failures = 0
			log.Println("HALFOPEN -> CLOSED")
		}
	}
}

func (cb *CircuitBreaker) Failure(err error) {
	log.Println("Failure")
	cb.mu.Lock()
	defer cb.mu.Unlock()
	now := time.Now()
	cb.failures++
	log.Println("failures:", cb.failures)
	if cb.status == HALFOPEN {
		cb.status = OPEN
		log.Println("HALFOPEN -> OPEN")
		cb.openUntil = now.Add(cb.resetTimeout)
		return
	}
	if cb.failures >= cb.maxFailures {
		cb.openUntil = now.Add(cb.resetTimeout)
		log.Println("CLOSE -> OPEN")
		cb.status = OPEN
		cb.failures = 0
	}
}
