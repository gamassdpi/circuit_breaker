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
	state            CircuitBreakerState
	failureCount     int
	successCount     int
	failureThreshold int
	probeInFlight    bool
	cooldownPeriod   time.Duration
	now              func() time.Time
	openedAt         time.Time
}

func NewCircuitBreaker(failureThreshold int, cooldownPeriod time.Duration, now func() time.Time) CircuiBreaker {
	return &circuitBreaker{failureThreshold: failureThreshold, cooldownPeriod: cooldownPeriod, now: now}
}

// TODO: implemement mutex for any shared state access (state, fail and success count, etc)
func (cb *circuitBreaker) Execute() error { return nil }

func (cb *circuitBreaker) CurrentState() CircuitBreakerState {
	return cb.state
}

func (cb *circuitBreaker) RecordFailure() {
	cb.failureCount++
}

func (cb *circuitBreaker) RecordSuccess() {
	cb.successCount++

	// if there's 1 request success while HalfOpen, set state to Closed
	if cb.CurrentState() == HalfOpen {
		cb.SetCurrentState(Closed)
		cb.probeInFlight = false
	}
}

func (cb *circuitBreaker) SetCurrentState(state CircuitBreakerState) {
	if cb.state == Open {
		cb.openedAt = cb.now()
	}
	cb.state = state
}

func (cb *circuitBreaker) Allow() bool {
	switch cb.CurrentState() {
	case Closed:
		return true
	case Open:
		if cb.now().Sub(cb.openedAt) < cb.cooldownPeriod {
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
