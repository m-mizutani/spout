//go:build !github_test
// +build !github_test

package frontend

import "embed"

//go:embed out/*
//go:embed out/_next
//go:embed out/_next/static/*/*.js
//go:embed out/_next/static/chunks/pages/*.js
var assets embed.FS

func Assets() *embed.FS {
	return &assets
}
