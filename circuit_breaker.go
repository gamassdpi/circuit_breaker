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
	RecordFailure()
	RecordSuccess()
}

// TODO: implement mutex
type circuitBreaker struct {
	state        CircuitBreakerState
	failCount    int
	successCount int
	threshold    int
}

func NewCircuitBreaker(threshold int) CircuiBreaker {
	return &circuitBreaker{threshold: threshold, state: Closed}
}

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
