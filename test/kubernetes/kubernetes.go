/*
Copyright 2020 The Kubernetes Authors.

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

package kubernetes

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	apierrs "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	utilnet "k8s.io/apimachinery/pkg/util/net"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	clientset "k8s.io/client-go/kubernetes"
	restclient "k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"sigs.k8s.io/yaml"

	// ensure auth plugins are loaded
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

const (
	// IngressClassKey indicates the class of an Ingress to be used
	// when determining which controller should implement the Ingress
	IngressClassKey = "kubernetes.io/ingress.class"
)

// IngressClassValue sets the value of the class of Ingresses
var IngressClassValue string

// KubeClient Kubernetes API client
var KubeClient *kubernetes.Clientset

// LoadClientset returns clientset for connecting to kubernetes clusters.
func LoadClientset() (*clientset.Clientset, error) {
	config, err := restclient.InClusterConfig()
	if err != nil {
		// Attempt to use local KUBECONFIG
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &clientcmd.ConfigOverrides{})
		// use the current context in kubeconfig
		var err error

		config, err = kubeconfig.ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	// TODO: add version information?
	config.UserAgent = fmt.Sprintf(
		"%s (%s/%s) ingress-conformance",
		filepath.Base(os.Args[0]),
		runtime.GOOS,
		runtime.GOARCH,
	)

	client, err := clientset.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// NewNamespace creates a new namespace using ingress-conformance- as prefix.
func NewNamespace(c kubernetes.Interface) (string, error) {
	ns := &corev1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "ingress-conformance-",
			Labels: map[string]string{
				"app.kubernetes.io/name": "ingress-conformance",
			},
		},
	}

	var err error

	ns, err = c.CoreV1().Namespaces().Create(context.TODO(), ns, metav1.CreateOptions{})
	if err != nil {
		return "", fmt.Errorf("unable to create namespace: %v", err)
	}

	return ns.Name, nil
}

// DeleteNamespace deletes a namespace and all the objects inside
func DeleteNamespace(c kubernetes.Interface, namespace string) error {
	grace := int64(0)
	pb := metav1.DeletePropagationBackground

	return c.CoreV1().Namespaces().Delete(context.TODO(), namespace, metav1.DeleteOptions{
		GracePeriodSeconds: &grace,
		PropagationPolicy:  &pb,
	})
}

// CleanupNamespaces removes namespaces created by conformance tests
func CleanupNamespaces(c kubernetes.Interface) error {
	namespaces, err := c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{
		LabelSelector: "app.kubernetes.io/name=ingress-conformance",
	})

	if err != nil {
		return err
	}

	for _, namespace := range namespaces.Items {
		err := DeleteNamespace(c, namespace.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

// NewIngress creates a new ingress
func NewIngress(c kubernetes.Interface, ingress *networking.Ingress) error {
	_, err := c.NetworkingV1().Ingresses(ingress.Namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
	if err != nil {
		return fmt.Errorf("creating Ingress: %w", err)
	}

	return nil
}

// IngressFromSpec deserializes an Ingress definition using an IngressSpec
func IngressFromSpec(name, namespace, ingressSpec string) (*networking.Ingress, error) {
	if namespace == metav1.NamespaceNone || namespace == metav1.NamespaceDefault {
		return nil, fmt.Errorf("Ingress definitions in the default namespace are not allowed (%v)", namespace)
	}

	ingress := &networking.Ingress{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
	}

	if err := yaml.Unmarshal([]byte(ingressSpec), &ingress.Spec); err != nil {
		return nil, fmt.Errorf("deserializing Ingress from spec: %w", err)
	}

	return ingress, nil
}

// IngressFromManifest deserializes an Ingress definition using an Ingress
func IngressFromManifest(namespace, manifest string) (*networking.Ingress, error) {
	if namespace == metav1.NamespaceNone || namespace == metav1.NamespaceDefault {
		return nil, fmt.Errorf("Ingress definitions in the default namespace are not allowed (%v)", namespace)
	}

	ingress := &networking.Ingress{}
	if err := yaml.Unmarshal([]byte(manifest), &ingress); err != nil {
		return nil, fmt.Errorf("deserializing Ingress from manifest: %w", err)
	}

	ingress.SetNamespace(namespace)
	return ingress, nil
}

const (
	// ingressWaitInterval time to wait between checks for a condition
	ingressWaitInterval = 5 * time.Second
)

var (
	// WaitForIngressAddressTimeout maximum wait time for valid ingress status value
	WaitForIngressAddressTimeout = 5 * time.Minute
)

// WaitForIngressAddress waits for the Ingress to acquire an address.
func WaitForIngressAddress(c clientset.Interface, namespace, name string) (string, error) {
	var address string
	err := wait.PollImmediate(ingressWaitInterval, WaitForIngressAddressTimeout, func() (bool, error) {
		ipOrNameList, err := getIngressAddress(c, namespace, name)
		if err != nil || len(ipOrNameList) == 0 {
			if isRetryableAPIError(err) {
				return false, nil
			}

			return false, err
		}

		address = ipOrNameList[0]
		return true, nil
	})

	if err != nil {
		return "", fmt.Errorf("waiting for ingress status update: %w", err)
	}

	return address, nil
}

// getIngressAddress returns the ips/hostnames associated with the Ingress.
func getIngressAddress(c clientset.Interface, ns, name string) ([]string, error) {
	ing, err := c.NetworkingV1().Ingresses(ns).Get(context.TODO(), name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	var addresses []string

	for _, a := range ing.Status.LoadBalancer.Ingress {
		if a.IP != "" {
			addresses = append(addresses, a.IP)
		}

		if a.Hostname != "" {
			addresses = append(addresses, a.Hostname)
		}
	}

	return addresses, nil
}

// isRetryableAPIError checks if an API error allows retries or not
func isRetryableAPIError(err error) bool {
	// These errors may indicate a transient error that we can retry in tests.
	if apierrs.IsInternalError(err) || apierrs.IsTimeout(err) || apierrs.IsServerTimeout(err) ||
		apierrs.IsTooManyRequests(err) || utilnet.IsProbableEOF(err) || utilnet.IsConnectionReset(err) {
		return true
	}

	// If the error sends the Retry-After header, we respect it as an explicit confirmation we should retry.
	if _, shouldRetry := apierrs.SuggestsClientDelay(err); shouldRetry {
		return true
	}

	return false
}
