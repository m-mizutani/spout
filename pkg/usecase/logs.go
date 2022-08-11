package usecase

import (
	"encoding/json"
	"strconv"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/spout/pkg/model"
)

func (x *Usecase) DumpLogs(ctx *model.Context) error {
	msgCh, err := x.clients.LogReader().Get(ctx)
	if err != nil {
		return err
	}

	encoder := json.NewEncoder(x.clients.Writer())
	encoder.SetIndent("", "  ")
	for msg := range msgCh {
		if msg.Error != nil {
			return msg.Error
		}

		if err := encoder.Encode(&msg.Log.Data); err != nil {
			return goerr.Wrap(err, "json.Encode log.Data")
		}
	}

	return nil
}

func (x *Usecase) ImportLogs(ctx *model.Context) error {
	msgCh, err := x.clients.LogReader().Get(ctx)
	if err != nil {
		return err
	}

	for msg := range msgCh {
		if msg.Error != nil {
			return msg.Error
		}

		if err := x.clients.Repository().Put(ctx, msg.Log); err != nil {
			return err
		}
	}

	return nil
}

type exportLogsOptions struct {
	offset uint64
	limit  uint64
	query  string
}

type ExportLogsOption func(opt *exportLogsOptions) error

func WithOffset(offset string) ExportLogsOption {
	return func(opt *exportLogsOptions) error {
		v, err := strconv.ParseUint(offset, 10, 64)
		if err != nil {
			return goerr.Wrap(err, "invalid offset")
		}
		opt.offset = v
		return nil
	}
}

func WithLimit(limit string) ExportLogsOption {
	return func(opt *exportLogsOptions) error {
		v, err := strconv.ParseUint(limit, 10, 64)
		if err != nil {
			return goerr.Wrap(err, "invalid limit")
		}
		opt.limit = v
		return nil
	}
}

func WithQuery(q string) ExportLogsOption {
	return func(opt *exportLogsOptions) error {
		opt.query = q
		return nil
	}
}

func (x *Usecase) ExportLogs(ctx *model.Context, options ...ExportLogsOption) ([]*model.Log, error) {
	opt := &exportLogsOptions{
		limit:  100,
		offset: 0,
	}

	for _, f := range options {
		if err := f(opt); err != nil {
			return nil, err
		}
	}

	logs, err := x.clients.Repository().Get(ctx, &model.RepositoryGetOption{
		Offset: opt.offset,
		Limit:  opt.limit,
	})
	if err != nil {
		return nil, err
	}

	return logs, nil
}
