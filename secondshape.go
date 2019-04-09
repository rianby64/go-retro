package retro

import "errors"

// Retry strategy
type Retry struct {
	Error          error
	MaxAttempts    int
	currentAttempt int
}

func (r *Retry) increaseAttempt() error {
	r.currentAttempt++
	if r.currentAttempt >= r.MaxAttempts {
		r.Error = errors.New("MaxAttempts reached")
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
		errCallback := fn()
		if errCallback != nil {
			if err != nil {
				return errCallback
			}
			continue
		}
		break
	}
	return nil
}
