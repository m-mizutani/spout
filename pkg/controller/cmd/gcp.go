package cmd

import (
	"os"
	"strings"

	"github.com/m-mizutani/gots/slice"
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
		Filter    cli.StringSlice
	}
	var periodOpt periodOptions
	var runOpt runOptions

	flags := slice.Flatten(
		[]cli.Flag{
			&cli.StringFlag{
				Name:        "project",
				Aliases:     []string{"p"},
				EnvVars:     []string{"SPOUT_GCP_PROJECT"},
				Usage:       "Google Cloud Project ID",
				Destination: (*string)(&localCfg.ProjectID),
				Required:    true,
			},
			&cli.StringSliceFlag{
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
				Value:       1000,
			},
		},
		periodOpt.Flags(),
		runOpt.Flags(),
	)

	return &cli.Command{
		Name:        "gcp",
		Aliases:     []string{"g"},
		Description: "Read Google Cloud logs",
		Flags:       flags,
		Action: func(c *cli.Context) error {
			utils.Logger.With("config", localCfg).Debug("fetching cloud logging")

			period, err := model.NewPeriod(periodOpt.baseTime, periodOpt.duration, model.RangeType(periodOpt.rangeType))
			if err != nil {
				return err
			}

			var options []gcp.Option
			filters := localCfg.Filter.Value()
			if len(filters) > 0 {
				options = append(options, gcp.WithFilter(strings.Join(filters, " ")))
			}
			reader := gcp.New(localCfg.ProjectID, localCfg.Limit, period, options...)

			clients := infra.New(
				infra.WithLogReader(reader),
				infra.WithWriter(os.Stdout),
			)
			ctx := model.NewContext(
				model.WithCtx(c.Context),
			)

			uc := usecase.New(clients)

			if err := run(ctx, uc, &runOpt); err != nil {
				return err
			}

			return nil
		},
	}
}
