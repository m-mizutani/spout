package file

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/spout/pkg/model"
	"github.com/m-mizutani/spout/pkg/utils"
)

type Client struct {
	path []string
}

func New(path []string) *Client {
	client := &Client{
		path: path,
	}

	return client
}

func readFile(path string, ch chan *model.Message) {
	fd, err := os.Open(filepath.Clean(path))
	if err != nil {
		ch <- &model.Message{
			Error: goerr.Wrap(err, "open local log file").With("path", path),
		}
		return
	}
	defer func() {
		if err := fd.Close(); err != nil {
			utils.Logger.With("error", err.Error()).Warn("failed to close local log file")
		}
	}()

	scanner := bufio.NewScanner(fd)
	for scanner.Scan() {
		line := scanner.Bytes()
		out := new(any)

		if err := json.Unmarshal(line, &out); err != nil {
			ch <- &model.Message{
				Error: goerr.Wrap(err, "parsing log file").With("path", path).With("original", string(line)),
			}
			return
		}

		ch <- &model.Message{
			Log: &model.Log{
				Data: out,
			},
		}
	}
}

func (x *Client) Get(ctx *model.Context) (chan *model.Message, error) {
	ch := make(chan *model.Message)

	go func() {
		defer close(ch)
		for _, filePath := range x.path {
			readFile(filePath, ch)
		}
	}()

	return ch, nil
}
