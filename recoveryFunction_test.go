package retro

import (
	"errors"
	"testing"
	"time"
)

func TestTryFnStrategyMaxAttemptWithRecoveryFunction(t *testing.T) {
	currentAttempt := 0
	currentRecovery := 0
	errExpected := errors.New("expected error")

	startTime := time.Now()

	strategy := Retry{
		MaxAttempts: 5,
		Delay:       time.Millisecond * 100,
		Execute: func() error {
			currentAttempt++
			return errExpected
		},
		Recovery: func(err error) error {
			currentRecovery++
			return nil
		},
	}

	errActual := strategy.Run()

	endTime := time.Now()

	if endTime.Sub(startTime) < time.Millisecond*4*100 {
		t.Error("Incorrect delay")
	}

	if errActual != errExpected {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error == nil {
		t.Error("strategy must hold an error")
	}
	if strategy.MaxAttempts != currentAttempt {
		t.Error("currentAttempt differs from MaxAttempts")
	}
	if strategy.MaxAttempts-1 != currentRecovery {
		t.Error("currentRecovery differs from MaxAttempts", currentRecovery)
	}
}

func TestTryFnStrategyRecoveryFunctionRecoversExecution(t *testing.T) {
	currentAttempt := 0
	currentRecovery := 0
	errExpected := errors.New("expected error")

	strategy := Retry{
		MaxAttempts: 5,
		Delay:       time.Millisecond * 100,
		Execute: func() error {
			currentAttempt++
			if currentRecovery == 2 {
				return nil
			}
			return errExpected
		},
		Recovery: func(err error) error {
			currentRecovery++
			return nil
		},
	}

	errActual := strategy.Run()

	if errActual != nil {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error != nil {
		t.Error("strategy must hold an error")
	}
	if currentAttempt >= strategy.MaxAttempts {
		t.Error("currentAttempt differs from MaxAttempts")
	}
	if currentRecovery != 2 {
		t.Error("currentRecovery differs from expected", currentRecovery)
	}
	if currentAttempt != 3 {
		t.Error("currentAttempt differs from expected", currentAttempt)
	}
}
