package usecase

import (
	"encoding/json"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/spout/pkg/infra"
	"github.com/m-mizutani/spout/pkg/model"
)

type Usecase struct {
	clients *infra.Clients
}

func New(clients *infra.Clients) *Usecase {
	return &Usecase{
		clients: clients,
	}
}

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
