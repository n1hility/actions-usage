/* 
Copyright (c) 2013 The go-github AUTHORS. All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are
met:

   * Redistributions of source code must retain the above copyright
notice, this list of conditions and the following disclaimer.
   * Redistributions in binary form must reproduce the above
copyright notice, this list of conditions and the following disclaimer
in the documentation and/or other materials provided with the
distribution.
   * Neither the name of Google Inc. nor the names of its
contributors may be used to endorse or promote products derived from
this software without specific prior written permission.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
"AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
(INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

/*
 * This file contains fixed up versions of the go-github project due
 * to an incomplete Action API. Therefore the license terms for this
 * file should be preserved until it can be removed entirely.
 */

package main
import (
	"context"
	"fmt"
	"net/url"
	"reflect"

	"github.com/google/go-github/v31/github"
	"github.com/google/go-querystring/query"
)

func addOptions(s string, opts interface{}) (string, error) {
	v := reflect.ValueOf(opts)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opts)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

func listRepositoryWorkflowRuns(ctx context.Context, client *github.Client, fullRepo string, opts *github.ListWorkflowRunsOptions) (*github.WorkflowRuns, *github.Response, error) {
	u := fmt.Sprintf("repos/%s/actions/runs", fullRepo)
	u, err := addOptions(u, opts)
	if err != nil {
		return nil, nil, err
	}

	req, err := client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}

	runs := new(github.WorkflowRuns)
	resp, err := client.Do(ctx, req, &runs)
	if err != nil {
		return nil, resp, err
	}

	return runs, resp, nil
}

func getWorkflow(ctx context.Context, client *github.Client, url string) (*github.Workflow, *github.Response, error) {
	req, err := client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}

	workflow := new(github.Workflow)
	resp, err := client.Do(ctx, req, workflow)
	if err != nil {
		return nil, resp, err
	}

	return workflow, resp, nil
}