package main

import (
	"os"

	"github.com/yborunov/umami-cli/internal/cmd"
)

func main() {
	os.Exit(cmd.Run())
}
