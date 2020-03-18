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
	pathRulesExactCheck.AddCheck(pathRulesExactFooCheck)
	pathRulesExactCheck.AddCheck(pathRulesExactBarCheck)
	pathRulesExactCheck.AddCheck(pathRulesExactBarNoSlashCheck)
	pathRulesExactCheck.AddCheck(pathRulesExactBarAaaCheck)
	pathRulesCheck.AddCheck(pathRulesExactCheck)

	pathRulesPrefixCheck.AddCheck(pathRulesPrefixAllPathsCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixFooCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixFooSlashCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixFoCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixFooxyzCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixAaaBbbCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixAaaBbbSlashCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixAaaBbbCccCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixAaaBbbcccCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixConsecutiveSlashesCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixConsecutiveSlashesNormalizedCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixInvalidCharactersCheck)
	pathRulesCheck.AddCheck(pathRulesPrefixCheck)

	checks.AllChecks.AddCheck(pathRulesCheck)
}

// placeholder check for the path-rules checks hierarchy
var pathRulesCheck = &checks.Check{
	Name: "path-rules",
}

// placeholder check for dividing the pathRulesCheck into a distinct hierarchy for Exact path tests
var pathRulesExactCheck = &checks.Check{
	Name: "path-rules-exact",
}

var pathRulesExactFooCheck = &checks.Check{
	Name:        "path-rules-exact-foo",
	Description: "Ingress with exact path rule is preferred to prefix match and should send traffic to the correct backend service (/foo matches /foo)",
	APIVersions: apiversion.NetworkingV1Beta1,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/foo",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-exact", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/foo", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesExactBarCheck = &checks.Check{
	Name:        "path-rules-exact-bar",
	Description: "Ingress with exact path rule should send traffic to the correct backend service (/bar/ matches /bar/)",
	APIVersions: apiversion.NetworkingV1Beta1,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/bar/",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-exact", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/bar/", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesExactBarNoSlashCheck = &checks.Check{
	Name:        "path-rules-exact-bar-no-slash",
	Description: "Ingress with exact path rule should not match partial paths (/bar/ does not match /bar)",
	APIVersions: apiversion.NetworkingV1Beta1,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/bar",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/bar", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesExactBarAaaCheck = &checks.Check{
	Name:        "path-rules-exact-bar-aaa",
	Description: "Ingress with exact path rule should not match partial paths (/bar/ does not match /bar/aaa)",
	APIVersions: apiversion.NetworkingV1Beta1,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/bar/aaa",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/bar/aaa", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

// placeholder check for dividing the pathRulesCheck into a distinct hierarchy for Prefix path tests
var pathRulesPrefixCheck = &checks.Check{
	Name: "path-rules-prefix",
}

var pathRulesPrefixAllPathsCheck = &checks.Check{
	Name:        "path-rules-prefix-all-paths",
	Description: "Ingress with prefix path rule '/' should match all paths",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
			a.E.DeepEquals(req.Path, "/", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixFooCheck = &checks.Check{
	Name:        "path-rules-prefix-foo",
	Description: "Ingress with prefix path rule without a trailing slash should send traffic to the correct backend service (/foo matches /foo)",
	APIVersions: apiversion.ExtensionsV1Beta1,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/foo",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-foo", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/foo", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixFooSlashCheck = &checks.Check{
	Name:        "path-rules-prefix-foo-slash",
	Description: "Ingress with prefix path rule without a trailing slash should send traffic to the correct backend service (/foo matches /foo/)",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/foo/",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-foo", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/foo/", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixFoCheck = &checks.Check{
	Name:        "path-rules-prefix-fo",
	Description: "Ingress with prefix path rule should not match partial paths (/foo does not match /fo)",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/fo",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/fo", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixFooxyzCheck = &checks.Check{
	Name:        "path-rules-prefix-fooxyz",
	Description: "Ingress with prefix path rule should not match string prefixes (/foo does not match /fooxyz)",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/fooxyz",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/fooxyz", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixAaaBbbCheck = &checks.Check{
	Name:        "path-rules-prefix-aaa-bbb",
	Description: "Ingress with prefix path rule with a trailing slash should send traffic to the correct backend service (/aaa/bbb/ matches /aaa/bbb)",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/aaa/bbb",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-aaa-bbb", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/aaa/bbb", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixAaaBbbSlashCheck = &checks.Check{
	Name:        "path-rules-prefix-aaa-bbb-slash",
	Description: "Ingress with prefix path rule with a trailing slash should send traffic to the correct backend service (/aaa/bbb/ matches /aaa/bbb/)",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/aaa/bbb/",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-aaa-bbb", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/aaa/bbb/", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixAaaBbbCccCheck = &checks.Check{
	Name:        "path-rules-prefix-aaa-bbb-ccc",
	Description: "Ingress with prefix path rule with a trailing slash should match subpath, send traffic to the correct backend service (/aaa/bbb/ matches /aaa/bbb/ccc)",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/aaa/bbb/ccc",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-aaa-bbb", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/aaa/bbb/ccc", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixAaaBbbcccCheck = &checks.Check{
	Name:        "path-rules-prefix-aaa-bbbccc",
	Description: "Ingress with prefix path rule with a trailing slash should not match string prefix (/aaa/bbb/ does not match /aaa/bbbccc)",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/aaa/bbbccc",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/aaa/bbbccc", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixConsecutiveSlashesCheck = &checks.Check{
	Name:        "path-rules-prefix-consecutive-slashes",
	Description: "Ingress with prefix path rule with consecutive slashes are ignored",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/routes/with/consecutive//slashes///are-ignored",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/routes/with/consecutive//slashes///are-ignored", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixConsecutiveSlashesNormalizedCheck = &checks.Check{
	Name:        "path-rules-prefix-consecutive-slashes-normalized",
	Description: "Ingress with prefix path rule with consecutive slashes are ignored with normalized request",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/routes/with/consecutive/slashes/are-ignored",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/routes/with/consecutive/slashes/are-ignored", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}

var pathRulesPrefixInvalidCharactersCheck = &checks.Check{
	Name:        "path-rules-prefix-invalid-characters",
	Description: "Ingress with prefix path rule with invalid characters are ignored",
	APIVersions: apiversion.All,
	RunRequest: &checks.Request{
		IngressName: "path-rules",
		Path:        "/routes with invalid characters are ignored!",
		Hostname:    "path-rules",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.Assertions, error) {
			a := &checks.Assertions{}
			// Assert the request received from the downstream service
			a.E.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
			a.W.DeepEquals(req.Path, "/routes%20with%20invalid%20characters%20are%20ignored%21", "expected the request path would be '%s' but was '%s'")
			// Assert the downstream service response
			a.E.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")

			return a, nil
		},
	},
}
