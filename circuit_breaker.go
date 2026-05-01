package circuitbreaker

type CircuitBreakerState int

const (
	Closed CircuitBreakerState = iota
	Open
	HalfOpen
)

type CircuiBreaker interface {
	Execute() error
	CurrentState() CircuitBreakerState
	RecordFailure() error
	RecordSuccess() error
}

type circuitBreaker struct {
	state        CircuitBreakerState
	failCount    int
	successCount int
	threshold    int
}

func NewCircuitBreaker(threshold int) CircuiBreaker {
	return &circuitBreaker{threshold: threshold}
}

func (cb *circuitBreaker) Execute() error { return nil }

func (cb *circuitBreaker) CurrentState() CircuitBreakerState { return Closed }

func (cb *circuitBreaker) RecordFailure() error { return nil }

func (cb *circuitBreaker) RecordSuccess() error { return nil }
