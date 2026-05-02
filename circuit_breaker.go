package circuitbreaker

import (
	"log"
	"time"
)

type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

type CircuiBreaker interface {
	Execute() error
	CurrentState() CircuitBreakerState
	RecordFailure()
	RecordSuccess()
	Allow() bool
}

// TODO: implement mutex
type circuitBreaker struct {
	state         CircuitBreakerState
	failCount     int
	successCount  int
	threshold     int
	probeInFlight bool
	now           func() time.Time
	timeout       time.Duration
}

func NewCircuitBreaker(threshold int, now func() time.Time) CircuiBreaker {
	return &circuitBreaker{state: Closed, threshold: threshold, now: now}
}

// TODO: implemement mutex for any shared state access (state, fail and success count, etc)
func (cb *circuitBreaker) Execute() error { return nil }

func (cb *circuitBreaker) CurrentState() CircuitBreakerState {
	return cb.state
}

func (cb *circuitBreaker) RecordFailure() {
	cb.failCount++
}

func (cb *circuitBreaker) RecordSuccess() {
	cb.successCount++
}

func (cb *circuitBreaker) SetCurrentState(state CircuitBreakerState) {
	cb.state = state
}

func (cb *circuitBreaker) Allow() bool {
	switch cb.CurrentState() {
	case Closed:
		return true
	case Open:
		if time.Since(cb.now()) < cb.timeout {
			return false
		}
		cb.SetCurrentState(HalfOpen)
		return true
	case HalfOpen:
		if cb.probeInFlight {
			return false
		}

		// allow only 1 probeInFlight
		cb.probeInFlight = true
		return true
	default:
		log.Printf("unknown circuit breaker state: %v", cb.CurrentState())
		return true
	}
}
