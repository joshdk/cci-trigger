// Copyright 2017 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package cmd

var version string

func Version() string {
	if version == "" {
		return "unknown"
	}

	return version
}
