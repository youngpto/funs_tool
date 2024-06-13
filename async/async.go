package async

import (
	"context"
	"github.com/youngpto/funs_tool/logger"
	"reflect"
)

/*
	future := async.Async(func() interface{} {
		time.Sleep(time.Second)
		return 1
	})
	cfg, err := future.Await()
*/

type Future[T any] interface {
	Await(ctx context.Context) (T, error)
}

type future[T any] struct {
	await func(ctx context.Context) (T, error)
}

func (f future[T]) Await(ctx context.Context) (T, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	return f.await(ctx)
}

func Exec[T any](f func() (T, error)) Future[T] {
	var result T
	var err error

	c := make(chan struct{})

	go func() {
		defer func() {
			err := recover()
			if err != nil {
				logger.Error("async error %v", err)
			}
			close(c)
		}()
		result, err = f()
	}()

	return future[T]{
		await: func(ctx context.Context) (T, error) {
			select {
			case <-ctx.Done():
				return result, err
			case <-c:
				return result, err
			}
		},
	}
}

func Go(function interface{}, param ...interface{}) {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				logger.Error("safe goroutine error %v", err)
			}
		}()
		fv := reflect.ValueOf(function)
		paramList := make([]reflect.Value, len(param))
		for i := 0; i < len(param); i++ {
			paramList[i] = reflect.ValueOf(param[i])
		}
		_ = fv.Call(paramList)
	}()
}
