//go:build github_test
// +build github_test

package frontend

import "embed"

var assets embed.FS

func Assets() *embed.FS {
	return &assets
}
