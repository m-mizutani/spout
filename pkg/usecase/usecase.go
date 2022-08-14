package usecase

import (
	"net/url"

	"github.com/m-mizutani/spout/pkg/infra"
)

type Usecase struct {
	clients *infra.Clients
}

func New(clients *infra.Clients) *Usecase {
	return &Usecase{
		clients: clients,
	}
}

func (x *Usecase) OpenBrowser(uri *url.URL) error {
	return x.clients.Browser().Open(uri)
}
