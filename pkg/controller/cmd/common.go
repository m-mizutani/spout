package cmd

import (
	"github.com/urfave/cli/v2"
)

type commonOptions struct {
	baseTime  string
	duration  string
	rangeType string

	mode string
	addr string
}

func (x *commonOptions) Flags() []cli.Flag {
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
