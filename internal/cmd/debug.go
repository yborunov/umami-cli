package cmd

import (
	"os"
	"strings"
)

func debugEnabled() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv("DEBUG")), "true")
}
