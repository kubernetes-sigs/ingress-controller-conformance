package cmd

import (
	"fmt"
	"github.com/datawire/ingress-conformance/internal/pkg/checks"
	"github.com/spf13/cobra"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"time"
)

func init() {
	verifyCmd.Flags().StringVarP(&checkName, "check", "c", "", "Verify only this specified check name.")

	rootCmd.AddCommand(verifyCmd)

	_ = message.Set(language.English, "%d success",
		plural.Selectf(1, "%d",
			"=0", "No checks passed...",
			"=1", "1 check passed,",
			"other", "%d checks passed!",
		),
	)
	_ = message.Set(language.English, "%d failure",
		plural.Selectf(1, "%d",
			"=0", "No failures!",
			"=1", "1 failure",
			"other", "%d failures!",
		),
	)
}

var (
	checkName = ""
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run Ingress verifications for conformance",
	Long:  "Run Ingress verifications for conformance",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		config := checks.Config{}
		successCount, failureCount, err := checks.Checks.Verify(checkName, config)
		if err != nil {
			fmt.Printf(err.Error())
		}

		elapsed := time.Since(start)

		p := message.NewPrinter(language.English)
		fmt.Printf("--- Verification completed ---\n%s %s\nin %s\n",
			p.Sprintf("%d success", successCount),
			p.Sprintf("%d failure", failureCount),
			elapsed)
	},
}
