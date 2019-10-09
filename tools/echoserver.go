package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

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

	response.Header().Set("Content-Type", "application/json")
	response.Write(js)
}
