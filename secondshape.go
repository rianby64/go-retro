package retro

import (
	"errors"
	"time"
)

// ErrorMaxAttemptsReached whatever
var ErrorMaxAttemptsReached = errors.New("MaxAttempts reached")

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

// Callback function
type Callback func() error

// Try a Callback
func Try(fn Callback) error {
	return fn()
}

// TryWithStrategy a Callback
func TryWithStrategy(fn Callback, strategy *Retry) error {
	for {
		err := strategy.increaseAttempt()
		errCallback := Try(fn)
		if errCallback != nil {
			if err != nil {
				return errCallback
			}

			if strategy.Delay > 0 {
				time.Sleep(strategy.increaseDelay())
			}
			continue
		}
		break
	}
	return nil
}
