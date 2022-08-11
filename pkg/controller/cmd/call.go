package cmd

import (
	"github.com/m-mizutani/goerr"
	"github.com/m-mizutani/gots/slice"
	"github.com/m-mizutani/spout/pkg/infra"
	"github.com/m-mizutani/spout/pkg/model"
	"github.com/m-mizutani/spout/pkg/usecase"
	"github.com/urfave/cli/v2"
)

func cmdCall(cfg *config) *cli.Command {
	app := &cli.App{
		Name:     "spout",
		Commands: readCommands(cfg),
	}

	var (
		cmdFilePath string
	)

	return &cli.Command{
		Name:        "call",
		Aliases:     []string{"c"},
		Description: "Call predefined commands",
		Usage:       "[cmd name] [additional options...]",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "spout-file",
				Aliases:     []string{"f"},
				Value:       ".spout.toml",
				Destination: &cmdFilePath,
			},
		},
		Action: func(c *cli.Context) error {
			args := c.Args().Slice()
			if len(args) <= 0 {
				return goerr.New("command name required")
			}

			ctx := model.NewContext(model.WithCtx(c.Context))
			uc := usecase.New(infra.New())

			commands, err := uc.LoadCommands(ctx, cmdFilePath)
			if err != nil {
				return err
			}

			cmd, ok := commands[args[0]]
			if !ok {
				return goerr.New("command not found").With("name", args[0])
			}

			newArgs := slice.Flatten(
				[]string{"spout", cmd.Command},
				cmd.Options,
				args[1:],
			)

			if err := app.Run(newArgs); err != nil {
				return goerr.Wrap(err)
			}

			return nil
		},
	}
}
