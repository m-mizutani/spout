package usecase

import (
	"encoding/json"
	"strconv"

	"github.com/itchyny/gojq"
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/spout/pkg/model"
	"github.com/m-mizutani/spout/pkg/utils"
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
	count := 0
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
		count++
	}

	utils.Logger.With("count", count).Debug("imported logs")

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

	var filter func(log *model.Log) []*model.Log
	if opt.query != "" {
		query, err := gojq.Parse(opt.query)
		if err != nil {
			return nil, goerr.Wrap(err, "failed to parse query")
		}
		code, err := gojq.Compile(query)
		if err != nil {
			return nil, goerr.Wrap(err, "failed to compile query")
		}

		filter = func(log *model.Log) []*model.Log {
			raw, err := json.Marshal(log.Data)
			if err != nil {
				utils.Logger.With("log", log.Data).With("err", err.Error()).Warn("failed to marshal log.Data")
				return nil
			}
			var obj any
			if err := json.Unmarshal(raw, &obj); err != nil {
				utils.Logger.With("log", log.Data).With("err", err.Error()).Warn("failed to unmarshal log.Data")
				return nil
			}

			iter := code.Run(obj)

			var resp []*model.Log
			for {
				v, ok := iter.Next()
				if !ok {
					break
				}
				if err, ok := v.(error); ok {
					utils.Logger.With("error", err.Error()).Debug("jq runtime error")
					continue // ignore
				}

				resp = append(resp, model.NewLog(log.Timestamp, log.Tag, v))
			}

			return resp
		}
	}

	logs, err := x.clients.Repository().Get(ctx, &model.RepositoryGetOption{
		Offset: opt.offset,
		Limit:  opt.limit,
		Filter: filter,
	})
	if err != nil {
		return nil, err
	}

	return logs, nil
}
