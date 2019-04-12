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

	// Execute defines the code to be wrapped under this strategy
	Execute Execute

	// DecideToRetry defines the function that states if execute the next attempt or not
	DecideToRetry DecideToRetry

	// Recovery defines the function that performs an extra code to help the next attempt to execute successful
	Recovery Recovery
}

func (r *RetryWaitIncrease) setError(err error) {
	r.Error = err
}

func (r *RetryWaitIncrease) increaseDelay() (time.Duration, error) {
	r.currentDelay = r.Delay * time.Duration(r.currentAttempt)
	if r.Delay == 0 {
		return 0, ErrorDelayIsZero
	}
	return r.currentDelay, nil
}

func (r *RetryWaitIncrease) increaseAttempt() error {
	r.currentAttempt++
	if r.currentAttempt >= r.MaxAttempts {
		return ErrorMaxAttemptsReached
	}
	return nil
}

func (r *RetryWaitIncrease) getExecute() (Execute, error) {
	if r.Execute == nil {
		return nil, ErrorExecuteFunctionNil
	}
	return r.Execute, nil
}

func (r *RetryWaitIncrease) getRecovery() Recovery {
	return r.Recovery
}

func (r *RetryWaitIncrease) getDecideToRetry() DecideToRetry {
	return r.DecideToRetry
}

// Run the execute function and behaves according to the strategy
func (r *RetryWaitIncrease) Run() error {
	return launcStrategy(r)
}