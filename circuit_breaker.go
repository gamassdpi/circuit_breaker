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

type CircuitBreaker interface {
	Execute() error
	CurrentState() CircuitBreakerState
	RecordFailure()
	Allow() bool

	GetFailureCount() int
	IsProbeInFlight() bool
}

// TODO: implement mutex
type circuitBreaker struct {
	state            CircuitBreakerState
	failureCount     int
	failureThreshold int
	probeInFlight    bool
	cooldownPeriod   time.Duration
	now              func() time.Time
	openedAt         time.Time
}

func NewCircuitBreaker(failureThreshold int, cooldownPeriod time.Duration, now func() time.Time) CircuitBreaker {
	return &circuitBreaker{failureThreshold: failureThreshold, cooldownPeriod: cooldownPeriod, now: now}
}

// TODO: implemement mutex for any shared state access (state, fail and success count, etc)
func (cb *circuitBreaker) Execute() error { return nil }

func (cb *circuitBreaker) CurrentState() CircuitBreakerState {
	return cb.state
}

func (cb *circuitBreaker) RecordFailure() {
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

func (cb *circuitBreaker) Allow() bool {
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

func (cb *circuitBreaker) GetFailureCount() int {
	return cb.failureCount
}

func (cb *circuitBreaker) IsProbeInFlight() bool {
	return cb.probeInFlight
}
