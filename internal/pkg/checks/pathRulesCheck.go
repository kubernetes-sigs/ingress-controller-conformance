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
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixAaaBbCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixAaaBbbcccCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixConsecutiveSlashesCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixConsecutiveSlashesNormalizedCheck)
	pathRulesPrefixCheck.AddCheck(pathRulesPrefixInvalidCharactersCheck)
	pathRulesCheck.AddCheck(pathRulesPrefixCheck)

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

var pathRulesExactCheck = &Check{
	Name: "path-rules-exact",
}

var pathRulesPrefixCheck = &Check{
	Name: "path-rules-prefix",
}

var pathRulesPrefixAllPathsCheck = &Check{
	Name:        "path-rules-prefix-all-paths",
	Description: "Ingress with prefix path rule '/' should match all paths",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s", pathRulesHost), "path-rules")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixFooCheck = &Check{
	Name:        "path-rules-prefix-foo",
	Description: "Ingress with prefix path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request path (/foo matches /foo)",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/foo", pathRulesHost), "path-rules")
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

var pathRulesPrefixFooSlashCheck = &Check{
	Name:        "path-rules-prefix-foo-slash",
	Description: "Ingress with prefix path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request path (/foo matches /foo/)",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/foo/", pathRulesHost), "path-rules")
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

var pathRulesPrefixFoCheck = &Check{
	Name:        "path-rules-prefix-fo",
	Description: "Ingress with prefix path rule without a trailing slash should not match partial paths (/foo does not match /fo)",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/fo", pathRulesHost), "path-rules")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/fo", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixAaaBbbCheck = &Check{
	Name:        "path-rules-prefix-aaa-bbb",
	Description: "Ingress with prefix path rule with a trailing slash should send traffic to the correct backend service, and preserve the original request path (/aaa/bbb/ matches /aaa/bbb)",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/aaa/bbb", pathRulesHost), "path-rules")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-aaa-bbb", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/aaa/bbb", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixAaaBbbSlashCheck = &Check{
	Name:        "path-rules-prefix-aaa-bbb-slash",
	Description: "Ingress with prefix path rule with a trailing slash should send traffic to the correct backend service, and preserve the original request path (/aaa/bbb/ matches /aaa/bbb/)",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/aaa/bbb/", pathRulesHost), "path-rules")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-aaa-bbb", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/aaa/bbb/", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixAaaBbbCccCheck = &Check{
	Name:        "path-rules-prefix-aaa-bbb-ccc",
	Description: "Ingress with prefix path rule with a trailing slash should match subpath, send traffic to the correct backend service, and preserve the original request path (/aaa/bbb/ matches /aaa/bbb/ccc)",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/aaa/bbb/ccc", pathRulesHost), "path-rules")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-aaa-bbb", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/aaa/bbb/ccc", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixAaaBbCheck = &Check{
	Name:        "path-rules-prefix-aaa-bb",
	Description: "Ingress with prefix path rule with a trailing slash should not match partial paths (/aaa/bbb/ does not match /aaa/bb)",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/aaa/bb", pathRulesHost), "path-rules")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/aaa/bb", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixAaaBbbcccCheck = &Check{
	Name:        "path-rules-prefix-aaa-bbbccc",
	Description: "Ingress with prefix path rule with a trailing slash should not match string prefix (/aaa/bbb/ does not match /aaa/bbbccc)",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/aaa/bbbccc", pathRulesHost), "path-rules")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/aaa/bbbccc", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixConsecutiveSlashesCheck = &Check{
	Name:        "path-rules-prefix-consecutive-slashes",
	Description: "Ingress with prefix path rule with consecutive slashes are ignored",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/routes/with/consecutive//slashes///are-ignored", pathRulesHost), "path-rules")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/routes/with/consecutive//slashes///are-ignored", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixConsecutiveSlashesNormalizedCheck = &Check{
	Name:        "path-rules-prefix-consecutive-slashes-normalized",
	Description: "Ingress with prefix path rule with consecutive slashes are ignored with normalized request",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/routes/with/consecutive/slashes/are-ignored", pathRulesHost), "path-rules")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/routes/with/consecutive/slashes/are-ignored", "expected the request path would be '%s' but was '%s'")
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

var pathRulesPrefixInvalidCharactersCheck = &Check{
	Name:        "path-rules-prefix-invalid-characters",
	Description: "Ingress with prefix path rule with invalid characters are ignored",
	Run: func(check *Check, config Config) (success bool, err error) {
		req, res, err := captureRequest(fmt.Sprintf("http://%s/routes with invalid characters are ignored!", pathRulesHost), "path-rules")
		if err != nil {
			return
		}

		a := new(assertionSet)
		// Assert the request received from the downstream service
		a.equals(req.TestId, "path-rules-catchall", "expected the downstream service would be '%s' but was '%s'")
		a.equals(req.Path, "/routes%20with%20invalid%20characters%20are%20ignored%21", "expected the request path would be '%s' but was '%s'")
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
