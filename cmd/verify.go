package cmd

import (
	"fmt"
	"github.com/datawire/ingress-conformance/internal/pkg"
	"github.com/spf13/cobra"
)

func init() {
	verifyCmd.Flags().StringVarP(&checkName, "check", "c", "", "Verify only this specified check name.")
	verifyCmd.Flags().StringVar(&host, "host", "localhost:80", "Target hostname or IP/port against which to run checks.")

	rootCmd.AddCommand(verifyCmd)
}

var (
	checkName = ""
	host      = ""
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run all Ingress verifications for conformance",
	Long:  "Run all Ingress verifications for conformance",
	Run: func(cmd *cobra.Command, args []string) {
		config := pkg.Config{Host: host}

		successCount, failureCount, err := pkg.Checks.Verify(checkName, config)

		if err != nil {
			fmt.Printf(err.Error())
		}

		fmt.Printf("%d checks passed! %d failures\n", successCount, failureCount)
	},
}
