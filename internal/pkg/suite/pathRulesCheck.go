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
	/*
			TODO: There are currently no 'exact' path types validations since it is unsupported in v1beta1
		          For now, we validate only the 'prefix' path types, which is closer to the v1beta1 'prefix'
		          path type assumption.
	*/
	pathRulesCheck.AddCheck(pathRulesExactCheck)

	pathRulesPrefixCheck.AddCheck(pathRulesPrefixAllPathsCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixFooCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixFooSlashCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixFoCheck)
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

var pathRulesHost string
var pathRulesCheck = &checks.Check{
	Name: "path-rules",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		if config.UseInsecureHost != "" {
			pathRulesHost = config.UseInsecureHost
			return true, nil
		}
		var err error
		pathRulesHost, err = k8s.GetIngressHost("default", "path-rules")
		if err != nil {
			return false, err
		}
		return true, nil
	},
}

// placeholder check for dividing the pathRulesCheck into a distinct hierarchy for Exact path tests
var pathRulesExactCheck = &checks.Check{
	Name: "path-rules-exact",
}

// placeholder check for dividing the pathRulesCheck into a distinct hierarchy for Prefix path tests
var pathRulesPrefixCheck = &checks.Check{
	Name: "path-rules-prefix",
}

var pathRulesPrefixAllPathsCheck = &checks.Check{
	Name:        "path-rules-prefix-all-paths",
	Description: "Ingress with prefix path rule '/' should match all paths",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixFooCheck = &checks.Check{
	Name:        "path-rules-prefix-foo",
	Description: "Ingress with prefix path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request path (/foo matches /foo)",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s/foo", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-foo", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/foo", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixFooSlashCheck = &checks.Check{
	Name:        "path-rules-prefix-foo-slash",
	Description: "Ingress with prefix path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request path (/foo matches /foo/)",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s/foo/", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-foo", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/foo/", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixFoCheck = &checks.Check{
	Name:        "path-rules-prefix-fo",
	Description: "Ingress with prefix path rule without a trailing slash should not match partial paths (/foo does not match /fo)",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s/fo", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/fo", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixAaaBbbCheck = &checks.Check{
	Name:        "path-rules-prefix-aaa-bbb",
	Description: "Ingress with prefix path rule with a trailing slash should send traffic to the correct backend service, and preserve the original request path (/aaa/bbb/ matches /aaa/bbb)",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s/aaa/bbb", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-aaa-bbb", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/aaa/bbb", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixAaaBbbSlashCheck = &checks.Check{
	Name:        "path-rules-prefix-aaa-bbb-slash",
	Description: "Ingress with prefix path rule with a trailing slash should send traffic to the correct backend service, and preserve the original request path (/aaa/bbb/ matches /aaa/bbb/)",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s/aaa/bbb/", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-aaa-bbb", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/aaa/bbb/", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixAaaBbbCccCheck = &checks.Check{
	Name:        "path-rules-prefix-aaa-bbb-ccc",
	Description: "Ingress with prefix path rule with a trailing slash should match subpath, send traffic to the correct backend service, and preserve the original request path (/aaa/bbb/ matches /aaa/bbb/ccc)",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s/aaa/bbb/ccc", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-aaa-bbb", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/aaa/bbb/ccc", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixAaaBbbcccCheck = &checks.Check{
	Name:        "path-rules-prefix-aaa-bbbccc",
	Description: "Ingress with prefix path rule with a trailing slash should not match string prefix (/aaa/bbb/ does not match /aaa/bbbccc)",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s/aaa/bbbccc", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/aaa/bbbccc", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixConsecutiveSlashesCheck = &checks.Check{
	Name:        "path-rules-prefix-consecutive-slashes",
	Description: "Ingress with prefix path rule with consecutive slashes are ignored",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s/routes/with/consecutive//slashes///are-ignored", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/routes/with/consecutive//slashes///are-ignored", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixConsecutiveSlashesNormalizedCheck = &checks.Check{
	Name:        "path-rules-prefix-consecutive-slashes-normalized",
	Description: "Ingress with prefix path rule with consecutive slashes are ignored with normalized request",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s/routes/with/consecutive/slashes/are-ignored", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/routes/with/consecutive/slashes/are-ignored", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixInvalidCharactersCheck = &checks.Check{
	Name:        "path-rules-prefix-invalid-characters",
	Description: "Ingress with prefix path rule with invalid characters are ignored",
	Run: func(check *checks.Check, config checks.Config) (bool, error) {
		req, res, err := checks.CaptureRoundTrip(fmt.Sprintf("http://%s/routes with invalid characters are ignored!", pathRulesHost), "path-rules")
		if err != nil {
			return false, err
		}

		a := &checks.AssertionSet{}
		// Assert the request received from the downstream service
		a.DeepEquals(req.DownstreamServiceId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.DeepEquals(req.Path, "/routes%20with%20invalid%20characters%20are%20ignored%21", "expected the request path would be '%s' but was '%s'")
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
