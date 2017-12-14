// Copyright 2017 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package main

import (
	"os"

	"github.com/joshdk/cci-trigger/cmd"
)

func main() {
	os.Exit(cmd.Cmd().Run(os.Args))
}
