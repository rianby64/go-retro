package retro

import "time"

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

	// DecideToRetry defines the function that states if execute the next attempt or not
	DecideToRetry DecideToRetry

	// Recovery defines the function that performs an extra code to help the next attempt to execute successful
	Recovery Recovery
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

func (r *RetryWaitIncrease) getRecovery() Recovery {
	return r.Recovery
}

func (r *RetryWaitIncrease) getDecideToRetry() DecideToRetry {
	return r.DecideToRetry
}

// Run the execute function and behaves according to the strategy
func (r *RetryWaitIncrease) Run(execute Execute) error {
	for {
		err := r.increaseAttempt()
		errExecute := execute()
		if errExecute != nil {
			if err != nil {
				return errExecute
			}

			duration, err := r.increaseDelay()
			if err != nil {
				return errExecute
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
