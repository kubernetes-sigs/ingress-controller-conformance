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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/k8s"
	"io/ioutil"
	"net/http"
	"time"
)

// Config contains test suite configuration fields
type Config struct {
	IngressAPIVersion string // IngressAPIVersion for only executing tests which are compatible with this APIVersion.
	UseInsecureHost   string // UseInsecureHost for cleartext requests when the infrastructure under test does not allow for auto-detecting the public IP associated with the Ingress resources.
	UseSecureHost     string // UseSecureHost for secure/encrypted requests when the infrastructure under test does not allow for auto-detecting the public IP associated with the Ingress resources.
}

// Check represents a test case. Checks are named, and must provide a
// description and a Run function. Checks are organized in a hierarchy.
type Check struct {
	Name        string
	Description string
	APIVersions []string // APIVersions for which this test is valid/compatible.

	//Specify either one of RunRequest or Run
	RunRequest *Request                                                    // Given a Request, check assertions without handling the boilerplate of request execution.
	Run        func(check *Check, config Config) (success bool, err error) // Generic Run function for advanced cases.

	checks []*Check // Child checks
	parent *Check   // Parent check
}

// Request allows defining a making and capturing a single HTTP request and build assertions with DoCheck.
type Request struct {
	IngressNamespace string
	IngressName      string

	Path     string
	Hostname string
	Insecure bool

	DoCheck func(*CapturedRequest, *CapturedResponse) (*Assertions, error) // For a given CapturedRequest and CapturedResponse, check some Assertions.
}

// CapturedRequest contains the original HTTP request metadata as received
// by the echoserver handling the test request.
type CapturedRequest struct {
	DownstreamServiceId string `json:"testId"` // DownstreamServiceId field contains the TEST_ID environment variable value of the downstream echoserver.
	Path                string
	Host                string
	Method              string
	Proto               string
	Headers             map[string][]string
}

// CapturedResponse contains the HTTP response metadata from the echoserver
type CapturedResponse struct {
	StatusCode    int
	ContentLength int64
	Proto         string
	Headers       map[string][]string
}

// CaptureRoundTrip will perform an HTTP request and return the CapturedRequest and CapturedResponse tuple
func CaptureRoundTrip(location string, hostOverride string) (*CapturedRequest, *CapturedResponse, error) {
	tr := &http.Transport{
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			// Skip all usual TLS verifications, since we are using a self-signed certificate.
			InsecureSkipVerify: true,
			VerifyPeerCertificate: func(certificates [][]byte, _ [][]*x509.Certificate) error {
				certs := make([]*x509.Certificate, len(certificates))
				for i, asn1Data := range certificates {
					cert, err := x509.ParseCertificate(asn1Data)
					if err != nil {
						return fmt.Errorf("tls: failed to parse certificate from server: " + err.Error())
					}
					certs[i] = cert
				}

				// Verify the certificate Hostname matches the request hostname.
				return certs[0].VerifyHostname(hostOverride)
			},
		},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 3,
	}
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		return nil, nil, err
	}
	if hostOverride != "" {
		req.Host = hostOverride
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	capReq := &CapturedRequest{}
	err = json.NewDecoder(resp.Body).Decode(capReq)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("unexpected response (statuscode: %d, length: %d): %s", resp.StatusCode, len(body), body)
		return nil, nil, err
	}

	capRes := &CapturedResponse{
		resp.StatusCode,
		resp.ContentLength,
		resp.Proto,
		resp.Header,
	}
	return capReq, capRes, nil
}

// Head of Check hierarchy
var AllChecks = &Check{
	Name: "all",
}

// Add child checks
func (c *Check) AddCheck(checks ...*Check) {
	for i, x := range checks {
		if checks[i] == c {
			panic("Checks can't be a child of itself")
		}
		checks[i].parent = c
		c.checks = append(c.checks, x)
	}
}

// List this check and its child's description
func (c *Check) List() {
	if c.Description != "" {
		fmt.Printf("- %s %v\n\t%s\n", c.Name, c.APIVersions, c.Description)
	}
	for _, check := range c.checks {
		check.List()
	}
}

// Run all checks, filtered by name and given a configuration
func (c *Check) Verify(filterOnCheckName string, config Config) (successCount int, failureCount int, err error) {
	if filterOnCheckName != c.Name && filterOnCheckName != "" {
		for _, check := range c.checks {
			s, f, err := check.Verify(filterOnCheckName, config)
			successCount += s
			failureCount += f
			if err != nil {
				fmt.Printf(err.Error())
			}
		}

		return
	}

	if c.Run == nil && c.RunRequest != nil {
		c.Run = func(check *Check, config Config) (bool, error) {
			var scheme = "https"
			var host = config.UseSecureHost
			if c.RunRequest.Insecure {
				scheme = "http"
				host = config.UseInsecureHost
			}
			if host == "" {
				var err error
				namespace := "default"
				if c.RunRequest.IngressNamespace != "" {
					namespace = c.RunRequest.IngressNamespace
				}
				host, err = k8s.GetIngressHost(namespace, c.RunRequest.IngressName)
				if err != nil {
					return false, err
				}
			}

			location := fmt.Sprintf("%s://%s%s", scheme, host, c.RunRequest.Path)
			req, res, err := CaptureRoundTrip(location, c.RunRequest.Hostname)
			if err != nil {
				return false, err
			}
			assertions, err := c.RunRequest.DoCheck(req, res)

			fmt.Print(assertions)
			return assertions != nil && assertions.Passed(), nil
		}
	}

	isAPIVersionCompatible := func() bool {
		for _, v := range c.APIVersions {
			if v == config.IngressAPIVersion {
				return true
			}
		}
		return false
	}

	runChildChecks := true
	if c.Run != nil && isAPIVersionCompatible() {
		fmt.Printf("Running '%s' verifications...\n", c.Name)
		success, err := c.Run(c, config)
		if err != nil {
			fmt.Printf("\t%s\n", err.Error())
		}

		if success {
			successCount++
		} else {
			failureCount++
			runChildChecks = false
			fmt.Printf("\t\033[1;31mCheck failed:\033[0m %s\n", c.Name)
		}
	}

	if runChildChecks {
		for _, check := range c.checks {
			s, f, err := check.Verify("", config)
			if err != nil {
				fmt.Printf(err.Error())
			}
			successCount += s
			failureCount += f
		}
	}

	return
}
