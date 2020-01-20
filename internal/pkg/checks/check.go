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

package checks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"
)

type Config struct {
}

type Check struct {
	Name        string
	Description string

	Run func(check *Check, config Config) (bool, error)

	// All checks
	checks []*Check

	// Parent check
	parent *Check
}

type CapturedRequest struct {
	TestId  string
	Path    string
	Host    string
	Method  string
	Proto   string
	Headers map[string][]string
}

type CapturedResponse struct {
	StatusCode    int
	ContentLength int64
	Proto         string
	Headers       map[string][]string
}

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

type assertionSet []error

func (a *assertionSet) equals(actual interface{}, expected interface{}, errorTemplate string) {
	if errorTemplate == "" {
		errorTemplate = "Expected '%s' but was '%s'"
	}
	if !reflect.DeepEqual(expected, actual) {
		err := fmt.Errorf(errorTemplate, expected, actual)
		*a = append(*a, err)
	}
}

func (a *assertionSet) containsKeys(actual map[string][]string, expected []string, errorTemplate string) {
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

func (a *assertionSet) containsOnlyKeys(actual map[string][]string, expected []string, errorTemplate string) {
	a.containsKeys(actual, expected, errorTemplate)
	if errorTemplate == "" {
		errorTemplate = "Expected to only contain '%s' but contained '%s'"
	}
	if len(actual) != len(expected) {
		err := fmt.Errorf(errorTemplate, expected, actual)
		*a = append(*a, err)
	}
}

func (a *assertionSet) Error() (err string) {
	for i, e := range *a {
		err += fmt.Sprintf("\t%d) Assertion failed: %s\n", i+1, e.Error())
	}
	return
}

func (c *Check) AddCheck(checks ...*Check) {
	for i, x := range checks {
		if checks[i] == c {
			panic("Checks can't be a child of itself")
		}
		checks[i].parent = c
		c.checks = append(c.checks, x)
	}
}

var Checks = &Check{
	Name: "all",
}

func (c Check) List() {
	if c.Description != "" {
		fmt.Printf("- %s (%s)\n", c.Description, c.Name)
	}
	for _, check := range c.checks {
		check.List()
	}
}

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
