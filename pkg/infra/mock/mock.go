package mock

import (
	"time"

	"github.com/m-mizutani/gots/ptr"
	"github.com/m-mizutani/spout/pkg/model"
)

type LogReader struct {
	logs []any
}

func NewLogReader(logs []any) *LogReader {
	return &LogReader{
		logs: logs,
	}
}

func (x *LogReader) Get(ctx *model.Context) (chan *model.Message, error) {
	ch := make(chan *model.Message)

	go func() {
		defer close(ch)
		for i := range x.logs {
			ch <- &model.Message{
				Log: &model.Log{
					Timestamp: ptr.To(time.Now()),
					Data:      x.logs[i],
				},
			}
		}
	}()

	return ch, nil
}
