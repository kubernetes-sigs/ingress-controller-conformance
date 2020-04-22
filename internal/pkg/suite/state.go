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

var (
	ingressEndpoint  string
	captureRoundTrip *CaptureRoundTrip
)

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

// CapturedResponse contains the HTTP request and response metadata from a round trip to echoserver
type CaptureRoundTrip struct {
	Request  *CapturedRequest
	Response *CapturedResponse
}
