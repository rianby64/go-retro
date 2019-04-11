package retro

import (
	"errors"
	"time"
)

// ErrorMaxAttemptsReached whatever
var ErrorMaxAttemptsReached = errors.New("MaxAttempts reached")

// ErrorDelayIsZero whatever
var ErrorDelayIsZero = errors.New("Delay is zero")

// Retryable whatever
type Retryable interface {
	increaseDelay() (time.Duration, error)
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

			duration, err := retryable.increaseDelay()
			if err != nil {
				return errCallback
			}
			if duration > 0 {
				time.Sleep(duration)
			}
			continue
		}
		break
	}
	return nil
}
