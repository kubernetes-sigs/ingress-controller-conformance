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
	Short: "Print the current Kubernetes context and server version",
	Long:  "Print the current Kubernetes context and server version then exits",
	Run: func(cmd *cobra.Command, args []string) {
		serverVersion, err := k8s.Client.ServerVersion()
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("Using active Kubernetes context '%s'\n", k8s.Config.CurrentContext)
		fmt.Printf("The target Kubernetes cluster is running verion %s\n", serverVersion)
	},
}
