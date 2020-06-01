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

	"github.com/spf13/cobra"
	"sigs.k8s.io/ingress-controller-conformance/internal/pkg/k8s"
)

func init() {
	rootCmd.AddCommand(contextCmd)
}

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Print the current Kubernetes context, server version, and supported Ingress APIVersions",
	Long:  "Print the current Kubernetes context, server version, and supported Ingress APIVersions then exits",
	Run: func(cmd *cobra.Command, args []string) {
		serverVersion, err := k8s.Client().ServerVersion()
		if err != nil {
			panic(err.Error())
		}

		resources, err := k8s.Client().ServerPreferredResources()
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("Using active Kubernetes context '%s'\n", k8s.Config().CurrentContext)
		fmt.Printf("The target Kubernetes cluster is running verion %s\n", serverVersion)
		for _, resource := range resources {
			for _, apiResource := range resource.APIResources {
				if apiResource.Kind == "Ingress" {
					fmt.Printf("  Supports Ingress kind APIVersion: '%s'\n", resource.GroupVersion)
				}
			}
		}
	},
}
