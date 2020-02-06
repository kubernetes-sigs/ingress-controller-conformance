/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Ingress conformance test harness machinery
package checks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

// Config contains test suite configuration fields
type Config struct {
}

// Check represents a test case. Checks are named, and must provide a
// description and a Run function. Checks are organized in a hierarchy.
type Check struct {
	Name        string
	Description string

	Run func(check *Check, config Config) (bool, error)

	// Child checks
	checks []*Check
	// Parent check
	parent *Check
}

// CapturedRequest contains the original HTTP request metadata as received
// by the echoserver handling the test request.
// The DownstreamServiceId field contains the TEST_ID environment variable
// value of the downstream echoserver.
type CapturedRequest struct {
	DownstreamServiceId string `json:"testId"`
	Path                string
	Host                string
	Method              string
	Proto               string
	Headers             map[string][]string
}

// CapturedResponse contains the HTTP response metadata from the echoserver
type CapturedResponse struct {
	StatusCode    int
	ContentLength int64
	Proto         string
	Headers       map[string][]string
}

// AssertionSet performs checks and accumulates assertion errors
type AssertionSet []error

func captureRequest(location string, hostOverride string) (capReq CapturedRequest, capRes CapturedResponse, err error) {
	tr := &http.Transport{
		DisableCompression: true,
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 3,
	}
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		return
	}
	if hostOverride != "" {
		req.Host = hostOverride
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(&capReq)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("unexpected response (statuscode: %d, length: %d): %s", resp.StatusCode, len(body), body)
		return
	}

	capRes = CapturedResponse{
		resp.StatusCode,
		resp.ContentLength,
		resp.Proto,
		resp.Header,
	}
	return
}

// Assert actual and expected parameters are deeply equal
func (a *AssertionSet) Equals(actual interface{}, expected interface{}, errorTemplate string) {
	if errorTemplate == "" {
		errorTemplate = "Expected '%s' but was '%s'"
	}
	if !reflect.DeepEqual(expected, actual) {
		err := fmt.Errorf(errorTemplate, expected, actual)
		*a = append(*a, err)
	}
}

// Assert the actual headers contains the expected headers key
func (a *AssertionSet) ContainsHeaders(actual map[string][]string, expected []string, errorTemplate string) {
	if errorTemplate == "" {
		errorTemplate = "Expected to contain '%s' but contained '%s'"
	}
	for _, expectedKey := range expected {
		if actual[expectedKey] == nil {
			err := fmt.Errorf(errorTemplate, expectedKey, actual)
			*a = append(*a, err)
		}
	}
}

// Assert the actual headers contains exactly the expected headers key and no more
func (a *AssertionSet) ContainsExactHeaders(actual map[string][]string, expected []string, errorTemplate string) {
	a.ContainsHeaders(actual, expected, errorTemplate)
	if errorTemplate == "" {
		errorTemplate = "Expected to only contain '%s' but contained '%s'"
	}
	if len(actual) != len(expected) {
		err := fmt.Errorf(errorTemplate, expected, actual)
		*a = append(*a, err)
	}
}

func (a *AssertionSet) Error() (err string) {
	for i, e := range *a {
		err += fmt.Sprintf("\t%d) Assertion failed: %s\n", i+1, e.Error())
	}
	return
}

// Head of Check hierarchy
var Checks = &Check{
	Name: "all",
}

// Add child checks
func (c *Check) AddCheck(checks ...*Check) {
	for i, x := range checks {
		if checks[i] == c {
			panic("Checks can't be a child of itself")
		}
		checks[i].parent = c
		c.checks = append(c.checks, x)
	}
}

// List this check and its child's description
func (c Check) List() {
	if c.Description != "" {
		fmt.Printf("- %s [%s]\n", c.Description, c.Name)
	}
	for _, check := range c.checks {
		check.List()
	}
}

// Run all checks, filtered by name and given a configuration
func (c Check) Verify(filterOnCheckName string, config Config) (successCount int, failureCount int, err error) {
	if filterOnCheckName != c.Name && filterOnCheckName != "" {
		for _, check := range c.checks {
			s, f, err := check.Verify(filterOnCheckName, config)
			successCount += s
			failureCount += f
			if err != nil {
				fmt.Printf(err.Error())
			}
		}

		return
	}

	fmt.Printf("Running '%s' verifications...\n", c.Name)
	runChildChecks := true
	if c.Run != nil {
		success, err := c.Run(&c, config)
		if err != nil {
			fmt.Printf("  %s\n", err.Error())
		}

		if success {
			successCount++
		} else {
			failureCount++
			runChildChecks = false
			fmt.Printf("  Check failed: %s\n", c.Name)
		}
	}

	if runChildChecks {
		for _, check := range c.checks {
			s, f, err := check.Verify("", config)
			if err != nil {
				fmt.Printf(err.Error())
			}
			successCount += s
			failureCount += f
		}
	}

	return
}
