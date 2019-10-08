package checks

import (
	"fmt"
	"github.com/datawire/ingress-conformance/internal/pkg/k8s"
)

func init() {
	Checks.AddCheck(pathRulesCheck)
}

var pathRulesCheck = &Check{
	Name: "path-rules",
	Run: func(check *Check, config Config) (success bool, err error) {
		host, err := k8s.GetIngressHost("default", "path-rules")
		if err != nil {
			return
		}

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
