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
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/apiversion"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/checks"
)

func init() {
	hostRulesCheck.AddCheck(hostRulesExactMatchCheck)
	hostRulesCheck.AddCheck(hostRulesExactNoMatchCheck)
	hostRulesCheck.AddCheck(hostRulesWildcardSingleLabelCheck)
	hostRulesCheck.AddCheck(hostRulesWildcardMultipleLabelsCheck)
	hostRulesCheck.AddCheck(hostRulesWildcardNoLabelCheck)
	checks.AllChecks.AddCheck(hostRulesCheck)
}

// placeholder check for the host-rules checks hierarchy
var hostRulesCheck = &checks.Check{
	Name: "host-rules",
}

var hostRulesExactMatchCheck = &checks.Check{
	Name:        "host-rules-exact-match",
	Description: "Ingress with exact host rule should send traffic to the correct backend service (foo.bar.com matches foo.bar.com)",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "host-rules",
		Path:        "",
		Hostname:    "foo.bar.com",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "host-rules-exact", "expected the downstream service would be '%s' but was '%s'")
			a.E.DeepEquals(req.Host, "foo.bar.com", "expected the request host would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var hostRulesExactNoMatchCheck = &checks.Check{
	Name:        "host-rules-exact-no-match",
	Description: "Ingress with exact host rule should send traffic to fallback to default-backend when hostname does not match (foo.bar.com does not match subdomain.bar.com)",
	APIVersions: apiversion.ExtensionsV1Beta1,
	RunRequest: &checks.Request{
		IngressName: "host-rules",
		Path:        "",
		Hostname:    "subdomain.bar.com",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "default-backend", "expected the downstream service would be '%s' but was '%s'")
			a.E.DeepEquals(req.Host, "subdomain.bar.com", "expected the request host would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var hostRulesWildcardSingleLabelCheck = &checks.Check{
	Name:        "host-rules-wildcard-single-label",
	Description: "Ingress with wildcard host rule should match a single label (*.foo.com matches wildcard.foo.com)",
	APIVersions: apiversion.NetworkingV1Beta1,
	RunRequest: &checks.Request{
		IngressName: "host-rules",
		Path:        "",
		Hostname:    "wildcard.foo.com",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "host-rules-wildcard", "expected the downstream service would be '%s' but was '%s'")
			a.E.DeepEquals(req.Host, "wildcard.foo.com", "expected the request host would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var hostRulesWildcardMultipleLabelsCheck = &checks.Check{
	Name:        "host-rules-wildcard-multiple-labels",
	Description: "Ingress with wildcard host rule should only match a single label & fallback to default-backend (*.foo.com does not match aaa.bbb.foo.com)",
	APIVersions: apiversion.NetworkingV1Beta1,
	RunRequest: &checks.Request{
		IngressName: "host-rules",
		Path:        "",
		Hostname:    "aaa.bbb.foo.com",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "default-backend", "expected the downstream service would be '%s' but was '%s'")
			a.E.DeepEquals(req.Host, "aaa.bbb.foo.com", "expected the request host would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var hostRulesWildcardNoLabelCheck = &checks.Check{
	Name:        "host-rules-wildcard-no-label",
	Description: "Ingress with wildcard host rule should match exactly one single label & fallback to default-backend (*.foo.com does not match foo.com)",
	APIVersions: apiversion.NetworkingV1Beta1,
	RunRequest: &checks.Request{
		IngressName: "host-rules",
		Path:        "",
		Hostname:    "foo.com",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "default-backend", "expected the downstream service would be '%s' but was '%s'")
			a.E.DeepEquals(req.Host, "foo.com", "expected the request host would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}
