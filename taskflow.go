package taskflow

import (
	"context"
	"fmt"

	"golang.org/x/sync/errgroup"
)

type State[T any] interface {
	Merge(T) T
}

type Handler[T State[T]] func(context.Context, T) (T, error)

func Sequence[T State[T]](handlers ...Handler[T]) Handler[T] {
	return func(ctx context.Context, s T) (T, error) {
		var err error
		for _, hanlder := range handlers {
			s, err = hanlder(ctx, s)
			if err != nil {
				return s, fmt.Errorf("sequence fail: %w", err)
			}
		}
		return s, nil
	}
}

func OrderedChoice[T State[T]](handlers ...Handler[T]) Handler[T] {
	return func(ctx context.Context, s T) (T, error) {
		var errs []error
		for _, hanlder := range handlers {
			newState, err := hanlder(ctx, s)
			if err == nil {
				return newState, nil
			} else {
				errs = append(errs, err)
			}
		}
		if len(errs) > 0 {
			return s, fmt.Errorf("orderedChoice fail: %v", errs)
		}
		return s, nil
	}
}

func Concurrent[T State[T]](handlers ...Handler[T]) Handler[T] {
	return func(ctx context.Context, s T) (T, error) {
		var wg errgroup.Group
		for _, handler := range handlers {
			handler := handler
			wg.Go(func() error {
				newState, err := handler(ctx, s)
				if err != nil {
					return err
				}
				s = s.Merge(newState)
				return nil
			})
		}
		err := wg.Wait()
		if err != nil {
			return s, fmt.Errorf("concurrent fail: %w", err)
		}
		return s, nil
	}
}
