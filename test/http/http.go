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

package http

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// CapturedRequest contains the original HTTP request metadata as received
// by the echoserver handling the test request.
type CapturedRequest struct {
	Path    string              `json:"path"`
	Host    string              `json:"host"`
	Method  string              `json:"method"`
	Proto   string              `json:"proto"`
	Headers map[string][]string `json:"headers"`

	Namespace string `json:"namespace"`
	Ingress   string `json:"ingress"`
	Service   string `json:"service"`
}

// CapturedResponse contains the HTTP response metadata from the echoserver.
type CapturedResponse struct {
	StatusCode    int
	ContentLength int64
	Proto         string
	Headers       map[string][]string
	TLSHostname   string
}

// CaptureRoundTrip will perform an HTTP request and return the CapturedRequest and CapturedResponse tuple
func CaptureRoundTrip(method, scheme, hostname, path, location string) (*CapturedRequest, *CapturedResponse, error) {
	var capturedTLSHostname string
	tr := &http.Transport{
		DisableCompression: true,
		TLSClientConfig: &tls.Config{
			// Skip all usual TLS verifications, since we are using self-signed certificates.
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
				capturedTLSHostname = hostname
				return certs[0].VerifyHostname(hostname)
			},
		},
	}
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Second * 3,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	url := fmt.Sprintf("%s://%s/%s", scheme, location, strings.TrimPrefix(path, "/"))
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, nil, err
	}
	if hostname != "" {
		req.Host = hostname
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
		capturedTLSHostname,
	}
	return capReq, capRes, nil
}
