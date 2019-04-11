package retro

import (
	"errors"
	"testing"
	"time"
)

func TestTryRetryWaitIncreaseFnStrategyMaxAttemptZeroDelay(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")

	strategy := RetryWaitIncrease{
		MaxAttempts: 5,
	}

	errActual := strategy.Run(func() error {
		currentAttempt++
		return errExpected
	})

	if errActual != errExpected {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error == nil {
		t.Error("strategy must hold an error")
	}
	if !(strategy.Error == ErrorDelayIsZero) {
		t.Error("strategy error should be Delay Is Zero")
	}
}

func TestTryRetryWaitIncreaseFnStrategyMaxAttemptNonZeroDelay(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")

	strategy := RetryWaitIncrease{
		MaxAttempts: 5,
		Delay:       time.Millisecond * 50,
	}
	startTime := time.Now()

	errActual := strategy.Run(func() error {
		currentAttempt++
		endTime := time.Now()
		if endTime.Sub(startTime) < strategy.Delay*time.Duration(currentAttempt-1) {
			t.Error("Incorrect increase delay")
		}
		startTime = time.Now()
		return errExpected
	})

	if errActual != errExpected {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error == nil {
		t.Error("strategy must hold an error")
	}
	if !(strategy.Error == ErrorMaxAttemptsReached) {
		t.Error("strategy error should be Delay Is Zero")
	}
}
