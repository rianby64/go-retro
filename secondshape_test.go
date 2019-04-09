package retro

import (
	"errors"
	"testing"
)

func TestSimplestTryFn(t *testing.T) {
	executed := false
	fn := func() error {
		executed = true
		return nil
	}

	Try(fn)
	if executed == false {
		t.Error("fn not called")
	}
}

func TestTryFnStrategyMaxAttempt(t *testing.T) {
	currentAttempt := 0
	errExpected := errors.New("expected error")
	fn := func() error {
		currentAttempt++
		return errExpected
	}

	strategy := Retry{
		MaxAttempts: 5,
	}

	errActual := TryWithStrategy(fn, &strategy)
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
