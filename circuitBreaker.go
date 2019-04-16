package retro

import (
	"time"
)

// CircuitBreaker strategy
type CircuitBreaker struct {
	Error     error
	lastError error

	lastTry           time.Time
	currentAttempt    int
	currentBanTimeout time.Duration
	BanTimeout        time.Duration
	MaxAttempts       int

	Execute     Execute
	Recovery    Recovery
	ShouldRetry ShouldRetry
}

func (r *CircuitBreaker) getExecute() (Execute, error) {
	if r.Execute == nil {
		return nil, ErrorExecuteFunctionNil
	}
	return r.Execute, nil
}

func (r *CircuitBreaker) getRecovery() Recovery {
	return r.Recovery
}

func (r *CircuitBreaker) increaseAttempt() error {
	r.currentAttempt++
	if r.currentAttempt >= r.MaxAttempts {
		return ErrorMaxAttemptsReached
	}
	return nil
}

func (r *CircuitBreaker) increaseBanTimeout() (time.Duration, error) {
	r.currentBanTimeout = r.BanTimeout
	if r.BanTimeout == 0 {
		return 0, ErrorBanTimeoutIsZero
	}
	return r.currentBanTimeout, nil
}

func (r *CircuitBreaker) setError(err error) {
	r.Error = err
}

// Run whatever
func (r *CircuitBreaker) Run() (err error) {
	var execute Execute
	var currentBanTimeout time.Duration
	execute, err = r.getExecute()
	if err != nil {
		r.setError(err)
		return r.lastError
	}
	currentBanTimeout, err = r.increaseBanTimeout()
	if err != nil {
		r.setError(err)
		return r.lastError
	}
	if !(r.lastTry.IsZero()) {
		now := time.Now()
		if now.Sub(r.lastTry) > currentBanTimeout {
			r.currentAttempt = 0
			r.setError(nil)
		}
	}
	r.lastTry = time.Now()
	if r.MaxAttempts == 0 {
		r.setError(ErrorMaxAttemptsIsZero)
		return r.lastError
	}
	if r.currentAttempt >= r.MaxAttempts {
		err = ErrorBanAttemptsReached
		r.setError(err)
		return r.lastError
	}
	err = r.increaseAttempt()
	if err != nil {
		r.setError(err)
	}
	r.lastError = execute()
	return r.lastError
}
