package circuitbreaker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type testClock struct {
	now time.Time
}

func (c *testClock) Now() time.Time {
	return c.now
}

const (
	failureThreshold int           = 3
	cooldownPeriod   time.Duration = 30 * time.Second
)

func newTestCB(t *testing.T) (*CircuitBreaker, *testClock) {
	t.Helper()
	clock := &testClock{now: time.Now()} // freeze time

	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod, clock.Now)
	return cb, clock
}

func triggerOpen(cb *CircuitBreaker) {
	for range failureThreshold {
		cb.RecordFailure()
	}
}

func TestNewCircuitBreaker(t *testing.T) {
	t.Run("initial state - is closed", func(t *testing.T) {
		cb, _ := newTestCB(t)
		assert.Equal(t, Closed, cb.CurrentState())
	})

	t.Run("initial failureCount - is 0", func(t *testing.T) {
		cb, _ := newTestCB(t)
		assert.Equal(t, 0, cb.GetFailureCount())
	})

	t.Run("initial probeInFligth - is false", func(t *testing.T) {
		cb, _ := newTestCB(t)
		assert.Equal(t, false, cb.IsProbeInFlight())
	})
}

func TestAllow(t *testing.T) {
	t.Run("when state is closed - return true", func(t *testing.T) {
		cb, _ := newTestCB(t)

		assert.Equal(t, true, cb.Allow())
	})

	t.Run("when open before cooldown - return false", func(t *testing.T) {
		cb, clock := newTestCB(t)
		triggerOpen(cb)

		clock.now = clock.Now().Add(cooldownPeriod).Add(-1 * time.Second) // ~1s before cooldown
		assert.Equal(t, false, cb.Allow())
	})

	t.Run("when open after cooldown - return true and state transition to HalfOpen", func(t *testing.T) {
		cb, clock := newTestCB(t)
		triggerOpen(cb)

		clock.now = clock.Now().Add(cooldownPeriod).Add(1 * time.Second)
		assert.Equal(t, true, cb.Allow())
		assert.Equal(t, HalfOpen, cb.CurrentState())
	})

	t.Run("when half open and probeInFlight false - return true", func(t *testing.T) {
		cb, clock := newTestCB(t)
		triggerOpen(cb)
		clock.now = clock.now.Add(cooldownPeriod).Add(1 * time.Second)
		// trigger state transition from Open to HalfOpen with probeInFlight is false
		cb.Allow()
		
		assert.Equal(t, true, cb.Allow())
	})

	t.Run("when half open and probeInFlight true - return false", func(t *testing.T) {
		cb, clock := newTestCB(t)
		triggerOpen(cb)
		clock.now = clock.now.Add(cooldownPeriod).Add(1 * time.Second)
		// trigger state transition from Open to HalfOpen with probeInFlight is false
		cb.Allow()
		// trigger probeInFlight to true
		cb.Allow()

		assert.Equal(t, false, cb.Allow())
	})
}

func TestRecordFailure(t *testing.T) {
	t.Run("when closed and failureCount below threshold - state stays closed", func(t *testing.T){
		cb, _ := newTestCB(t)
		
		for range failureThreshold - 1 {
			cb.RecordFailure()
		}

		assert.Equal(t, Closed, cb.CurrentState())
	})

	t.Run("when closed and reaches failureThreshold - state transitions to Open", func(t *testing.T) {
		cb, _ := newTestCB(t)
		triggerOpen(cb)

		assert.Equal(t, Open, cb.CurrentState())
	})
}
