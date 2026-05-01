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
