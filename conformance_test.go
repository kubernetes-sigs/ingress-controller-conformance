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

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/cucumber/godog"
	"k8s.io/apimachinery/pkg/util/sets"
	clientset "k8s.io/client-go/kubernetes"
	"k8s.io/klog"

	"sigs.k8s.io/ingress-controller-conformance/test/conformance/defaultbackend"
	"sigs.k8s.io/ingress-controller-conformance/test/files"
	"sigs.k8s.io/ingress-controller-conformance/test/kubernetes"
)

var (
	godogFormat        string
	godogTags          string
	godogStopOnFailure bool
	godogNoColors      bool
	godogOutput        string

	manifests string
)

func TestMain(m *testing.M) {
	// register flags from klog (client-go verbose logging)
	klog.InitFlags(nil)

	flag.StringVar(&godogFormat, "format", "pretty", "Set godog format to use. Valid values are pretty and cucumber")
	flag.StringVar(&godogTags, "tags", "", "Tags for conformance test")
	flag.BoolVar(&godogStopOnFailure, "stop-on-failure ", false, "Stop when failure is found")
	flag.BoolVar(&godogNoColors, "no-colors", false, "Disable colors in godog output")
	flag.StringVar(&manifests, "manifests", "./manifests", "Directory where manifests for test applications or scenerarios are located")
	flag.StringVar(&godogOutput, "output-directory", ".", "Output directory for test reports")
	flag.StringVar(&kubernetes.IngressClassValue, "ingress-class", "conformance", "Sets the value of the annotation kubernetes.io/ingress.class in Ingress definitions")

	flag.Parse()

	validFormats := sets.NewString("cucumber", "pretty")
	if !validFormats.Has(godogFormat) {
		klog.Fatalf("the godog format '%v' is not supported", godogFormat)
	}

	manifestsPath, err := filepath.Abs(manifests)
	if err != nil {
		klog.Fatal(err)
	}

	if !files.IsDir(manifestsPath) {
		klog.Fatalf("The specified value in the flag --manifests-directory (%v) is not a directory", manifests)
	}

	kubernetes.KubeClient, err = setupSuite()
	if err != nil {
		klog.Fatal(err)
	}

	if err := kubernetes.CleanupNamespaces(kubernetes.KubeClient); err != nil {
		klog.Fatalf("error deleting temporal namespaces: %v", err)
	}

	os.Exit(m.Run())
}

func setupSuite() (*clientset.Clientset, error) {
	c, err := kubernetes.LoadClientset()
	if err != nil {
		return nil, fmt.Errorf("error loading client: %v", err)
	}

	dc := c.DiscoveryClient

	serverVersion, serverErr := dc.ServerVersion()
	if serverErr != nil {
		return nil, fmt.Errorf("unexpected server error retrieving version: %v", serverErr)
	}

	if serverVersion != nil {
		// TODO: check minimum k8s version?
		klog.Infof("kube-apiserver version: %s", serverVersion.GitVersion)
	}

	return c, nil
}

// Generated code. DO NOT EDIT.
var (
	features = map[string]func(*godog.Suite){
		"features/default_backend.feature": defaultbackend.FeatureContext,
	}
)

func TestSuite(t *testing.T) {
	for feature, featureContext := range features {
		err := testFeature(feature, featureContext)
		if err != nil {
			t.Fatal(err)
		}
	}
}

func testFeature(feature string, featureContext func(*godog.Suite)) error {
	var testOutput io.Writer
	// default output is stdout
	testOutput = os.Stdout

	if godogFormat == "cucumber" {
		rf := path.Join(godogOutput, fmt.Sprintf("%v-report.json", filepath.Base(feature)))
		file, err := os.Create(rf)
		if err != nil {
			return fmt.Errorf("error creating report file %v: %w", rf, err)
		}

		defer file.Close()

		writer := bufio.NewWriter(file)
		defer writer.Flush()

		testOutput = writer
	}

	exitCode := godog.RunWithOptions("conformance", func(s *godog.Suite) {
		featureContext(s)
	}, godog.Options{
		Format:        godogFormat,
		Paths:         []string{feature},
		Tags:          godogTags,
		StopOnFailure: godogStopOnFailure,
		NoColors:      godogNoColors,
		Output:        testOutput,
		Concurrency:   1, // do not run tests concurrently
	})
	if exitCode > 0 {
		return fmt.Errorf("unexpected exit code running test: %v", exitCode)
	}

	return nil
}
