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

	// DecideToRetry defines the function that states if execute the next attempt or not
	DecideToRetry DecideToRetry

	// Recovery defines the function that performs an extra code to help the next attempt to execute successful
	Recovery Recovery
}

func (r *Retry) increaseDelay() (time.Duration, error) {
	r.currentDelay = r.Delay
	return r.currentDelay, nil
}

func (r *Retry) increaseAttempt() error {
	r.currentAttempt++
	if r.currentAttempt >= r.MaxAttempts {
		r.Error = ErrorMaxAttemptsReached
		return r.Error
	}
	return nil
}

func (r *Retry) getRecovery() Recovery {
	return r.Recovery
}

func (r *Retry) getDecideToRetry() DecideToRetry {
	return r.DecideToRetry
}

// Run the execute function and behaves according to the strategy
func (r *Retry) Run(execute Execute) error {
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
