package retro

import (
	"errors"
	"time"
)

// ErrorBanAttemptsReached occurs when the execute function was executed BanAttempts times
var ErrorBanAttemptsReached = errors.New("BanAttempts reached")

// ErrorShouldRetryFunctionError occurs when ShouldRetry function returns false
var ErrorShouldRetryFunctionError = errors.New("ShouldRetry returned false")

// ErrorRecoveryFunctionError occurs when Recovery function returns an error
var ErrorRecoveryFunctionError = errors.New("Recovery function finished with error")

// ErrorExecuteFunctionNil occurs when Run() is executed without Execute function defined
var ErrorExecuteFunctionNil = errors.New("Execute function has not beed defined")

// ErrorMaxAttemptsReached occurs when the execute function was executed MaxAttempts times
var ErrorMaxAttemptsReached = errors.New("MaxAttempts reached")

// ErrorMaxAttemptsIsZero occurs as Max Attempts is supposed to be > 0
var ErrorMaxAttemptsIsZero = errors.New("Max Attempts is zero")

// ErrorDelayIsZero occurs as Delay is supposed to be > 0
var ErrorDelayIsZero = errors.New("Delay is zero")

// ErrorBanTimeoutIsZero occurs as Ban Timeout is supposed to be > 0
var ErrorBanTimeoutIsZero = errors.New("Ban Timeout is zero")

// Execute function
type Execute func() error

// Recovery function
type Recovery func(error) error

// ShouldRetry function
type ShouldRetry func(error) bool

type strategy interface {
	getExecute() (Execute, error)
	getShouldRetry() ShouldRetry
	getRecovery() Recovery
	increaseAttempt() error
	increaseDelay() (time.Duration, error)
	setError(error)
}

func launcStrategy(s strategy) error {
	execute, err := s.getExecute()
	if err != nil {
		s.setError(err)
		return err
	}
	for {
		err := s.increaseAttempt()
		errExecute := execute()
		if errExecute != nil {
			if err != nil {
				s.setError(err)
				return errExecute
			}

			duration, err := s.increaseDelay()
			if err != nil {
				s.setError(err)
				return errExecute
			}
			shouldRetry := s.getShouldRetry()
			if shouldRetry != nil {
				if shouldRetry(errExecute) == false {
					s.setError(ErrorShouldRetryFunctionError)
					return errExecute
				}
			}
			if duration > 0 {
				time.Sleep(duration)
			}
			recover := s.getRecovery()
			if recover != nil {
				errRecovery := recover(errExecute)
				if errRecovery != nil {
					s.setError(ErrorRecoveryFunctionError)
					return errRecovery
				}
			}
			continue
		}
		break
	}
	return nil
}
