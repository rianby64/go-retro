package retro

import (
	"errors"
	"time"
)

// ErrorRecoveryFunctionError whatever
var ErrorRecoveryFunctionError = errors.New("Recovery function finished with error")

// ErrorExecuteFunctionNil whatever
var ErrorExecuteFunctionNil = errors.New("Execute function has not beed defined")

// ErrorMaxAttemptsReached whatever
var ErrorMaxAttemptsReached = errors.New("MaxAttempts reached")

// ErrorDelayIsZero whatever
var ErrorDelayIsZero = errors.New("Delay is zero")

// Execute function
type Execute func() error

// Recovery function
type Recovery func(error) error

// DecideToRetry function
type DecideToRetry func(error) bool

type strategy interface {
	getExecute() (Execute, error)
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
