package infra

import (
	"net/url"
	"os/exec"
	"runtime"

	"github.com/m-mizutani/goerr"
)

type browser struct{}

func (x *browser) Open(uri *url.URL) error {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", uri.String()).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", uri.String()).Start()
	case "darwin":
		err = exec.Command("open", uri.String()).Start()
	default:
		err = goerr.New("unsupported platform to open browser")
	}

	if err != nil {
		return goerr.Wrap(err)
	}

	return nil
}
