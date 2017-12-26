// Copyright 2017 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package cmd

import (
	"errors"
	"fmt"

	"github.com/joshdk/cci-trigger/cci"
)

type action uint

type handler func(client cci.Client, vcs string, username string, project string) (*cci.BuildResponse, error)

const (
	unknown action = iota
	buildDefault
	buildBranch
	buildBranchAtRef
	buildRef
	buildTag
	rebuild
	rebuildWithSSH
)

// getAction converts the given flags into the correct action to take, if possible.
func getAction(build string, ssh bool, tag string, branch string, ref string, params map[string]string) (action, error) {
	type flagSet [6]bool

	has := flagSet{
		branch != "",
		ref != "",
		tag != "",
		build != "",
		ssh,
		len(params) != 0,
	}

	switch has {

	// cci-trigger <project> [params...]
	case flagSet{false, false, false, false, false, false}:
		fallthrough
	case flagSet{false, false, false, false, false, true}:
		return buildDefault, nil

	// cci-trigger <project> --branch X [params...]
	case flagSet{true, false, false, false, false, false}:
		fallthrough
	case flagSet{true, false, false, false, false, true}:
		return buildBranch, nil

	// cci-trigger <project> --branch X --ref X [params...]
	case flagSet{true, true, false, false, false, false}:
		fallthrough
	case flagSet{true, true, false, false, false, true}:
		return buildBranchAtRef, nil

	// cci-trigger <project> --ref X [params...]
	case flagSet{false, true, false, false, false, false}:
		fallthrough
	case flagSet{false, true, false, false, false, true}:
		return buildRef, nil

	// cci-trigger <project> --tag X [params...]
	case flagSet{false, false, true, false, false, false}:
		fallthrough
	case flagSet{false, false, true, false, false, true}:
		return buildTag, nil

	// cci-trigger <project> --build X
	case flagSet{false, false, false, true, false, false}:
		return rebuild, nil

	// cci-trigger <project> --build X --ssh
	case flagSet{false, false, false, true, true, false}:
		return rebuildWithSSH, nil

	default:
		return unknown, errors.New("invalid flag combination")
	}
}

func getHandler(action action, build string, ssh bool, tag string, branch string, ref string, params map[string]string) (string, handler) {
	switch action {
	case buildDefault:
		return "build default branch",
			func(client cci.Client, vcs string, username string, project string) (*cci.BuildResponse, error) {
				return client.BuildDefault(vcs, username, project, params)
			}
	case buildBranch:
		return fmt.Sprintf("build branch %s", branch),
			func(client cci.Client, vcs string, username string, project string) (*cci.BuildResponse, error) {
				return client.BuildBranch(vcs, username, project, branch, params)
			}
	case buildBranchAtRef:
		return fmt.Sprintf("build branch %s at %s", branch, ref),
			func(client cci.Client, vcs string, username string, project string) (*cci.BuildResponse, error) {
				return client.BuildBranchAtRef(vcs, username, project, branch, ref, params)
			}
	case buildRef:
		return fmt.Sprintf("build ref %s", ref),
			func(client cci.Client, vcs string, username string, project string) (*cci.BuildResponse, error) {
				return client.BuildRef(vcs, username, project, ref, params)
			}
	case buildTag:
		return fmt.Sprintf("build tag %s", tag),
			func(client cci.Client, vcs string, username string, project string) (*cci.BuildResponse, error) {
				return client.BuildTag(vcs, username, project, tag, params)
			}
	case rebuild:
		return fmt.Sprintf("rebuild #%s", build),
			func(client cci.Client, vcs string, username string, project string) (*cci.BuildResponse, error) {
				return client.Rebuild(vcs, username, project, build)
			}
	case rebuildWithSSH:
		return fmt.Sprintf("rebuild #%s with SSH", build),
			func(client cci.Client, vcs string, username string, project string) (*cci.BuildResponse, error) {
				return client.RebuildWithSSH(vcs, username, project, build)
			}
	default:
		return "invalid flag combination", nil
	}
}
