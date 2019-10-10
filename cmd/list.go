package cmd

import (
	"github.com/datawire/ingress-conformance/internal/pkg/checks"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Ingress verifications",
	Long:  "List all Ingress verifications",
	Run: func(cmd *cobra.Command, args []string) {
		checks.Checks.List()
	},
}
