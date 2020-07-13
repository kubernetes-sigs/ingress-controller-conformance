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

package pathrules

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"

	"sigs.k8s.io/ingress-controller-conformance/test/kubernetes"
	tstate "sigs.k8s.io/ingress-controller-conformance/test/state"
)

var (
	state *tstate.Scenario
)

// IMPORTANT: Steps definitions are generated and should not be modified
// by hand but rather through make codegen. DO NOT EDIT.
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^an Ingress resource in a new random namespace$`, anIngressResourceInANewRandomNamespace)
	ctx.Step(`^The Ingress status shows the IP address or FQDN where it is exposed$`, theIngressStatusShowsTheIPAddressOrFQDNWhereItIsExposed)
	ctx.Step(`^I send a "([^"]*)" request to "([^"]*)"$`, iSendARequestTo)
	ctx.Step(`^the response status-code must be (\d+)$`, theResponseStatuscodeMustBe)
	ctx.Step(`^the response must be served by the "([^"]*)" service$`, theResponseMustBeServedByTheService)
	ctx.Step(`^the request path must be "([^"]*)"$`, theRequestPathMustBe)

	ctx.BeforeScenario(func(*godog.Scenario) {
		state = tstate.New(nil)
	})

	ctx.AfterScenario(func(*messages.Pickle, error) {
		// delete namespace an all the content
		_ = kubernetes.DeleteNamespace(kubernetes.KubeClient, state.Namespace)
	})
}

func anIngressResourceInANewRandomNamespace(arg1 *messages.PickleStepArgument_PickleDocString) error {
	return godog.ErrPending
}

func theIngressStatusShowsTheIPAddressOrFQDNWhereItIsExposed() error {
	return godog.ErrPending
}

func iSendARequestTo(arg1 string, arg2 string) error {
	return godog.ErrPending
}

func theResponseStatuscodeMustBe(arg1 int) error {
	return godog.ErrPending
}

func theResponseMustBeServedByTheService(arg1 string) error {
	return godog.ErrPending
}

func theRequestPathMustBe(arg1 string) error {
	return godog.ErrPending
}
