package infra

import (
	"net/url"

	"github.com/m-mizutani/spout/pkg/model"
)

type LogReader interface {
	Get(ctx *model.Context) (chan *model.Message, error)
}

type Repository interface {
	Put(ctx *model.Context, logs ...*model.Log) error
	Get(ctx *model.Context, input *model.RepositoryGetInput) (*model.RepositoryGetOutput, error)
}

type Browser interface {
	Open(uri *url.URL) error
}
