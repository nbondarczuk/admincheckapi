package config

import (
	"fmt"
	"os"

	"admincheckapi/api/version"
)

func printVersionInfoAndExit() {
	fmt.Printf("%s\n", version.OneLineInfo())
	os.Exit(0)
}
