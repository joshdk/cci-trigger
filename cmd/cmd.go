// Copyright 2017 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package cmd

import (
	"fmt"
	"os"

	"github.com/palantir/pkg/cli"
	"github.com/palantir/pkg/cli/flag"
)

const (
	CircleTokenEnvVar = "CIRCLE_TOKEN"
	CircleHostEnvVar  = "CIRCLE_HOST"
	CirclePublicHost  = "circleci.com"
)

var (
	projectParam = flag.StringParam{
		Name: "project",
	}
	buildFlag = flag.StringFlag{
		Name: "build",
	}
	sshFlag = flag.BoolFlag{
		Name: "ssh",
	}
	tagFlag = flag.StringFlag{
		Name: "tag",
	}
	branchFlag = flag.StringFlag{
		Name: "branch",
	}
	refFlag = flag.StringFlag{
		Name: "ref",
	}
	buildParams = flag.StringSlice{
		Name:     "params",
		Optional: true,
	}
)

func Cmd() *cli.App {

	app := cli.NewApp()

	app.Name = "cci-trigger"
	app.Description = "Trigger CircleCI builds programmatically"

	app.Flags = []flag.Flag{
		projectParam,
		buildFlag,
		sshFlag,
		tagFlag,
		branchFlag,
		refFlag,
		buildParams,
	}

	app.ErrorHandler = func(ctx cli.Context, err error) int {
		fmt.Fprintf(os.Stderr, "%s: %s\n", app.Name, err.Error())
		return 1
	}

	app.Action = func(ctx cli.Context) error {

		var (
			token   string
			host    string
			project = ctx.String(projectParam.Name)
			branch  = ctx.String(branchFlag.Name)
			ref     = ctx.String(refFlag.Name)
			tag     = ctx.String(tagFlag.Name)
			build   = ctx.String(buildFlag.Name)
			ssh     = ctx.Bool(sshFlag.Name)
			params  = ctx.Slice(buildParams.Name)
		)

		// The CIRCLE_TOKEN environment variable is required for operation
		token, found := os.LookupEnv(CircleTokenEnvVar)
		if !found {
			return fmt.Errorf("no %s in working environment", CircleTokenEnvVar)
		}

		// The CIRCLE_HOST environment variable is optional, and overrides the default
		host, found = os.LookupEnv(CircleHostEnvVar)
		if !found {
			host = CirclePublicHost
		}

		fmt.Printf("host:    %q\n", host)
		fmt.Printf("token:   %q\n", token)
		fmt.Printf("project: %q\n", project)
		fmt.Printf("build:   %q\n", build)
		fmt.Printf("ssh:     %t\n", ssh)
		fmt.Printf("tag:     %q\n", tag)
		fmt.Printf("branch:  %q\n", branch)
		fmt.Printf("ref:     %q\n", ref)
		fmt.Printf("params:  %v\n", params)

		return nil
	}

	return app
}
