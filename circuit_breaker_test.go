package circuitbreaker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNewCircuitBreaker(t *testing.T) {
	failureThreshold := 3
	cooldownPeriod := 30 * time.Second
	baseTime := time.Now

	cb := NewCircuitBreaker(failureThreshold, cooldownPeriod, baseTime)

	assert.Equal(t, Closed, cb.CurrentState())
	assert.Equal(t, 0, cb.GetFailureCount())
	assert.Equal(t, false, cb.IsProbeInFlight())
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
