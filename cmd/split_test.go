// Copyright 2017 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package cmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSplitParams(t *testing.T) {

	tests := []struct {
		title  string
		args   []string
		params map[string]string
		err    string
	}{
		{
			title: "no params",
		},
		{
			title: "single param",
			args:  []string{"key=value"},
			params: map[string]string{
				"key": "value",
			},
		},
		{
			title: "multiple params",
			args: []string{
				"key1=value1",
				"key2=value2",
				"key3=value3",
			},
			params: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
			},
		},
		{
			title: "duplicate params",
			args: []string{
				"key1=value1",
				"key2=value2",
				"key3=value3",
				"key2=value4",
			},
			params: map[string]string{
				"key1": "value1",
				"key2": "value4",
				"key3": "value3",
			},
		},
		{
			title: "blank param",
			args:  []string{""},
			err:   `invalid build parameter ""`,
		},
		{
			title: "single equals",
			args:  []string{"="},
			err:   `invalid build parameter "="`,
		},
		{
			title: "single equals whitespace",
			args:  []string{" = "},
			err:   `invalid build parameter " = "`,
		},
		{
			title: "single equals prefix",
			args:  []string{"key="},
			err:   `invalid build parameter "key="`,
		},
		{
			title: "single equals prefix whitespace",
			args:  []string{" key= "},
			err:   `invalid build parameter " key= "`,
		},
		{
			title: "single equals suffix",
			args:  []string{"=value"},
			err:   `invalid build parameter "=value"`,
		},
		{
			title: "single equals suffix whitespace",
			args:  []string{" =value "},
			err:   `invalid build parameter " =value "`,
		},
		{
			title: "double equals consecutive",
			args:  []string{"key==value"},
			params: map[string]string{
				"key": "=value",
			},
		},
		{
			title: "double equals",
			args:  []string{"key=value=foo"},
			params: map[string]string{
				"key": "value=foo",
			},
		},
		{
			title: "whitespace",
			args: []string{
				"key1   =value1",
				"key2=   value2",
				"key3   =   value3",
				"key4=value4   ",
				"   key5=value5",
				"   key6=value6   ",
				"   key7   =   value7   ",
				"   key8   =   value   8   ",
			},
			params: map[string]string{
				"key1": "value1",
				"key2": "value2",
				"key3": "value3",
				"key4": "value4",
				"key5": "value5",
				"key6": "value6",
				"key7": "value7",
				"key8": "value   8",
			},
		},
		{
			title: "key invalid character",
			args:  []string{"key.custom=value"},
			err:   `invalid build parameter "key.custom=value"`,
		},
		{
			title: "key digit prefix",
			args:  []string{"1key=value"},
			err:   `invalid build parameter "1key=value"`,
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {
			actual, err := splitParams(test.args)

			if test.err != "" {
				require.EqualError(t, err, test.err)
			}

			require.Equal(t, test.params, actual)
		})
	}
}

func TestSplitProject(t *testing.T) {

	tests := []struct {
		title    string
		arg      string
		vcs      string
		username string
		project  string
		err      string
	}{
		{
			title: "empty",
			err:   `invalid project name ""`,
		},
		{
			title: "single field",
			arg:   "example",
			err:   `invalid project name "example"`,
		},
		{
			title:    "two fields",
			arg:      "alice/example",
			vcs:      "github",
			username: "alice",
			project:  "example",
		},
		{
			title:    "short github vcs",
			arg:      "gh/alice/example",
			vcs:      "github",
			username: "alice",
			project:  "example",
		},
		{
			title:    "long github vcs",
			arg:      "github/alice/example",
			vcs:      "github",
			username: "alice",
			project:  "example",
		},
		{
			title:    "short bitbucket vcs",
			arg:      "bb/bob/example",
			vcs:      "bitbucket",
			username: "bob",
			project:  "example",
		},
		{
			title:    "long bitbucket vcs",
			arg:      "bitbucket/bob/example",
			vcs:      "bitbucket",
			username: "bob",
			project:  "example",
		},
		{
			title: "unknown vcs",
			arg:   "svn/carol/example",
			err:   `invalid project name "svn/carol/example"`,
		},
		{
			title:    "many slashes",
			arg:      "github/carol/example/a/b/c",
			vcs:      "github",
			username: "carol",
			project:  "example/a/b/c",
		},
		{
			title:    "other characters",
			arg:      "dave-user/example_repo",
			vcs:      "github",
			username: "dave-user",
			project:  "example_repo",
		},
	}

	for index, test := range tests {
		name := fmt.Sprintf("Case #%d - %s", index, test.title)

		t.Run(name, func(t *testing.T) {
			vcs, username, project, err := splitProject(test.arg)

			if test.err != "" {
				require.EqualError(t, err, test.err)
			}

			require.Equal(t, test.vcs, vcs)
			require.Equal(t, test.username, username)
			require.Equal(t, test.project, project)
		})
	}
}
