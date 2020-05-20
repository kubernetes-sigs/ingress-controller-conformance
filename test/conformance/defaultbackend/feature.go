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

package defaultbackend

import (
	"github.com/cucumber/godog"
	"github.com/cucumber/messages-go/v10"

	"github.com/kubernetes-sigs/ingress-controller-conformance/test/kubernetes"
	tstate "github.com/kubernetes-sigs/ingress-controller-conformance/test/state"
)

var (
	state *tstate.Scenario
)

// IMPORTANT: Steps definitions are generated and should not be modified
// by hand but rather through make codegen. DO NOT EDIT.
func FeatureContext(s *godog.Suite) { 
	s.Step(`^a new random namespace$`, aNewRandomNamespace)
	s.Step(`^reading Ingress from manifest "([^"]*)"$`, readingIngressFromManifest)
	s.Step(`^creating Ingress from manifest returns an error message containing "([^"]*)"$`, creatingIngressFromManifestReturnsAnErrorMessageContaining)
	s.Step(`^creating Ingress from manifest$`, creatingIngressFromManifest)
	s.Step(`^The ingress status shows the IP address or FQDN where is exposed$`, theIngressStatusShowsTheIPAddressOrFQDNWhereIsExposed)
	s.Step(`^Header "([^"]*)" with value "([^"]*)"$`, headerWithValue)
	s.Step(`^Send HTTP request with method "([^"]*)"$`, sendHTTPRequestWithMethod)
	s.Step(`^Response status code is (\d+)$`, responseStatusCodeIs)
	s.Step(`^Send HTTP request with <path> and <method> checking response status code is (\d+):$`, sendHTTPRequestWithPathAndMethodCheckingResponseStatusCodeIs)
	s.Step(`^creating objects from directory "([^"]*)"$`, creatingObjectsFromDirectory)
	s.Step(`^With path "([^"]*)"$`, withPath)

	s.BeforeScenario(func(this *messages.Pickle) {
		state = tstate.New(nil)
	})

	s.AfterScenario(func(*messages.Pickle, error) {
		// delete namespace an all the content
		_ = kubernetes.DeleteNamespace(kubernetes.KubeClient, state.Namespace)
	})
}


func aNewRandomNamespace() error {
	return godog.ErrPending
}

func readingIngressFromManifest(arg1 string) error {
	return godog.ErrPending
}

func creatingIngressFromManifestReturnsAnErrorMessageContaining(arg1 string) error {
	return godog.ErrPending
}

func creatingIngressFromManifest() error {
	return godog.ErrPending
}

func theIngressStatusShowsTheIPAddressOrFQDNWhereIsExposed() error {
	return godog.ErrPending
}

func headerWithValue(arg1 string, arg2 string) error {
	return godog.ErrPending
}

func sendHTTPRequestWithMethod(arg1 string) error {
	return godog.ErrPending
}

func responseStatusCodeIs(arg1 int) error {
	return godog.ErrPending
}

func sendHTTPRequestWithPathAndMethodCheckingResponseStatusCodeIs(arg1 int, arg2 *messages.PickleStepArgument_PickleTable) error {
	return godog.ErrPending
}

func creatingObjectsFromDirectory(arg1 string) error {
	return godog.ErrPending
}

func withPath(arg1 string) error {
	return godog.ErrPending
}

