package cmd

import (
	"runtime"

	"github.com/yborunov/umami-cli/internal/out"
)

var version = "dev"

var commit = "none"

var date = "unknown"

type VersionCmd struct{}

func (c *VersionCmd) Run(*Context) error {
	out.Printf("umami version %s\n", version)
	out.Printf("commit: %s\n", commit)
	out.Printf("built: %s\n", date)
	out.Printf("go: %s\n", runtime.Version())
	return nil
}
