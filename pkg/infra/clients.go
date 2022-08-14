package infra

import (
	"io"
	"os"

	"github.com/m-mizutani/spout/pkg/infra/memdb"
)

type Clients struct {
	logReader LogReader
	repo      Repository
	writer    io.Writer
	browser   Browser
}

func (x *Clients) LogReader() LogReader   { return x.logReader }
func (x *Clients) Repository() Repository { return x.repo }
func (x *Clients) Writer() io.Writer      { return x.writer }
func (x *Clients) Browser() Browser       { return x.browser }

func New(options ...Option) *Clients {
	clients := &Clients{
		repo:    memdb.New(),
		writer:  os.Stdout,
		browser: &browser{},
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

func WithRepository(repo Repository) Option {
	return func(c *Clients) {
		c.repo = repo
	}
}
