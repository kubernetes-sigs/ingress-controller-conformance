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

package state

import (
	"io/ioutil"
	"net/http"

	v1beta1 "k8s.io/api/networking/v1beta1"
)

// Scenario holds state for a test scenario
type Scenario struct {
	client *http.Client

	RequestPath string

	RequestHeaders http.Header

	ResponseBody    []byte
	ResponseHeaders http.Header

	StatusCode int

	Namespace string

	IngressManifest string

	Ingress *v1beta1.Ingress
	Address string
}

// New creates a new state to use in a test Scenario
func New(client *http.Client) *Scenario {
	if client == nil {
		client = &http.Client{}
	}

	return &Scenario{
		client:         client,
		RequestPath:    "/",
		RequestHeaders: make(http.Header),
	}
}

// SendRequest sends an HTTP request and updates the
// state. In case of an error, the HTTP state is
// removed and returns an error.
func (f *Scenario) SendRequest(req *http.Request) error {
	req.Header = f.RequestHeaders

	resp, err := f.client.Do(req)
	if err != nil {
		f.ResponseBody = nil
		f.StatusCode = 0
		f.ResponseHeaders = nil

		return err
	}

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	f.ResponseBody = bodyBytes
	f.ResponseHeaders = resp.Header.Clone()
	f.StatusCode = resp.StatusCode

	defer resp.Body.Close()

	return nil
}

// AddRequestHeader Add adds the key, value pair to the header.
// It appends to any existing values associated with key.
func (f *Scenario) AddRequestHeader(header, value string) {
	f.RequestHeaders.Add(header, value)
}
