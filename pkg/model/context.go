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
