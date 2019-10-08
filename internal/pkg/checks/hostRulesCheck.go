package checks

import (
	"fmt"
	"github.com/datawire/ingress-conformance/internal/pkg/k8s"
)

func init() {
	Checks.AddCheck(hostRulesCheck)
}

var hostRulesCheck = &Check{
	Name: "host-rules",
	Run: func(check *Check, config Config) (success bool, err error) {
		host, err := k8s.GetIngressHost("default", "host-rules")
		if err != nil {
			return
		}

		resp, err := captureRequest(fmt.Sprintf("http://%s", host), "foo.bar.com")
		if err != nil {
			return
		}

		a := new(assertionSet)
		a.equals(assert{resp.StatusCode, 200, "Expected StatusCode to be %s but was %s"})
		a.equals(assert{resp.TestId, "host-rules", "Expected the responding service would be '%s' but was '%s'"})
		a.equals(assert{resp.Host, "foo.bar.com", "Expected the responding host would be '%s' but was '%s'"})

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}
