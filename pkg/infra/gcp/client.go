package gcp

import (
	"fmt"
	"time"

	"cloud.google.com/go/logging/logadmin"
	"google.golang.org/api/iterator"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/spout/pkg/model"
	"github.com/m-mizutani/spout/pkg/utils"
)

type Client struct {
	projectID model.GoogleProjectID
	limit     int
	filter    string
	begin     time.Time
	end       time.Time
}

func New(projectID model.GoogleProjectID, limit int, period model.Period, options ...Option) *Client {
	client := &Client{
		projectID: projectID,
		limit:     limit,
		begin:     period.Begin(),
		end:       period.End(),
	}

	for _, opt := range options {
		opt(client)
	}

	return client
}

type Option func(x *Client)

func WithFilter(filter string) Option {
	return func(x *Client) {
		x.filter = filter
	}
}

type mapper interface {
	AsMap() map[string]interface{}
}

func (x *Client) Get(ctx *model.Context) (chan *model.Message, error) {
	adminClient, err := logadmin.NewClient(ctx, string(x.projectID))
	if err != nil {
		return nil, goerr.Wrap(err, "creating logadmin client")
	}

	filter := fmt.Sprintf(`timestamp >= "%s" timestamp <= "%s"`,
		x.begin.Format("2006-01-02T15:04:05Z"),
		x.end.Format("2006-01-02T15:04:05Z"),
	)
	if x.filter != "" {
		filter += " " + x.filter
	}
	utils.Logger.With("filter", filter).Debug("starting log download")

	iter := adminClient.Entries(ctx,
		logadmin.Filter(filter),
		logadmin.NewestFirst(),
	)

	ch := make(chan *model.Message)
	go func() {
		defer close(ch)
		defer adminClient.Close()

		for i := 0; i < x.limit; i++ {
			entry, err := iter.Next()
			if err == iterator.Done {
				return
			}
			if err != nil {
				ch <- &model.Message{Error: goerr.Wrap(err)}
				return
			}
			if entry == nil {
				utils.Logger.Warn("entry is null")
				continue
			}

			var payload any
			if m, ok := entry.Payload.(mapper); ok {
				payload = m.AsMap()
			} else {
				payload = fmt.Sprintf("%T", entry.Payload)
			}

			ch <- &model.Message{
				Log: &model.Log{
					Timestamp: entry.Timestamp,
					Data: &model.CloudLoggingLog{
						Severity: entry.Severity.String(),
						InsertID: entry.InsertID,
						Payload:  payload,
						LogName:  entry.LogName,
						Resource: entry.Resource,
						Trace:    entry.Trace,
						SpanID:   entry.SpanID,
					},
				},
			}
		}
	}()

	return ch, nil
}
