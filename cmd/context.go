package cmd

import (
	"fmt"
	"github.com/datawire/ingress-conformance/internal/pkg/k8s"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(contextCmd)
}

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Print the current Kubernetes context, server version, and supported Ingress APIVersions",
	Long:  "Print the current Kubernetes context, server version, and supported Ingress APIVersions then exits",
	Run: func(cmd *cobra.Command, args []string) {
		serverVersion, err := k8s.Client.ServerVersion()
		if err != nil {
			panic(err.Error())
		}

		resources, err := k8s.Client.ServerPreferredResources()
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("Using active Kubernetes context '%s'\n", k8s.Config.CurrentContext)
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
