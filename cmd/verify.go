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

package cmd

import (
	"fmt"
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/apiversion"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/assets"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/suite"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	verifyCmd.Flags().StringVar(&apiVersionTag, "api-version", "",
		fmt.Sprintf("run conformance tests for a specific api-version %v", apiversion.All))
	if err := verifyCmd.MarkFlagRequired("api-version"); err != nil {
		panic(err)
	}

	rootCmd.AddCommand(verifyCmd)
}

var (
	apiVersionTag = ""
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run Ingress verifications for conformance",
	Long:  "Run Ingress verifications for conformance",
	Run: func(cmd *cobra.Command, args []string) {
		featuresDir := ".features"
		err := assets.RestoreAssets(featuresDir, "features")
		if err != nil {
			panic(err)
		}
		status := godog.RunWithOptions("ingress-controller-conformance", func(s *godog.Suite) {
			suite.FeatureContext(s)
		}, godog.Options{
			Output:      colors.Colored(os.Stdout),
			Format:      "progress",
			Paths:       []string{featuresDir},
			Tags:        apiVersionTag,
			Strict:      true,
			Concurrency: 1,
			Randomize:   -1, // Let godog generate a random seed
		})

		os.Exit(status)
	},
}
