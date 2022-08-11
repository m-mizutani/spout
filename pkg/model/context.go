package model

import "context"

type Context struct {
	context.Context
}

func NewContext(options ...CtxOption) *Context {
	ctx := &Context{
		Context: context.Background(),
	}

	return ctx
}

type CtxOption func(ctx *Context)

func WithCtx(base context.Context) CtxOption {
	return func(ctx *Context) {
		ctx.Context = base
	}
}

func (x *Context) New(options ...CtxOption) *Context {
	newCtx := *x

	for _, opt := range options {
		opt(&newCtx)
	}

	return &newCtx
}
