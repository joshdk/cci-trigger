package main

import (
	"os"

	"github.com/joshdk/cci-trigger/cmd"
)

func main() {
	os.Exit(cmd.Cmd().Run(os.Args))
}
