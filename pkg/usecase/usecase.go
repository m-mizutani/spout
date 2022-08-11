package usecase

import (
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
