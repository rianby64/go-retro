package retro

import (
	"errors"
	"time"
)

// ErrorMaxAttemptsReached whatever
var ErrorMaxAttemptsReached = errors.New("MaxAttempts reached")

// Retryable whatever
type Retryable interface {
	getDelay() time.Duration
	increaseDelay() time.Duration
	increaseAttempt() error
}

// Callback function
type Callback func() error

// Try a Callback
func Try(fn Callback, retryable Retryable) error {
	for {
		err := retryable.increaseAttempt()
		errCallback := fn()
		if errCallback != nil {
			if err != nil {
				return errCallback
			}

			if retryable.getDelay() > 0 {
				time.Sleep(retryable.increaseDelay())
			}
			continue
		}
		break
	}
	return nil
}
