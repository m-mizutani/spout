package cmd

import (
	"net/url"

	"github.com/m-mizutani/spout/pkg/controller/server"
	"github.com/m-mizutani/spout/pkg/model"
	"github.com/m-mizutani/spout/pkg/usecase"
	"github.com/m-mizutani/spout/pkg/utils"
	"github.com/urfave/cli/v2"
	"golang.org/x/sync/errgroup"
)

type periodOptions struct {
	baseTime  string
	duration  string
	rangeType string
}

type runOptions struct {
	mode string
	addr string
}

func (x *periodOptions) Flags() []cli.Flag {
	return []cli.Flag{
		// for period
		&cli.StringFlag{
			Name:        "base-time",
			Aliases:     []string{"t"},
			EnvVars:     []string{"SPOUT_BASE_TIME"},
			Usage:       "Base time",
			Destination: &x.baseTime,
		},
		&cli.StringFlag{
			Name:        "duration",
			Aliases:     []string{"d"},
			EnvVars:     []string{"SPOUT_DURATION"},
			Usage:       "Duration, e.g. 10m, 30s",
			Value:       "10m",
			Destination: &x.duration,
		},
		&cli.StringFlag{
			Name:        "range",
			Aliases:     []string{"r"},
			EnvVars:     []string{"SPOUT_RANGE"},
			Usage:       "Range type [before|after|around]",
			Value:       "before",
			Destination: &x.rangeType,
		},
	}
}

func (x *runOptions) Flags() []cli.Flag {
	return []cli.Flag{
		// mode
		&cli.StringFlag{
			Name:        "mode",
			Aliases:     []string{"m"},
			EnvVars:     []string{"SPOUT_MODE"},
			Usage:       "Run mode [console|browser]",
			Value:       "browser",
			Destination: &x.mode,
		},
		&cli.StringFlag{
			Name:        "addr",
			Aliases:     []string{"a"},
			EnvVars:     []string{"SPOUT_ADDR"},
			Usage:       "Server address for browser mode",
			Value:       "127.0.0.1:3280",
			Destination: &x.addr,
		},
	}
}

func run(ctx *model.Context, uc *usecase.Usecase, opt *runOptions) error {
	mode, err := model.ToRunMode(opt.mode)
	if err != nil {
		return err
	}

	switch mode {
	case model.ConsoleMode:
		if err := uc.DumpLogs(ctx); err != nil {
			return err
		}

	case model.BrowserMode:
		errgp, ectx := errgroup.WithContext(ctx)
		ctx = ctx.New(model.WithCtx(ectx))
		uri := &url.URL{
			Scheme: "http",
			Host:   opt.addr,
		}

		errgp.Go(func() error {
			if err := uc.ImportLogs(ctx); err != nil {
				utils.Logger.With("error", err.Error()).Error("failed to import logs")
			}
			return nil
		})

		errgp.Go(func() error {
			utils.Logger.Info("starting server %s", uri.String())
			if err := server.New(uc).Listen(opt.addr); err != nil {
				utils.Logger.With("error", err.Error()).Error("failed to import logs")
				return err
			}
			return nil
		})

		errgp.Go(func() error {
			if err := uc.OpenBrowser(uri); err != nil {
				utils.Logger.With("error", err.Error()).Error("failed to import logs")
			}
			return nil
		})

		if err := errgp.Wait(); err != nil {
			return err
		}
	}

	return nil
}
