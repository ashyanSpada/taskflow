package taskflow

import (
	"context"
	"errors"
	"testing"
)

type testState struct {
	A int
	B int
}

func (t testState) Merge(b testState) testState {
	if t.A == 0 && b.A != 0 {
		t.A = b.A
	}
	if t.B == 0 && b.B != 0 {
		t.B = b.B
	}
	return t
}

func TestSequence(t *testing.T) {
	ctx := context.TODO()
	s := testState{}
	a, err := Sequence(
		func(ctx context.Context, s testState) (testState, error) {
			return s, errors.New("err happen")
		},
		func(ctx context.Context, s testState) (testState, error) {
			return s, nil
		},
	)(ctx, s)
	if err == nil {
		t.Error("should be err")
		return
	}
	a, err = Sequence(
		func(ctx context.Context, s testState) (testState, error) {
			s.A = 123
			s.B = 1234
			return s, nil
		},
		func(ctx context.Context, s testState) (testState, error) {
			s.B = 1234
			return s, nil
		},
	)(ctx, s)
	if err != nil {
		t.Error(err)
		return
	}
	if a.A != 123 || a.B != 1234 {
		t.Error("should equal")
		return
	}
}

func TestOrderedChoice(t *testing.T) {
	ctx := context.TODO()
	s := testState{}
	a, err := OrderedChoice(
		func(ctx context.Context, s testState) (testState, error) {
			return s, errors.New("err happen")
		},
		func(ctx context.Context, s testState) (testState, error) {
			s.B = 12345
			return s, nil
		},
	)(ctx, s)
	if err != nil {
		t.Error("should be err")
		return
	}
	if a.B != 12345 {
		t.Error("should equal")
		return
	}
	a, err = OrderedChoice(
		func(ctx context.Context, s testState) (testState, error) {
			s.A = 123
			return s, nil
		},
		func(ctx context.Context, s testState) (testState, error) {
			s.B = 1234
			return s, nil
		},
	)(ctx, s)
	if err != nil {
		t.Error(err)
		return
	}
	if a.A != 123 && a.B == 1234 {
		t.Error("should equal")
		return
	}
	a, err = OrderedChoice(
		func(ctx context.Context, s testState) (testState, error) {
			return s, errors.New("err happen")
		},
		func(ctx context.Context, s testState) (testState, error) {
			return s, errors.New("err happen")
		},
	)(ctx, s)
	if err == nil {
		t.Error(err)
		return
	}
	a, err = OrderedChoice[testState]()(ctx, s)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestConcurrent(t *testing.T) {
	ctx := context.TODO()
	s := testState{}
	a, err := Concurrent(
		func(ctx context.Context, s testState) (testState, error) {
			return s, errors.New("err happen")
		},
		func(ctx context.Context, s testState) (testState, error) {
			return s, nil
		},
	)(ctx, s)
	if err == nil {
		t.Error("should be err")
		return
	}
	a, err = Concurrent(
		func(ctx context.Context, s testState) (testState, error) {
			s.A = 123
			s.B = 1234
			return s, nil
		},
		func(ctx context.Context, s testState) (testState, error) {
			s.B = 1234
			return s, nil
		},
	)(ctx, s)
	if err != nil {
		t.Error(err)
		return
	}
	if a.A != 123 || a.B != 1234 {
		t.Error("should equal")
		return
	}
}
