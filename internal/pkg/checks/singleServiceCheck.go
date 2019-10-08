package checks

import (
	"fmt"
	"github.com/datawire/ingress-conformance/internal/pkg/k8s"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
)

func init() {
	Checks.AddCheck(singleServiceCheck)
}

var singleServiceCheck = &Check{
	Name: "single-service",
	Run: func(check *Check, config Config) (success bool, err error) {
		ingressInterface, err := k8s.Client.NetworkingV1beta1().Ingresses("default").Get("single-service", v1.GetOptions{})
		if err != nil {
			fmt.Printf(err.Error())
		}
		host := ingressInterface.Status.LoadBalancer.Ingress[0].Hostname

		resp, err := http.Get(fmt.Sprintf("http://%s", host))
		if err != nil {
			return
		}
		if resp.StatusCode == 200 {
			success = true
		}
		return
	},
}
