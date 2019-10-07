package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func init() {
	rootCmd.AddCommand(contextCmd)
}

var contextCmd = &cobra.Command{
	Use:   "context",
	Short: "Print the current Kubernetes context and server version",
	Long:  "Print the current Kubernetes context and server version then exits",
	Run: func(cmd *cobra.Command, args []string) {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()

		config, err := loadingRules.Load()
		if err != nil {
			panic(err.Error())
		}
		kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})

		// use the current context in kubeconfig
		clientConfig, err := kubeconfig.ClientConfig()
		if err != nil {
			panic(err.Error())
		}

		// create the clientset
		clientset, err := kubernetes.NewForConfig(clientConfig)
		if err != nil {
			panic(err.Error())
		}
		serverVersion, err := clientset.ServerVersion()
		if err != nil {
			panic(err.Error())
		}

		fmt.Printf("Targetting Kubernetes cluster under active context '%s'\n", config.CurrentContext)
		fmt.Printf("The target Kubernetes cluster is running verion %s\n", serverVersion)
	},
}
