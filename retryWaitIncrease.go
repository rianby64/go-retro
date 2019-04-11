package retro

import (
	"time"
)

// RetryWaitIncrease strategy
type RetryWaitIncrease struct {
	currentAttempt int
	currentDelay   time.Duration

	// Error holds the error ocurred in the RetryWaitIncrease' scope
	Error error

	// MaxAttempts defines the maximum number of attempts
	MaxAttempts int

	// Delay defines the amount of time between attempts
	Delay time.Duration
}

func (r *RetryWaitIncrease) increaseDelay() (time.Duration, error) {
	r.currentDelay = r.Delay * time.Duration(r.currentAttempt)
	if r.Delay == 0 {
		r.Error = ErrorDelayIsZero
		return 0, r.Error
	}
	return r.currentDelay, nil
}

func (r *RetryWaitIncrease) increaseAttempt() error {
	r.currentAttempt++
	if r.currentAttempt >= r.MaxAttempts {
		r.Error = ErrorMaxAttemptsReached
		return r.Error
	}
	return nil
}
