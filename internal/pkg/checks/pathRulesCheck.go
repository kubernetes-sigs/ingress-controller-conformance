package checks

import (
	"fmt"
	"github.com/datawire/ingress-conformance/internal/pkg/k8s"
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
	Name: "path-rules-foo",
	Run: func(check *Check, config Config) (success bool, err error) {
		resp, err := captureRequest(fmt.Sprintf("http://%s/foo", host), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		a.equals(assert{resp.StatusCode, 200, "Expected StatusCode to be %s but was %s"})
		a.equals(assert{resp.TestId, "path-rules-foo", "Expected the responding service would be '%s' but was '%s'"})
		a.equals(assert{resp.Path, "/foo", "Expected the request path would be '%s' but was '%s'"})

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}

var pathRulesFooTrailingSlashCheck = &Check{
	Name: "path-rules-foo-trailing",
	Run: func(check *Check, config Config) (success bool, err error) {
		resp, err := captureRequest(fmt.Sprintf("http://%s/foo/", host), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		a.equals(assert{resp.StatusCode, 200, "Expected StatusCode to be %s but was %s"})
		a.equals(assert{resp.TestId, "path-rules-foo", "Expected the responding service would be '%s' but was '%s'"})
		a.equals(assert{resp.Path, "/foo/", "Expected the request path would be '%s' but was '%s'"})

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}

var pathRulesBarCheck = &Check{
	Name: "path-rules-bar",
	Run: func(check *Check, config Config) (success bool, err error) {
		resp, err := captureRequest(fmt.Sprintf("http://%s/bar/", host), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		a.equals(assert{resp.StatusCode, 200, "Expected StatusCode to be %s but was %s"})
		a.equals(assert{resp.TestId, "path-rules-bar", "Expected the responding service would be '%s' but was '%s'"})
		a.equals(assert{resp.Path, "/bar/", "Expected the request path would be '%s' but was '%s'"})

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}

var pathRulesBarSubpathCheck = &Check{
	Name: "path-rules-bar-subpath",
	Run: func(check *Check, config Config) (success bool, err error) {
		resp, err := captureRequest(fmt.Sprintf("http://%s/bar//bershop", host), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		a.equals(assert{resp.StatusCode, 200, "Expected StatusCode to be %s but was %s"})
		a.equals(assert{resp.TestId, "path-rules-bar", "Expected the responding service would be '%s' but was '%s'"})
		a.equals(assert{resp.Path, "/bar//bershop", "Expected the request path would be '%s' but was '%s'"})

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}

// TODO: Implement more checks on edge cases like leading `/`, query params, and encoding
