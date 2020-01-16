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
	pathRulesCheck.AddCheck(prefixPathRulesFooCheck)
	pathRulesCheck.AddCheck(prefixPathRulesFooTrailingSlashCheck)
	pathRulesCheck.AddCheck(prefixPathRulesFooSubpathCheck)
	pathRulesCheck.AddCheck(prefixPathRulesFooNomatchPrefixCheck)
	pathRulesCheck.AddCheck(prefixPathRulesBarCheck)
	pathRulesCheck.AddCheck(prefixPathRulesBarTrailingSlashCheck)
	pathRulesCheck.AddCheck(prefixPathRulesBarSubpathCheck)
	pathRulesCheck.AddCheck(prefixPathRulesBarNomatchPrefixCheck)
	Checks.AddCheck(pathRulesCheck)
}

var pathRulesHost string
var pathRulesCheck = &Check{
	Name: "path-rules",
	Run: func(check *Check, config Config) (success bool, err error) {
		pathRulesHost, err = k8s.GetIngressHost("default", "path-rules")
		if err == nil {
			success = true
		}

		return
	},
}

var prefixPathRulesFooCheck = &Check{
	Name:        "prefix-path-rules-foo",
	Description: "Ingress with prefix path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request path",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/foo", pathRulesHost), "")
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

var prefixPathRulesFooTrailingSlashCheck = &Check{
	Name:        "prefix-path-rules-foo-trailing",
	Description: "Ingress with prefix path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request including sub-paths",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/foo/", pathRulesHost), "")
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

var prefixPathRulesFooSubpathCheck = &Check{
	Name:        "prefix-path-rules-foo-subpath",
	Description: "Ingress with prefix path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/foo/bar", pathRulesHost), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-foo", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/foo/bar", "expected the request path would be '%s' but was '%s'")
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

var prefixPathRulesFooNomatchPrefixCheck = &Check{
	Name:        "prefix-path-rules-foo-nomatch",
	Description: "Ingress with prefix path rule with a trailing slash should not match on a partial path and fallback to catch-all",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/foobar", pathRulesHost), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/foobar", "expected the request path would be '%s' but was '%s'")
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

var prefixPathRulesBarCheck = &Check{
	Name:        "prefix-path-rules-bar",
	Description: "Ingress with prefix path rule with a trailing slash should send traffic to the correct backend service, and preserve the original request path",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/bar/", pathRulesHost), "")
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

var prefixPathRulesBarTrailingSlashCheck = &Check{
	Name:        "prefix-path-rules-bar-trailing-slash-ignored",
	Description: "Ingress with prefix path rule with a trailing slash is ignored and should send traffic to the correct backend service, and preserve the original request path",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/bar", pathRulesHost), "")
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

var prefixPathRulesBarSubpathCheck = &Check{
	Name:        "prefix-path-rules-bar-subpath",
	Description: "Ingress with prefix path rule with a trailing slash should send traffic to the correct backend service, and preserve the original request including sub-paths and double '/'",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/bar//bershop", pathRulesHost), "")
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

var prefixPathRulesBarNomatchPrefixCheck = &Check{
	Name:        "prefix-path-rules-bar-nomatch",
	Description: "Ingress with prefix path rule with a trailing slash should not match on a partial path and fallback to catch-all",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/barbershop", pathRulesHost), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/barbershop", "expected the request path would be '%s' but was '%s'")
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
