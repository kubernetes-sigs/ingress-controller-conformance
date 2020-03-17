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

type Assertions struct {
	E AssertionSet
	W AssertionSet
}
type AssertionSet []error

// DeepEquals asserts actual and expected object parameters are deeply equal.
//
// The errorTemplate format specifier must include a first %v for the expected value
// and a second %v for the actual value.
func (a *AssertionSet) DeepEquals(actual interface{}, expected interface{}, errorTemplate string) {
	if errorTemplate == "" {
		errorTemplate = "expected '%v' but was '%v'"
	}
	if !reflect.DeepEqual(expected, actual) {
		err := fmt.Errorf(errorTemplate, expected, actual)
		*a = append(*a, err)
	}
}

// ContainsHeaders asserts the actual headers contains the expected header keys.
func (a *AssertionSet) ContainsHeaders(actual map[string][]string, expected []string) {
	errorTemplate := "expected headers to contain '%v' but contained '%v'"
	for _, expectedKey := range expected {
		if actual[expectedKey] == nil {
			err := fmt.Errorf(errorTemplate, expectedKey, actual)
			*a = append(*a, err)
		}
	}
}

// ContainsExactHeaders asserts the actual headers contains exactly the expected
// header keys and nothing more.
func (a *AssertionSet) ContainsExactHeaders(actual map[string][]string, expected []string) {
	a.ContainsHeaders(actual, expected)
	errorTemplate := "expected headers to only contain '%v' but contained '%v'"
	if len(actual) != len(expected) {
		err := fmt.Errorf(errorTemplate, expected, actual)
		*a = append(*a, err)
	}
}

func (a *Assertions) String() string {
	var err string
	for i, e := range a.E {
		err += fmt.Sprintf("\tERROR %d) Assertion failed: %s\n", i+1, e.Error())
	}
	for i, w := range a.W {
		err += fmt.Sprintf("\tWARN  %d) Assertion failed: %s\n", i+1, w.Error())
	}
	return err
}

func (a *Assertions) Passed() bool {
	if len(a.E) > 0 {
		return false
	}
	return true
}
