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
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/checks"
)

func init() {
	checks.AllChecks.AddCheck(singleServiceCheck)
}

var singleServiceCheck = &checks.Check{
	Name:        "default-backend",
	Description: "Ingress with a single default backend should send traffic to the correct backend service",
	RunRequest: &checks.Request{
		IngressName: "default-backend",
		Path:        "",
		Host:        "",
		Insecure:    true,
		DoCheck: func(req *checks.CapturedRequest, res *checks.CapturedResponse) (*checks.AssertionSet, error) {
			a := &checks.AssertionSet{}

			a.DeepEquals(req.DownstreamServiceId, "default-backend", "expected the downstream service would be '%s' but was '%s'")
			a.DeepEquals(req.Method, "GET", "expected the originating request method would be '%s' but was '%s'")
			a.DeepEquals(req.Proto, "HTTP/1.1", "expected the originating request protocol would be '%s' but was '%s'")
			a.ContainsHeaders(req.Headers, []string{"User-Agent"})
			// Assert the downstream service response
			a.DeepEquals(res.StatusCode, 200, "expected statuscode to be %s but was %s")
			a.DeepEquals(res.Proto, "HTTP/1.1", "expected the response protocol would be %s but was %s")
			a.ContainsExactHeaders(res.Headers, []string{"Content-Length", "Content-Type", "Date", "Server"})

			return a, nil
		},
	},
}
