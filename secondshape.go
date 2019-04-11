package retro

import (
	"errors"
)

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
