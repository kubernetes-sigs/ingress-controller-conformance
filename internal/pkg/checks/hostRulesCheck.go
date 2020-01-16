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
	hostRulesCheck.AddCheck(hostRulesExactMatchCheck)
	hostRulesCheck.AddCheck(hostRulesMatchingWildcardCheck)
	hostRulesCheck.AddCheck(hostRulesTopLevelWildcardCheck)
	hostRulesCheck.AddCheck(hostRulesMultilevelWildcardCheck)
	Checks.AddCheck(hostRulesCheck)
}

var hostRulesHost string
var hostRulesCheck = &Check{
	Name: "host-rules",
	Run: func(check *Check, config Config) (success bool, err error) {
		hostRulesHost, err = k8s.GetIngressHost("default", "host-rules")
		if err == nil {
			success = true
		}

		return
	},
}

var hostRulesExactMatchCheck = &Check{
	Name:        "host-rules-exact-match",
	Description: "Ingress with exact host rule should send traffic to the correct backend service",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s", hostRulesHost), "foo.bar.com")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "host-rules-exact", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Host, "foo.bar.com", "expected the request host would be '%s' but was '%s'")
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

var hostRulesMatchingWildcardCheck = &Check{
	Name:        "host-rules-wildcard",
	Description: "Ingress with wildcard host rule should match single-level wildcard requests",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s", hostRulesHost), "wildcard.bar.com")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "host-rules-wildcard", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Host, "wildcard.bar.com", "expected the request host would be '%s' but was '%s'")
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

var hostRulesTopLevelWildcardCheck = &Check{
	Name:        "host-rules-toplevel-wildcard",
	Description: "Ingress with wildcard host rule should not match top level requests & fallback to single-service",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s", hostRulesHost), "bar.com")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "single-service", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Host, "bar.com", "expected the request host would be '%s' but was '%s'")
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

var hostRulesMultilevelWildcardCheck = &Check{
	Name:        "host-rules-multilevel-wildcard",
	Description: "Ingress with wildcard host rule should not match multi-level wildcard requests & fallback to single-service",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s", hostRulesHost), "wildcard.foo.bar.com")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "single-service", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Host, "wildcard.foo.bar.com", "expected the request host would be '%s' but was '%s'")
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
