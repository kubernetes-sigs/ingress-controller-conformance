/*
Copyright 2020 The Kubernetes Authors.
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
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
)

func FeatureContext(s *godog.Suite) {
	s.Step(`^I have an Ingress named "([^"]*)" in the "([^"]*)" namespace$`, iHaveAnIngressNamedInTheNamespace)
	s.Step(`^I send a "([^"]*)" "([^"]*)" request$`, iSendARequest)
	s.Step(`^the response status-code must be (\d+)$`, theResponseStatuscodeMustBe)
	s.Step(`^the response must be served by the "([^"]*)" service$`, theResponseMustBeServedByTheService)
	s.Step(`^the response proto must be "([^"]*)"$`, theResponseProtoMustBe)
	s.Step(`^the response headers must contain <key> with matching <value>$`, theResponseHeadersMustContainHeaderWithMatchingValue)
	s.Step(`^the request method must be "([^"]*)"$`, theRequestMethodMustBe)
	s.Step(`^the request proto must be "([^"]*)"$`, theRequestProtoMustBe)
	s.Step(`^the request host must be "([^"]*)"$`, theRequestHostMustBe)
	s.Step(`^the request path must be "([^"]*)"$`, theRequestPathMustBe)
	s.Step(`^the request headers must contain <key> with matching <value>$`, theRequestHeadersMustContainHeaderWithMatchingValue)

	// clean the state before every scenario
	s.BeforeScenario(func(p *messages.Pickle) {
		ingressEndpoint = ""
		captureRoundTrip = nil
	})
}
