package circuitbreaker

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

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
