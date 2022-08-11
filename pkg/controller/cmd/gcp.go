package cmd

import (
	"os"

	"github.com/m-mizutani/spout/pkg/infra"
	"github.com/m-mizutani/spout/pkg/infra/gcp"
	"github.com/m-mizutani/spout/pkg/model"
	"github.com/m-mizutani/spout/pkg/usecase"
	"github.com/m-mizutani/spout/pkg/utils"
	"github.com/urfave/cli/v2"
)

func cmdGCP(globalCfg *config) *cli.Command {
	var localCfg struct {
		ProjectID model.GoogleProjectID
		Limit     int
		Filter    string
	}
	var commonOpt commonOptions

	flags := append([]cli.Flag{
		&cli.StringFlag{
			Name:        "project",
			Aliases:     []string{"p"},
			EnvVars:     []string{"SPOUT_GCP_PROJECT"},
			Usage:       "Google Cloud Project ID",
			Destination: (*string)(&localCfg.ProjectID),
			Required:    true,
		},
		&cli.StringFlag{
			Name:        "filter",
			Aliases:     []string{"f"},
			EnvVars:     []string{"SPOUT_GCP_FILTER"},
			Usage:       "Google Cloud Logging filter",
			Destination: &localCfg.Filter,
		},
		&cli.IntFlag{
			Name:        "limit",
			EnvVars:     []string{"SPOUT_GCP_LIMIT"},
			Usage:       "Limit of fetching log",
			Destination: &localCfg.Limit,
			Value:       100,
		},
	}, commonOpt.Flags()...)

	return &cli.Command{
		Name:        "gcp",
		Aliases:     []string{"g"},
		Description: "Read Google Cloud logs",
		Flags:       flags,
		Action: func(c *cli.Context) error {
			utils.Logger.With("config", localCfg).Debug("fetching cloud logging")

			period, err := model.NewPeriod(commonOpt.baseTime, commonOpt.duration, model.RangeType(commonOpt.rangeType))
			if err != nil {
				return err
			}

			var options []gcp.Option
			if localCfg.Filter != "" {
				options = append(options, gcp.WithFilter(localCfg.Filter))
			}
			reader := gcp.New(localCfg.ProjectID, localCfg.Limit, period, options...)

			clients := infra.New(
				infra.WithLogReader(reader),
				infra.WithWriter(os.Stdout),
			)
			ctx := model.NewContext(
				model.WithCtx(c.Context),
			)

			if err := usecase.New(clients).DumpLogs(ctx); err != nil {
				return err
			}

			return nil
		},
	}
}
