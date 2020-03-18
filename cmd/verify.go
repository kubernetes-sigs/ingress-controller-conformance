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
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/apiversion"
	"github.com/kubernetes-sigs/ingress-controller-conformance/internal/pkg/checks"
	"github.com/spf13/cobra"
	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"os"
	"time"
)

func init() {
	verifyCmd.Flags().StringVar(&verifyIngressAPIVersion,
		"api-version", "",
		fmt.Sprintf("verify using assertions for the given apiVersion %s", apiversion.All))
	verifyCmd.Flags().StringVarP(&checkName, "check", "c", "",
		"verify only this specified check name")
	verifyCmd.Flags().StringVar(&useInsecureHost, "use-insecure-host", "",
		"endpoint to use for testing cleartext requests, such as 'localhost:8080', when the Ingress"+
			" resources cannot be associated with a load balancer interface due to infrastructure restrictions")
	verifyCmd.Flags().StringVar(&useSecureHost, "use-secure-host", "",
		"endpoint to use for testing secure/encrypted requests, such as 'localhost:8443', when the Ingress"+
			" resources cannot be associated with a load balancer interface due to infrastructure restrictions")

	if err := verifyCmd.MarkFlagRequired("api-version"); err != nil {
		panic(err)
	}
	rootCmd.AddCommand(verifyCmd)

	_ = message.Set(language.English, "%d success",
		plural.Selectf(1, "%d",
			"=0", "\033[1;36mNo checks passed...\033[0m",
			"=1", "\033[1;34m1 check passed,\033[0m",
			"other", "\033[1;34m%d checks passed!\033[0m",
		),
	)
	_ = message.Set(language.English, "%d failure",
		plural.Selectf(1, "%d",
			"=0", "\033[1;34mNo failures!\033[0m",
			"=1", "\033[1;31m1 failure\033[0m",
			"other", "\033[1;31m%d failures!\033[0m",
		),
	)
}

var (
	verifyIngressAPIVersion = ""
	checkName               = ""
	useInsecureHost         = ""
	useSecureHost           = ""
)

var verifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run Ingress verifications for conformance",
	Long:  "Run Ingress verifications for conformance",
	Run: func(cmd *cobra.Command, args []string) {
		start := time.Now()

		config := checks.Config{
			IngressAPIVersion: verifyIngressAPIVersion,
			UseInsecureHost:   useInsecureHost,
			UseSecureHost:     useSecureHost,
		}
		successCount, failureCount, err := checks.AllChecks.Verify(checkName, config)
		if err != nil {
			fmt.Printf(err.Error())
		}

		elapsed := time.Since(start)

		p := message.NewPrinter(language.English)
		fmt.Printf("\n--- Verification completed ---\nAPIVersion: %s\n%s %s\nin %s\n",
			verifyIngressAPIVersion,
			p.Sprintf("%d success", successCount),
			p.Sprintf("%d failure", failureCount),
			elapsed)

		// should exit with a non-zero status if we have test failures
		if err != nil || failureCount > 0 {
			os.Exit(1)
		}
	},
}
