package out

import (
	"encoding/json"
	"fmt"
	"os"
)

func PrintJSON(v any) error {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func Printf(format string, args ...any) {
	fmt.Printf(format, args...)
}
