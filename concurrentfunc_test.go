package concurrentfunc_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/nurjamil/concurrentfunc"
	"github.com/stretchr/testify/assert"
)

type AnimalType int

const (
	Type AnimalType = iota
	DOG
)

func (a AnimalType) String() string {
	switch a {
	case DOG:
		return "DOG"
	}

	return "UNKNOWN"
}

func TestConcurrentFunc(t *testing.T) {
	type People struct {
		Name string
	}

	type Animal struct {
		Type AnimalType
	}

	firstFunc := func(ctx context.Context) (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		firstFuncRes := People{Name: "example"}

		return firstFuncRes, nil
	}
	secondFunc := func(ctx context.Context) (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		secondFuncRes := Animal{Type: DOG}

		return secondFuncRes, nil
	}

	res, errs, err := concurrentfunc.Exec(context.Background(), time.Second, firstFunc, secondFunc)

	assert.Equal(t, nil, err)
	assert.Equal(t, 2, len(res))
	assert.Equal(t, 2, len(errs))
	assert.Equal(t, "example", res[0].(People).Name)
	assert.Equal(t, "DOG", res[1].(Animal).Type.String())
	assert.Equal(t, nil, errs[0])
	assert.Equal(t, nil, errs[1])

}

func TestConcurrent__ContextDeadlineError(t *testing.T) {
	firstFunc := func(ctx context.Context) (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		firstFuncRes := "success"

		return firstFuncRes, nil
	}
	secondFunc := func(ctx context.Context) (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		secondFuncRes := 20

		return secondFuncRes, nil
	}

	res, _, err := concurrentfunc.Exec(context.Background(), time.Nanosecond, firstFunc, secondFunc)

	assert.Equal(t, context.DeadlineExceeded, err)
	assert.Equal(t, nil, res[0])
	assert.Equal(t, nil, res[1])

}

func TestConcurrent__OneFuncReturnError(t *testing.T) {
	mockError := errors.New(`failed exec`)
	firstFunc := func(ctx context.Context) (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, mockError
	}
	secondFunc := func(ctx context.Context) (any, error) {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		return nil, nil
	}

	_, _, err := concurrentfunc.Exec(context.Background(), time.Second, firstFunc, secondFunc)

	assert.Equal(t, mockError, err)

}
