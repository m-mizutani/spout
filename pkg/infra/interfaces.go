package infra

import (
	"github.com/m-mizutani/spout/pkg/model"
)

type LogReader interface {
	Get(ctx *model.Context) (chan *model.Message, error)
}
