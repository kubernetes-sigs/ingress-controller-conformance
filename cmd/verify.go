package cmd

import (
	"fmt"
	"github.com/datawire/ingress-conformance/internal/pkg/checks"
	"github.com/spf13/cobra"
)

func init() {
	verifyCmd.Flags().StringVarP(&checkName, "check", "c", "", "Verify only this specified check name.")

	rootCmd.AddCommand(verifyCmd)
}

var (
	checkName = ""
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run Ingress verifications for conformance",
	Long:  "Run Ingress verifications for conformance",
	Run: func(cmd *cobra.Command, args []string) {
		config := checks.Config{}
		successCount, failureCount, err := checks.Checks.Verify(checkName, config)

		if err != nil {
			fmt.Printf(err.Error())
		}

		fmt.Printf("%d checks passed! %d failures\n", successCount, failureCount)
	},
}
