package model

import "github.com/m-mizutani/goerr"

type (
	RunMode int
)

const (
	ConsoleMode RunMode = iota + 1
	BrowserMode
)

var runModeMap = map[string]RunMode{
	"console": ConsoleMode,
	"browser": BrowserMode,
}

func ToRunMode(mode string) (RunMode, error) {
	v, ok := runModeMap[mode]
	if !ok {
		return 0, goerr.New("invalid mode, choose 'console' or 'browser'").With("mode", mode)
	}

	return v, nil
}
