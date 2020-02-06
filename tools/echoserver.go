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

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

// RequestAssertions contains the HTTP response which can be asserted
// by checks.CapturedRequest
type RequestAssertions struct {
	TestId  string
	Path    string
	Host    string
	Method  string
	Proto   string
	Headers map[string][]string
}

type preserveSlashes struct {
	mux http.Handler
}

func (s *preserveSlashes) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// preserve the original URL Path, don't let Go normalize it
	r.URL.Path = strings.Replace(r.URL.Path, "//", "/", -1)
	s.mux.ServeHTTP(w, r)
}

var TestId string

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func main() {
	c := struct {
		Port   string
		TestId string
	}{
		Port:   getEnv("PORT", "3000"),
		TestId: getEnv("TEST_ID", ""),
	}
	TestId = c.TestId

	fmt.Printf("Starting server, listening on port %s\n", c.Port)
	fmt.Printf("Reporting TestId '%s'\n", TestId)
	httpMux := http.NewServeMux()
	httpMux.HandleFunc("/", RequestHandler)
	err := http.ListenAndServe(fmt.Sprintf(":%s", c.Port), &preserveSlashes{httpMux})
	if err != nil {
		panic(fmt.Sprintf("Failed to start listening: %s\n", err.Error()))
	}
}

func RequestHandler(response http.ResponseWriter, request *http.Request) {
	fmt.Printf("Echoing back request made to %s to client (%s)\n", request.RequestURI, request.RemoteAddr)

	requestAssertions := RequestAssertions{
		TestId,
		request.RequestURI,
		request.Host,
		request.Method,
		request.Proto,
		request.Header,
	}

	js, err := json.Marshal(requestAssertions)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	// Go libs have no gzip support for HTTP responses, sending it uncompressed.
	response.Header().Set("Content-Type", "application/json")
	response.Write(js)
}
