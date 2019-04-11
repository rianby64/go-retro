package retro

import (
	"errors"
	"testing"
	"time"
)

func TestTryFnStrategyMaxAttempt(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")

	strategy := Retry{
		MaxAttempts: 5,
	}

	errActual := Try(func() error {
		currentAttempt++
		return errExpected
	}, &strategy)

	if errActual != errExpected {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error == nil {
		t.Error("strategy must hold an error")
	}
	if strategy.MaxAttempts != currentAttempt {
		t.Error("currentAttempt differs from MaxAttempts")
	}
}

func TestTryFnStrategyMaxAttemptWithDelay(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")

	startTime := time.Now()

	strategy := Retry{
		MaxAttempts: 5,
		Delay:       time.Millisecond * 100,
	}

	errActual := Try(func() error {
		currentAttempt++
		return errExpected
	}, &strategy)

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
}

func TestTryFnStrategyMaxAttemptWithDelayBelowMaxAttempts(t *testing.T) {
	currentAttempt := 0
	StopAttempt := 3
	errExpected := errors.New("expected error")

	startTime := time.Now()

	strategy := Retry{
		MaxAttempts: 5,
		Delay:       time.Millisecond * 100,
	}

	errActual := Try(func() error {
		currentAttempt++
		if currentAttempt == StopAttempt {
			return nil
		}
		return errExpected
	}, &strategy)

	endTime := time.Now()

	if endTime.Sub(startTime) > time.Millisecond*4*100 {
		t.Error("Incorrect delay")
	}

	if errActual != nil {
		t.Error("fn returns an unexpected error")
	}
	if strategy.Error != nil {
		t.Error("strategy must hold an error")
	}
	if StopAttempt != currentAttempt {
		t.Error("currentAttempt differs from MaxAttempts")
	}
}
