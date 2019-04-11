package retro

import "time"

// Retry strategy
type Retry struct {
	currentAttempt int
	currentDelay   time.Duration

	// Error holds the error ocurred in the Retry' scope
	Error error

	// MaxAttempts defines the maximum number of attempts
	MaxAttempts int

	// Delay defines the amount of time between attempts
	Delay time.Duration
}

func (r *Retry) getDelay() time.Duration {
	return r.Delay
}

func (r *Retry) increaseDelay() time.Duration {
	r.currentDelay = r.Delay
	return r.currentDelay
}

func (r *Retry) increaseAttempt() error {
	r.currentAttempt++
	if r.currentAttempt >= r.MaxAttempts {
		r.Error = ErrorMaxAttemptsReached
		return r.Error
	}
	return nil
}
