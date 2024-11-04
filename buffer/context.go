package buffer

import (
	"context"
	"errors"

	internal "github.com/PlayerR9/go-safe/buffer/internal"
	"github.com/PlayerR9/go-safe/common"
)

type contextKey struct{}

func fromContext[T any](ctx context.Context) (*Context[T], error) {
	if ctx == nil {
		return nil, common.NewErrNilParam("ctx")
	}

	v, ok := ctx.Value(contextKey{}).(*Context[T])
	if !ok || v == nil {
		return nil, errors.New("expected non-nil *Context[T] in context")
	}

	return v, nil
}

type Context[T any] struct {
	buffer *internal.Buffer[T]
}

func NewContext[T any](parent context.Context) (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(parent)

	c := &Context[T]{}

	pc, err := fromContext[T](parent)
	if err == nil {
		c.buffer = pc.buffer
	} else {
		c.buffer = new(internal.Buffer[T])

		err := c.buffer.Start()
		if err != nil {
			panic(err)
		}
	}

	ctx = context.WithValue(parent, contextKey{}, c)

	cancelFn := func() {
		cancel()

		c.buffer.Close()
	}

	return ctx, cancelFn
}
