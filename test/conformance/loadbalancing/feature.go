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

package loadbalancing

import (
	"fmt"
	"net/url"

	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"
	"k8s.io/apimachinery/pkg/util/sets"

	"sigs.k8s.io/ingress-controller-conformance/test/http"
	"sigs.k8s.io/ingress-controller-conformance/test/kubernetes"
	tstate "sigs.k8s.io/ingress-controller-conformance/test/state"
)

var (
	state *tstate.Scenario

	resultStatus map[int]sets.String
)

// IMPORTANT: Steps definitions are generated and should not be modified
// by hand but rather through make codegen. DO NOT EDIT.

// InitializeScenario configures the Feature to test
func InitializeScenario(ctx *godog.ScenarioContext) {
	ctx.Step(`^a new random namespace$`, aNewRandomNamespace)
	ctx.Step(`^an Ingress resource named "([^"]*)" with this spec:$`, anIngressResourceNamedWithThisSpec)
	ctx.Step(`^The Ingress status shows the IP address or FQDN where it is exposed$`, theIngressStatusShowsTheIPAddressOrFQDNWhereItIsExposed)
	ctx.Step(`^The backend deployment "([^"]*)" for the ingress resource is scaled to (\d+)$`, theBackendDeploymentForTheIngressResourceIsScaledTo)
	ctx.Step(`^I send (\d+) requests to "([^"]*)"$`, iSendRequestsTo)
	ctx.Step(`^all the responses status-code must be (\d+) and the response body should contain the IP address of (\d+) different Kubernetes pods$`, allTheResponsesStatuscodeMustBeAndTheResponseBodyShouldContainTheIPAddressOfDifferentKubernetesPods)

	ctx.BeforeScenario(func(*godog.Scenario) {
		state = tstate.New()
		resultStatus = make(map[int]sets.String, 0)
	})

	ctx.AfterScenario(func(*messages.Pickle, error) {
		// delete namespace an all the content
		_ = kubernetes.DeleteNamespace(kubernetes.KubeClient, state.Namespace)
	})
}

func aNewRandomNamespace() error {
	ns, err := kubernetes.NewNamespace(kubernetes.KubeClient)
	if err != nil {
		return err
	}

	state.Namespace = ns
	return nil
}

func anIngressResourceNamedWithThisSpec(name string, spec *messages.PickleStepArgument_PickleDocString) error {
	ingress, err := kubernetes.IngressFromSpec(name, state.Namespace, spec.GetContent())
	if err != nil {
		return err
	}

	err = kubernetes.DeploymentsFromIngress(kubernetes.KubeClient, ingress)
	if err != nil {
		return err
	}

	err = kubernetes.NewIngress(kubernetes.KubeClient, state.Namespace, ingress)
	if err != nil {
		return err
	}

	state.IngressName = name

	return nil
}

func theIngressStatusShowsTheIPAddressOrFQDNWhereItIsExposed() error {
	ingress, err := kubernetes.WaitForIngressAddress(kubernetes.KubeClient, state.Namespace, state.IngressName)
	if err != nil {
		return err
	}

	state.IPOrFQDN = ingress
	return err
}

func iSendRequestsTo(totalRequest int, rawURL string) error {
	u, err := url.Parse(rawURL)
	if err != nil {
		return err
	}

	for iteration := 1; iteration <= totalRequest; iteration++ {
		capturedRequest, capturedResponse, err := http.CaptureRoundTrip("GET", u.Scheme, u.Host, u.Path, state.IPOrFQDN)
		if err != nil {
			return err
		}

		if resultStatus[capturedResponse.StatusCode] == nil {
			resultStatus[capturedResponse.StatusCode] = sets.NewString()
		}

		resultStatus[capturedResponse.StatusCode].Insert(capturedRequest.Pod)
	}

	return nil
}

func allTheResponsesStatuscodeMustBeAndTheResponseBodyShouldContainTheIPAddressOfDifferentKubernetesPods(statusCode int, pods int) error {
	results, ok := resultStatus[statusCode]
	if !ok {
		return fmt.Errorf("no reponses for status code %v returned", statusCode)
	}

	if results.Len() != pods {
		return fmt.Errorf("expected %v different POD IP addresses/FQDN for status code %v but %v was returned", pods, statusCode, results.Len())
	}

	return nil
}

func theBackendDeploymentForTheIngressResourceIsScaledTo(deployment string, replicas int) error {
	return kubernetes.ScaleIngressBackendDeployment(kubernetes.KubeClient, state.Namespace, state.IngressName, deployment, replicas)
}
