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
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"github.com/cucumber/messages-go/v10"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/k8s"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func iHaveAnIngressNamedInTheNamespace(ingressName string, namespace string) error {
	// Fixtures are externalized.
	// We assume all resources were applied with `ingress-controller-conformance apply` prior to executing the test suite

	host, err := k8s.GetIngressHost(namespace, ingressName)
	if err != nil {
		return err
	}
	ingressEndpoint = host

	return nil
}

func iSendARequest(requestMethod string, requestURL string) error {
	if ingressEndpoint == "" {
		return fmt.Errorf("undefined ingress host location")
	}

	parsedRequestURL, err := url.Parse(requestURL)
	if err != nil {
		return err
	}

	requestHost := parsedRequestURL.Host    // record the true request Host to set headers
	parsedRequestURL.Host = ingressEndpoint // replace the request Host with the ingress enpoint
	requestLocation := parsedRequestURL.String()

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
				return certs[0].VerifyHostname(requestHost)
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
	req, err := http.NewRequest(requestMethod, requestLocation, nil)
	if err != nil {
		return err
	}
	req.Host = requestHost // set the Host header according to our desired requestURL

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	capReq := &CapturedRequest{}
	err = json.NewDecoder(resp.Body).Decode(capReq)
	if err != nil {
		body, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("unexpected response (statuscode: %d, length: %d): %s", resp.StatusCode, len(body), body)
		return err
	}

	capRes := &CapturedResponse{
		resp.StatusCode,
		resp.ContentLength,
		resp.Proto,
		resp.Header,
	}

	captureRoundTrip = &CaptureRoundTrip{
		Request:  capReq,
		Response: capRes,
	}

	return nil
}

func theResponseStatuscodeMustBe(expectedResponseStatusCode int) error {
	if captureRoundTrip.Response.StatusCode != expectedResponseStatusCode {
		return fmt.Errorf("expected the status code to be %v but was %v", expectedResponseStatusCode, captureRoundTrip.Response.StatusCode)
	}
	return nil
}

func theResponseProtoMustBe(expectedResponseProto string) error {
	if captureRoundTrip.Response.Proto != expectedResponseProto {
		return fmt.Errorf("expected the response protocol to be %v but was %v", expectedResponseProto, captureRoundTrip.Response.Proto)
	}
	return nil
}

func theResponseMustBeServedByTheService(expectedDownstreamService string) error {
	if captureRoundTrip.Request.DownstreamServiceId != expectedDownstreamService {
		return fmt.Errorf("expected the responding service would be %v but was %v", expectedDownstreamService, captureRoundTrip.Request.DownstreamServiceId)
	}
	return nil
}

func theResponseHeadersMustContainHeaderWithMatchingValue(headers *messages.PickleStepArgument_PickleTable) error {
	if len(headers.Rows) < 1 {
		return fmt.Errorf("expected a table with at least one <key>|<value> row")
	}

	for index, row := range headers.Rows {
		if index == 0 {
			continue // Skip the header row
		}

		key := row.Cells[0].Value   // <key>
		value := row.Cells[1].Value // <value>

		headerValues := captureRoundTrip.Response.Headers[key]
		if headerValues == nil {
			return fmt.Errorf("expected the response headers to contain %v", key)
		}
		if value != "*" {
			found := false
			for _, v := range headerValues {
				if v == value {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("expected the response headers to contain %v; with a value of %s but was %v", key, value, headerValues)
			}
		}
	}

	return nil
}

func theRequestHostMustBe(expectedRequestHost string) error {
	if captureRoundTrip.Request.Host != expectedRequestHost {
		return fmt.Errorf("expected the request host would be %v but was %v", expectedRequestHost, captureRoundTrip.Request.Host)
	}
	return nil
}

func theRequestMethodMustBe(expectedRequestMethod string) error {
	if captureRoundTrip.Request.Method != expectedRequestMethod {
		return fmt.Errorf("expected the request method would be %v but was %v", expectedRequestMethod, captureRoundTrip.Request.Method)
	}
	return nil
}

func theRequestPathMustBe(expectedRequestPath string) error {
	if captureRoundTrip.Request.Path != expectedRequestPath {
		return fmt.Errorf("expected the request path would be %v but was %v", expectedRequestPath, captureRoundTrip.Request.Path)
	}
	return nil
}

func theRequestProtoMustBe(expectedRequestProto string) error {
	if captureRoundTrip.Request.Proto != expectedRequestProto {
		return fmt.Errorf("expected the request protocol would be %v but was %v", expectedRequestProto, captureRoundTrip.Request.Proto)
	}
	return nil
}

func theRequestHeadersMustContainHeaderWithMatchingValue(headers *messages.PickleStepArgument_PickleTable) error {
	if len(headers.Rows) < 1 {
		return fmt.Errorf("expected a table with at least one <key>|<value> row")
	}

	for index, row := range headers.Rows {
		if index == 0 {
			continue // Skip the header row
		}

		key := row.Cells[0].Value   // <key>
		value := row.Cells[1].Value // <value>

		headerValues := captureRoundTrip.Request.Headers[key]
		if headerValues == nil {
			return fmt.Errorf("expected the request headers to contain %v", key)
		}
		if value != "*" {
			found := false
			for _, v := range headerValues {
				if v == value {
					found = true
				}
			}
			if !found {
				return fmt.Errorf("expected the request headers to contain %v; with a value of %s but was %v", key, value, headerValues)
			}
		}
	}

	return nil
}
