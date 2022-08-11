package infra

import (
	"io"
	"os"
)

type Clients struct {
	logReader LogReader
	writer    io.Writer
}

func (x *Clients) LogReader() LogReader { return x.logReader }
func (x *Clients) Writer() io.Writer    { return x.writer }

func New(options ...Option) *Clients {
	clients := &Clients{
		writer: os.Stdout,
	}

	for _, opt := range options {
		opt(clients)
	}

	return clients
}

type Option func(c *Clients)

func WithLogReader(LogReader LogReader) Option {
	return func(c *Clients) {
		c.logReader = LogReader
	}
}

func WithWriter(w io.Writer) Option {
	return func(c *Clients) {
		c.writer = w
	}
}
