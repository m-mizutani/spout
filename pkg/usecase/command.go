package usecase

import (
	"os"
	"path/filepath"

	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/spout/pkg/model"
	"github.com/m-mizutani/spout/pkg/utils"

	"github.com/BurntSushi/toml"
)

func (x *Usecase) LoadCommands(ctx *model.Context, filePath string) (map[string]*model.Command, error) {
	path := filepath.Clean(filePath)

	fd, err := os.Open(path)
	if err != nil {
		return nil, goerr.Wrap(err, "can not open command file")
	}
	defer func() {
		if err := fd.Close(); err != nil {
			utils.Logger.With("error", err.Error()).Warn("failed to close command file")
		}
	}()

	var commands map[string]*model.Command
	if _, err := toml.NewDecoder(fd).Decode(&commands); err != nil {
		return nil, goerr.Wrap(err)
	}

	return commands, nil
}
