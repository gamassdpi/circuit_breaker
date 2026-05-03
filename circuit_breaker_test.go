package circuitbreaker

import (
	// "fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCircuitBreaker_InitialState(t *testing.T) {
	failureThreshold := 3
	cooldownPeriod := 30 * time.Second
	baseTime := time.Now()

	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod, func() time.Time { return baseTime })

	assert.Equal(t, Closed, cb.CurrentState())
	assert.Equal(t, 0, cb.GetFailureCount())
	assert.Equal(t, false, cb.IsProbeInFlight())
}

func TestAllow_WhenClosed_ReturnsTrue(t *testing.T) {
	failureThreshold := 3
	cooldownPeriod := 30 * time.Second
	baseTime := time.Now()

	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod, func() time.Time { return baseTime })

	assert.Equal(t, true, cb.Allow())
}

func TestRecordFailure_BelowThreshold_StaysClosed(t *testing.T) {
	failureThreshold := 3
	cooldownPeriod := 30 * time.Second
	baseTime := time.Now()

	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod, func() time.Time { return baseTime })

	cb.RecordFailure()
	cb.RecordFailure()

	assert.Equal(t, 2, cb.GetFailureCount())
	assert.Equal(t, Closed, cb.CurrentState())
}

func TestRecordFailure_ReachesTreshold_TransitionsToOpen(t *testing.T) {
	failureThreshold := 3
	cooldownPeriod := 30 * time.Second
	baseTime := time.Now()

	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod, func() time.Time { return baseTime })

	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordFailure()

	assert.Equal(t, Open, cb.CurrentState())
}

func TestAllow_WhenOpen_BeforeCooldown_ReturnsFalse(t *testing.T) {
	failureThreshold := 3
	cooldownPeriod := 30 * time.Second
	baseTime := time.Now()

	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod, func() time.Time { return baseTime })

	// trigger state transitions to Open
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordFailure()

	// add time before cooldown period
	baseTime = baseTime.Add(29 * time.Second)

	assert.Equal(t, false, cb.Allow())
}

func TestAllow_WhenOpen_AfterCooldown_TransitionToHalfOpen_ReturnTrue(t *testing.T) {
	failureThreshold := 3
	cooldownPeriod := 30 * time.Second
	baseTime := time.Now()

	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod, func() time.Time { return baseTime })

	// trigger state transitions to Open
	cb.RecordFailure()
	cb.RecordFailure()
	cb.RecordFailure()

	// add time before cooldown period
	baseTime = baseTime.Add(31 * time.Second)
	// fmt.Println(cb.IsProbeInFlight()) // false
	assert.Equal(t, true, cb.Allow())
	assert.Equal(t, HalfOpen, cb.CurrentState())
}

func TestAllow_ShouldReturnFalse_WhenOpen_AndTimeoutNotExpired(t *testing.T) {
	baseTime := time.Now()
	cb := &circuitBreaker{
		now:            func() time.Time { return baseTime },
		state:          Open,
		openedAt:       baseTime,
		cooldownPeriod: 30 * time.Second,
	}

	baseTime = baseTime.Add(29 * time.Second)
	assert.Equal(t, false, cb.Allow())
}
