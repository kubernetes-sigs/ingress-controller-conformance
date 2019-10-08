package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

func init() {
	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
	// use the current context in kubeconfig
	clientConfig, err := kubeconfig.ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	Config, err = loadingRules.Load()
	if err != nil {
		panic(err.Error())
	}

	Client, err = kubernetes.NewForConfig(clientConfig)
	if err != nil {
		panic(err.Error())
	}
}

var (
	Client *kubernetes.Clientset
	Config *api.Config
)
