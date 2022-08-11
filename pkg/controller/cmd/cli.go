package cmd

import (
	"github.com/m-mizutani/spout/pkg/utils"
	"github.com/m-mizutani/zlog"
	"github.com/urfave/cli/v2"
)

type config struct {
	LogLevel string
}

func Run(argv []string) error {
	var cfg config
	app := &cli.App{
		Name:  "spout",
		Usage: "A log analysis tool",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "log-level",
				Aliases:     []string{"l"},
				EnvVars:     []string{"SPOUT_LOG_LEVEL"},
				Usage:       "Log level [trace|debug|info|warn|error]",
				Value:       "info",
				Destination: &cfg.LogLevel,
			},
		},
		Before: func(ctx *cli.Context) error {
			utils.Logger = utils.Logger.Clone(zlog.WithLogLevel(cfg.LogLevel))
			return nil
		},
		Commands: append(readCommands(&cfg), cmdCall(&cfg)),
	}

	if err := app.Run(argv); err != nil {
		utils.Logger.Error(err.Error())
		utils.Logger.Err(err).Debug("error detail:")
		return err
	}
	return nil
}

func readCommands(cfg *config) []*cli.Command {
	return []*cli.Command{
		cmdGCP(cfg),
		cmdFile(cfg),
	}
}
