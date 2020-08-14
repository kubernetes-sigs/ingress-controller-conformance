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

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	networking "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func applyDeployment(kubeClientSet kubernetes.Interface, namespace string, deployment *appsv1.Deployment) (*appsv1.Deployment, error) {
	existing, err := kubeClientSet.AppsV1().Deployments(namespace).Get(context.TODO(), deployment.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := kubeClientSet.AppsV1().Deployments(namespace).Create(context.TODO(), deployment, metav1.CreateOptions{})
		return actual, err
	}

	if err != nil {
		return nil, err
	}

	return existing, nil
}

func applySecret(kubeClientSet kubernetes.Interface, namespace string, secret *corev1.Secret) (*corev1.Secret, error) {
	existing, err := kubeClientSet.CoreV1().Secrets(namespace).Get(context.TODO(), secret.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := kubeClientSet.CoreV1().Secrets(namespace).Create(context.TODO(), secret, metav1.CreateOptions{})
		return actual, err
	}

	if err != nil {
		return nil, err
	}

	return existing, nil
}

func applyService(kubeClientSet kubernetes.Interface, namespace string, service *corev1.Service) (*corev1.Service, error) {
	existing, err := kubeClientSet.CoreV1().Services(namespace).Get(context.TODO(), service.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := kubeClientSet.CoreV1().Services(namespace).Create(context.TODO(), service, metav1.CreateOptions{})
		return actual, err
	}

	if err != nil {
		return nil, err
	}

	return existing, nil
}

func applyIngress(kubeClientSet kubernetes.Interface, namespace string, ingress *networking.Ingress) (*networking.Ingress, error) {
	existing, err := kubeClientSet.NetworkingV1().Ingresses(namespace).Get(context.TODO(), ingress.Name, metav1.GetOptions{})
	if apierrors.IsNotFound(err) {
		actual, err := kubeClientSet.NetworkingV1().Ingresses(namespace).Create(context.TODO(), ingress, metav1.CreateOptions{})
		return actual, err
	}

	if err != nil {
		return nil, err
	}

	return existing, nil
}
