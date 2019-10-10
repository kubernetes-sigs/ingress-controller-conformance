package checks

import (
	"fmt"
	"github.com/datawire/ingress-conformance/internal/pkg/k8s"
)

func init() {
	Checks.AddCheck(singleServiceCheck)
}

var singleServiceCheck = &Check{
	Name:        "single-service",
	Description: "Ingress with no rules should send traffic to the correct backend service",
	Run: func(check *Check, config Config) (success bool, err error) {
		host, err := k8s.GetIngressHost("default", "single-service")
		if err != nil {
			return
		}

		resp, err := captureRequest(fmt.Sprintf("http://%s", host), "")
		if err != nil {
			return
		}

		a := new(assertionSet)
		a.equals(assert{resp.StatusCode, 200, "Expected StatusCode to be %s but was %s"})
		a.equals(assert{resp.TestId, "single-service", "Expected the responding service would be '%s' but was '%s'"})

		if a.Error() == "" {
			success = true
		} else {
			fmt.Print(a)
		}
		return
	},
}
