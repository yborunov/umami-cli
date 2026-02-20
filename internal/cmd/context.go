package cmd

import "github.com/yborunov/umami-cli/internal/config"

type Context struct {
	Config *config.Config
	JSON   bool
}
