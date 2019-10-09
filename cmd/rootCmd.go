package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

var (
	VERSION string
)

var rootCmd = &cobra.Command{
	Use:     "ingress-conformance",
	Short:   "Ingress conformance test suite",
	Version: VERSION,
	Long: `Kubernetes ingress controller conformance test suite in Go.
  Complete documentation is available at https://github.com/datawire/ingress-conformance`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
