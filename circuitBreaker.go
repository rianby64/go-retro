package retro

import (
	"time"
)

// CircuitBreaker strategy
type CircuitBreaker struct {
	Error     error
	lastError error

	currentAttempt int
	BanTimeout     time.Duration
	lastTry        time.Time
	MaxAttempts    int
	Delay          time.Duration

	Execute     Execute
	Recovery    Recovery
	ShouldRetry ShouldRetry
}

func (r *CircuitBreaker) getExecute() (Execute, error) {
	return r.Execute, nil
}

func (r *CircuitBreaker) getRecovery() Recovery {
	return r.Recovery
}

func (r *CircuitBreaker) getShouldRetry() ShouldRetry {
	return r.ShouldRetry
}

func (r *CircuitBreaker) increaseAttempt() error {
	r.currentAttempt++
	if r.currentAttempt >= r.MaxAttempts {
		return ErrorMaxAttemptsReached
	}
	return nil
}

func (r *CircuitBreaker) increaseDelay() (time.Duration, error) {
	return r.Delay, nil
}

func (r *CircuitBreaker) setError(err error) {
	r.Error = err
}

// Run whatever
func (r *CircuitBreaker) Run() (err error) {
	var execute Execute
	if !(r.lastTry.IsZero()) {
		now := time.Now()
		if now.Sub(r.lastTry) > r.BanTimeout {
			r.currentAttempt = 0
			r.setError(nil)
		}
	}
	r.lastTry = time.Now()
	if r.currentAttempt >= r.MaxAttempts {
		err = ErrorBanAttemptsReached
		r.setError(err)
		return r.lastError
	}
	err = r.increaseAttempt()
	if err != nil {
		r.setError(err)
	}
	execute, err = r.getExecute()
	if err != nil {
		r.setError(err)
		return r.lastError
	}
	r.lastError = execute()
	return r.lastError
}
