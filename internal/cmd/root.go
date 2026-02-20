package cmd

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/yborunov/umami-cli/internal/config"
)

type Globals struct {
	Endpoint string `help:"Umami base URL (e.g. https://analytics.example.com)" env:"UMAMI_URL"`
	Token    string `help:"API token (overrides stored config)" env:"UMAMI_TOKEN"`
	JSON     bool   `help:"Output raw JSON"`
}

type CLI struct {
	Globals

	Auth      AuthCmd      `cmd:"" help:"Authenticate and manage tokens"`
	Analytics AnalyticsCmd `cmd:"" help:"Analytics operations"`
	Teams     TeamsCmd     `cmd:"" help:"Team operations"`
	Websites  WebsitesCmd  `cmd:"" help:"Website operations"`
	Version   VersionCmd   `cmd:"" help:"Print version"`
}

func Run() int {
	cli := &CLI{}
	kctx := kong.Parse(cli,
		kong.Name("umami"),
		kong.Description("CLI for Umami Analytics API"),
		kong.UsageOnError(),
	)

	cfg, err := config.Load(cli.Endpoint, cli.Token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	ctx := &Context{
		Config: cfg,
		JSON:   cli.JSON,
	}

	if err := kctx.Run(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	return 0
}
