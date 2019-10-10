package k8s

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
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

func GetIngressHost(namespace string, ingressName string) (host string, err error) {
	ingressInterface, err := Client.NetworkingV1beta1().Ingresses(namespace).Get(ingressName, v1.GetOptions{})
	if err != nil {
		return
	}
	if ingressInterface.Status.LoadBalancer.Ingress != nil {
		ingressInterface := ingressInterface.Status.LoadBalancer.Ingress[0]
		if ingressInterface.Hostname != "" {
			host = ingressInterface.Hostname
		} else {
			host = ingressInterface.IP
		}
		if host == "" {
			err = fmt.Errorf("ingresses.networking.k8s.io \"%s\" has no hostname or IP", ingressName)
		}
	} else {
		err = fmt.Errorf("ingresses.networking.k8s.io \"%s\" has no load balancer interface", ingressName)
	}
	return
}
