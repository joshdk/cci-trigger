// Copyright 2017 Josh Komoroske. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE.txt file.

package cci

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	// PublicHostname is the public CircleCI service endpoint
	PublicHostname = "circleci.com"
)

type Client struct {
	token string
	host  string
}

type BuildResponse struct {
	BuildURL string `json:"build_url"`
}

func New(token string) Client {
	return Client{token, PublicHostname}
}

func NewWithHost(token string, host string) Client {
	return Client{token, host}
}

// BuildDefault triggers a build on the HEAD of the default branch. This branch
// is typically master, and can be customized in your VCS platform.
//
// See https://circleci.com/docs/api/v1-reference/#new-build for details on
// this API action.
func (client Client) BuildDefault(vcs string, username string, project string, params map[string]string) (*BuildResponse, error) {
	// https://circleci.com/api/v1.1/project/:vcs-type/:username/:project
	path := fmt.Sprintf("project/%s/%s/%s", vcs, username, project)

	return client.do(path, "", "", params)
}

// BuildTag triggers a build on the given tag.
//
// See https://circleci.com/docs/api/v1-reference/#new-build for details on
// this API action.
func (client Client) BuildTag(vcs string, username string, project string, tag string, params map[string]string) (*BuildResponse, error) {
	// https://circleci.com/api/v1.1/project/:vcs-type/:username/:project
	path := fmt.Sprintf("project/%s/%s/%s", vcs, username, project)

	return client.do(path, tag, "", params)
}

// BuildRef triggers a build on the given ref.
//
// See https://circleci.com/docs/api/v1-reference/#new-build for details on
// this API action.
func (client Client) BuildRef(vcs string, username string, project string, ref string, params map[string]string) (*BuildResponse, error) {
	// https://circleci.com/api/v1.1/project/:vcs-type/:username/:project
	path := fmt.Sprintf("project/%s/%s/%s", vcs, username, project)

	return client.do(path, "", ref, params)
}

// BuildBranch triggers a build on the HEAD of the given branch.
//
// See https://circleci.com/docs/api/v1-reference/#new-build-branch for details
// on this API action.
func (client Client) BuildBranch(vcs string, username string, project string, branch string, params map[string]string) (*BuildResponse, error) {
	// https://circleci.com/api/v1.1/project/:vcs-type/:username/:project/tree/:branch
	path := fmt.Sprintf("project/%s/%s/%s/tree/%s", vcs, username, project, branch)

	return client.do(path, "", "", params)
}

// BuildBranchAtRef triggers a build on the given branch at the given ref.
//
// See https://circleci.com/docs/api/v1-reference/#new-build-branch for details
// on this API action.
func (client Client) BuildBranchAtRef(vcs string, username string, project string, branch string, ref string, params map[string]string) (*BuildResponse, error) {
	// https://circleci.com/api/v1.1/project/:vcs-type/:username/:project/tree/:branch
	path := fmt.Sprintf("project/%s/%s/%s/tree/%s", vcs, username, project, branch)

	return client.do(path, "", ref, params)
}

// Rebuild triggers a rebuild on the given build number.
//
// See https://circleci.com/docs/api/v1-reference/#retry-build for details on
// this API action.
func (client Client) Rebuild(vcs string, username string, project string, build string) (*BuildResponse, error) {
	// https://circleci.com/api/v1.1/project/:vcs-type/:username/:project/:build_num/retry
	path := fmt.Sprintf("project/%s/%s/%s/%s/retry", vcs, username, project, build)

	return client.do(path, "", "", nil)
}

// RebuildWithSSH triggers a rebuild on the given build number, and enables SSH.
//
// See https://circleci.com/docs/api/v1-reference/#retry-build for details on
// this API action.
func (client Client) RebuildWithSSH(vcs string, username string, project string, build string) (*BuildResponse, error) {
	// https://circleci.com/api/v1.1/project/:vcs-type/:username/:project/:build_num/ssh
	path := fmt.Sprintf("project/%s/%s/%s/%s/ssh", vcs, username, project, build)

	return client.do(path, "", "", nil)
}

func (client Client) do(path string, tag string, revision string, buildParams map[string]string) (*BuildResponse, error) {

	url := fmt.Sprintf("https://%s/api/v1.1/%s", client.host, path)

	var postParams = struct {
		Tag         string            `json:"tag,omitempty"`
		Revision    string            `json:"revision,omitempty"`
		BuildParams map[string]string `json:"build_parameters,omitempty"`
	}{
		tag,
		revision,
		buildParams,
	}

	postBody, err := json.Marshal(postParams)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(postBody))
	if err != nil {
		return nil, err
	}

	// Indicate that we are sending (and want to receive) JSON
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Include the api token as a URL parameter (...?circle-token=xxx)
	q := req.URL.Query()
	q.Add("circle-token", client.token)
	req.URL.RawQuery = q.Encode()

	// Perform the HTTP POST
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}

	// Cleanup function
	defer func() {
		if err := resp.Body.Close(); err != nil {
			panic(err)
		}
	}()

	// Return error if request was not "successful"
	if resp.StatusCode < 200 || 300 <= resp.StatusCode {
		return nil, errors.New(resp.Status)
	}

	// Read the entire response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var br BuildResponse
	if err := json.Unmarshal(body, &br); err != nil {
		return nil, err
	}

	return &br, nil
}
