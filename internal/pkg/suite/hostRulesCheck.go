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

package suite

import (
	"fmt"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/checks"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/k8s"
)

func init() {
	hostRulesCheck.AddCheck(hostRulesExactMatchCheck)
	hostRulesCheck.AddCheck(hostRulesWildcardSingleLabelCheck)
	hostRulesCheck.AddCheck(hostRulesWildcardMultipleLabelsCheck)
	hostRulesCheck.AddCheck(hostRulesWildcardNoLabelCheck)
	checks.AllChecks.AddCheck(hostRulesCheck)
}

var hostRulesHost string
var hostRulesCheck = &checks.Check{
	Name: "host-rules",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		if config.UseInsecureHost != "" {
			hostRulesHost = config.UseInsecureHost
			return true, nil
		}
		var err error
		hostRulesHost, err = k8s.GetIngressHost("default", "host-rules")
		if err != nil {
			return false, err
		}
		return true, nil
	},
}

var hostRulesExactMatchCheck = &checks.Check{
	Name:        "host-rules-exact-match",
	Description: "Ingress with exact host rule should send traffic to the correct backend service",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s", hostRulesHost), "foo.bar.com")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "host-rules-exact", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Host, "foo.bar.com", "expected the request host would be '%s' but was '%s'")
		// Assert the downstream service response
		a.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

		if a.Error() == "" {
			return true, nil
		} else {
			fmt.Print(a)
		}
		return false, nil
	},
}

var hostRulesWildcardSingleLabelCheck = &checks.Check{
	Name:        "host-rules-wildcard-single-label",
	Description: "Ingress with wildcard host rule should match a single label",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s", hostRulesHost), "wildcard.foo.com")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "host-rules-wildcard", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Host, "wildcard.foo.com", "expected the request host would be '%s' but was '%s'")
		// Assert the downstream service response
		a.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

		if a.Error() == "" {
			return true, nil
		} else {
			fmt.Print(a)
		}
		return false, nil
	},
}

var hostRulesWildcardMultipleLabelsCheck = &checks.Check{
	Name:        "host-rules-wildcard-multiple-labels",
	Description: "Ingress with wildcard host rule should only match a single label & fallback to default-backend",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s", hostRulesHost), "aaa.bbb.foo.com")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "default-backend", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Host, "aaa.bbb.foo.com", "expected the request host would be '%s' but was '%s'")
		// Assert the downstream service response
		a.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

		if a.Error() == "" {
			return true, nil
		} else {
			fmt.Print(a)
		}
		return false, nil
	},
}

var hostRulesWildcardNoLabelCheck = &checks.Check{
	Name:        "host-rules-wildcard-no-label",
	Description: "Ingress with wildcard host rule should match exactly one single label & fallback to default-backend",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s", hostRulesHost), "foo.com")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "default-backend", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Host, "foo.com", "expected the request host would be '%s' but was '%s'")
		// Assert the downstream service response
		a.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

		if a.Error() == "" {
			return true, nil
		} else {
			fmt.Print(a)
		}
		return false, nil
	},
}
