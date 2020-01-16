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
	"fmt"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/k8s"
)

func init() {
	pathRulesCheck.AddCheck(pathRulesFooCheck)
	pathRulesCheck.AddCheck(pathRulesFooTrailingSlashCheck)
	pathRulesCheck.AddCheck(pathRulesBarCheck)
	pathRulesCheck.AddCheck(pathRulesBarSubpathCheck)
	Checks.AddCheck(pathRulesCheck)
}

var host string

var pathRulesCheck = &Check{
	Name: "path-rules",
	Run: func(check *Check, config Config) (success bool, err error) {
		host, err = k8s.GetIngressHost("default", "path-rules")
		if err == nil {
			success = true
		}

		return
	},
}

var pathRulesFooCheck = &Check{
	Name:        "path-rules-foo",
	Description: "[SAMPLE] Ingress with path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request path",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/foo", host), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-foo", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/foo", "expected the request path would be '%s' but was '%s'")
		// Assert the downstream service response
		a.equals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}

var pathRulesFooTrailingSlashCheck = &Check{
	Name:        "path-rules-foo-trailing",
	Description: "[SAMPLE] Ingress with path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request including sub-paths",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/foo/", host), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-foo", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/foo/", "expected the request path would be '%s' but was '%s'")
		// Assert the downstream service response
		a.equals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}

var pathRulesBarCheck = &Check{
	Name:        "path-rules-bar",
	Description: "[SAMPLE] Ingress with path rule with a trailing slash should send traffic to the correct backend service, and preserve the original request path",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/bar/", host), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-bar", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/bar/", "expected the request path would be '%s' but was '%s'")
		// Assert the downstream service response
		a.equals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}

var pathRulesBarSubpathCheck = &Check{
	Name:        "path-rules-bar-subpath",
	Description: "[SAMPLE] Ingress with path rule with a trailing slash should send traffic to the correct backend service, and preserve the original request including sub-paths and double '/'",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/bar//bershop", host), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-bar", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/bar//bershop", "expected the request path would be '%s' but was '%s'")
		// Assert the downstream service response
		a.equals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}

// TODO: Implement more checks on edge cases like leading `/`, query params, and encoding
