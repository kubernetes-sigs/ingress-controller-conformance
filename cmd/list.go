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
	"github.com/cucumber/godog"
	"github.com/cucumber/godog/colors"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/suite"
	"github.com/spf13/cobra"
	"os"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Ingress verification step definitions",
	Long:  "List all Ingress verification step definitions",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godog.RunWithOptions("ingress-controller-conformance", func(s *godog.Suite) {
			suite.FeatureContext(s)
		}, godog.Options{
			Output:              colors.Colored(os.Stdout),
			ShowStepDefinitions: true,
		})
	},
}
