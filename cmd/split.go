// Copyright 2017 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package cmd

import (
	"fmt"
	"regexp"
	"strings"
)

func splitParams(args []string) (map[string]string, error) {
	regexBuildVar := regexp.MustCompile("^[a-zA-Z_]+[a-zA-Z0-9_]*$")

	params := make(map[string]string, len(args))

	if len(args) == 0 {
		return nil, nil
	}

	for _, arg := range args {
		chunks := strings.SplitN(arg, "=", 2)

		for index, chunk := range chunks {
			chunks[index] = strings.TrimSpace(chunk)
		}

		switch {
		case len(chunks) != 2:
			fallthrough
		case chunks[0] == "":
			fallthrough
		case chunks[1] == "":
			fallthrough
		case !regexBuildVar.MatchString(chunks[0]):
			return nil, fmt.Errorf("invalid build parameter %q", arg)

		default:
			params[chunks[0]] = chunks[1]
		}
	}

	return params, nil
}
