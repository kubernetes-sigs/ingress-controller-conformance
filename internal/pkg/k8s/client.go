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

package k8s

import (
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
	"os"
	"sync"
)

var (
	cachedClient *kubernetes.Clientset
	cachedConfig *api.Config

	once sync.Once
)

func initClient() {
	err := func() error {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
		// use the current context in kubeconfig
		clientConfig, err := kubeconfig.ClientConfig()
		if err != nil {
			return err
		}

		cachedConfig, err = loadingRules.Load()
		if err != nil {
			return err
		}

		cachedClient, err = kubernetes.NewForConfig(clientConfig)
		if err != nil {
			return err
		}

		return nil
	}()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to configure Kubernetes cluster context: %s\n", err)
		os.Exit(1)
	}
}

func Client() *kubernetes.Clientset {
	once.Do(initClient)
	return cachedClient
}

func Config() *api.Config {
	once.Do(initClient)
	return cachedConfig
}

func GetIngressHost(namespace string, ingressName string) (host string, err error) {
	ingressInterface, err := Client().NetworkingV1beta1().Ingresses(namespace).Get(ingressName, v1.GetOptions{})
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
			err = fmt.Errorf("ingresses/status/loadBalancer \"%s\" has no hostname or IP", ingressName)
		}
	} else {
		err = fmt.Errorf("ingresses/status \"%s\" has no load balancer interface; use '--use-insecure-host' "+
			"and '--use-secure-host' if this is a limitation from the infrastructure", ingressName)
	}
	return
}
