package circuitbreaker

import (
	"errors"
	"log"
	"time"
)

type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

var ErrCircuitOpen = errors.New("circuit breaker is open")

// TODO: implement mutex
type CircuitBreaker struct {
	state            CircuitBreakerState
	failureCount     int
	failureThreshold int
	probeInFlight    bool
	cooldownPeriod   time.Duration
	now              func() time.Time
	openedAt         time.Time
}

func NewCircuitBreaker(failureThreshold int, cooldownPeriod time.Duration, now func() time.Time) *CircuitBreaker {
	return &CircuitBreaker{failureThreshold: failureThreshold, cooldownPeriod: cooldownPeriod, now: now}
}

// TODO: implemement mutex for any shared state access (state, fail and success count, etc)
func (cb *CircuitBreaker) Execute(fn func() error) error {
	if !cb.Allow() {
		return ErrCircuitOpen
	}

	err := fn()
	if err != nil {
		cb.RecordFailure()
	}

	return err
}

func (cb *CircuitBreaker) CurrentState() CircuitBreakerState {
	return cb.state
}

func (cb *CircuitBreaker) RecordFailure() {
	switch cb.state {
	case Closed:
		cb.failureCount++
		if cb.failureCount >= cb.failureThreshold {
			cb.state = Open
			cb.openedAt = cb.now()
		}
	case HalfOpen:
		cb.state = Open
		cb.openedAt = cb.now()
		cb.probeInFlight = false
	}
}

func (cb *CircuitBreaker) Allow() bool {
	switch cb.state {
	case Closed:
		return true
	case Open:
		if cb.now().Sub(cb.openedAt) < cb.cooldownPeriod {
			return false
		}

		cb.state = HalfOpen
		return true
	case HalfOpen:
		if cb.probeInFlight {
			return false
		}

		// allow only 1 probeInFlight
		cb.probeInFlight = true
		return true
	default:
		log.Printf("unknown circuit breaker state: %v", cb.state)
		return true
	}
}

func (cb *CircuitBreaker) GetFailureCount() int {
	return cb.failureCount
}

func (cb *CircuitBreaker) IsProbeInFlight() bool {
	return cb.probeInFlight
}
