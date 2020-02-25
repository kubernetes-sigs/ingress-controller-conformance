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

// Ingress conformance test harness machinery
package checks

import (
	"fmt"
	"reflect"
)

// AssertionSet performs checks and accumulates assertion errors
type AssertionSet []error

// Assert actual and expected parameters are deeply equal
func (a *AssertionSet) Equals(actual interface{}, expected interface{}, errorTemplate string) {
	if errorTemplate == "" {
		errorTemplate = "Expected '%s' but was '%s'"
	}
	if !reflect.DeepEqual(expected, actual) {
		err := fmt.Errorf(errorTemplate, expected, actual)
		*a = append(*a, err)
	}
}

// Assert the actual headers contains the expected headers key
func (a *AssertionSet) ContainsHeaders(actual map[string][]string, expected []string, errorTemplate string) {
	if errorTemplate == "" {
		errorTemplate = "Expected to contain '%s' but contained '%s'"
	}
	for _, expectedKey := range expected {
		if actual[expectedKey] == nil {
			err := fmt.Errorf(errorTemplate, expectedKey, actual)
			*a = append(*a, err)
		}
	}
}

// Assert the actual headers contains exactly the expected headers key and no more
func (a *AssertionSet) ContainsExactHeaders(actual map[string][]string, expected []string, errorTemplate string) {
	a.ContainsHeaders(actual, expected, errorTemplate)
	if errorTemplate == "" {
		errorTemplate = "Expected to only contain '%s' but contained '%s'"
	}
	if len(actual) != len(expected) {
		err := fmt.Errorf(errorTemplate, expected, actual)
		*a = append(*a, err)
	}
}

func (a *AssertionSet) Error() (err string) {
	for i, e := range *a {
		err += fmt.Sprintf("\t%d) Assertion failed: %s\n", i+1, e.Error())
	}
	return
}
