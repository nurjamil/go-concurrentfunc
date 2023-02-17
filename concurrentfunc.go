package concurrentfunc

import (
	"context"
	"sync"
	"time"
)

// this library can be use if we have move than 1 function to run at the same time
// it's has timeout
// see the md file to see full example how to properly used it

type ConcurrentFunc func(context.Context) (any, error)

type Entity struct {
	Entity any
	Err    error
}

func Exec(ctx context.Context, timeout time.Duration, funcs ...ConcurrentFunc) ([]any, []error, error) {
	newCtx, cancelNewCtx := context.WithTimeout(ctx, timeout)
	defer cancelNewCtx()
	res := make([]any, len(funcs))
	err := make([]error, len(funcs))
	var rootCause error
	doneChannel := make(chan struct{})

	wg := new(sync.WaitGroup)
	wg.Add(len(funcs))
	mutex := new(sync.Mutex)

	for idx, func_ := range funcs {
		func_ := func_
		idx := idx
		go func() {
			defer wg.Done()
			var tempEntity Entity
			tempEntity.Entity, tempEntity.Err = func_(newCtx)

			// cancel the context to cancel the other function
			// and put
			if tempEntity.Err != nil {
				mutex.Lock()
				if rootCause == nil {
					rootCause = tempEntity.Err
				}
				mutex.Unlock()
				cancelNewCtx()
			}
			res[idx] = tempEntity.Entity
			err[idx] = tempEntity.Err
		}()
	}

	go func() {
		defer close(doneChannel)
		wg.Wait()
	}()

	select {
	case <-time.After(timeout):
		cancelNewCtx()
		return nil, nil, newCtx.Err()
	case <-doneChannel:
	}

	return res, err, rootCause
}
