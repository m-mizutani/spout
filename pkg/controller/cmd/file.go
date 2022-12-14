package cmd

import (
	"github.com/m-mizutani/spout/pkg/infra"
	"github.com/m-mizutani/spout/pkg/infra/file"
	"github.com/m-mizutani/spout/pkg/model"
	"github.com/m-mizutani/spout/pkg/usecase"
	"github.com/urfave/cli/v2"
)

func cmdFile(globalCfg *config) *cli.Command {
	var runOpt runOptions
	return &cli.Command{
		Name:        "file",
		Aliases:     []string{"f"},
		Usage:       "[file1, [file2, ...]]",
		Description: "Read local file logs",
		Flags:       runOpt.Flags(),
		Action: func(c *cli.Context) error {
			reader := file.New(c.Args().Slice())
			clients := infra.New(
				infra.WithLogReader(reader),
			)
			uc := usecase.New(clients)
			ctx := model.NewContext(model.WithCtx(c.Context))

			if err := run(ctx, uc, &runOpt); err != nil {
				return err
			}

			return nil
		},
	}
}
