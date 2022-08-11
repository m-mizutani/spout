package cmd

import (
	"github.com/urfave/cli/v2"
)

type commonOptions struct {
	baseTime  string
	duration  string
	rangeType string
	query     string
}

func (x *commonOptions) Flags() []cli.Flag {
	return []cli.Flag{
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
		&cli.StringFlag{
			Name:        "jq",
			Aliases:     []string{"j"},
			EnvVars:     []string{"SPOUT_JQ"},
			Usage:       "jq ",
			Destination: &x.query,
		},
	}
}
