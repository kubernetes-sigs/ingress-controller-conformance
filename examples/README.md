# Ingress conformance tests examples

Testing conformance of an ingress controller should only require a Kubernetes cluster with a running pod of the ingress controller to use.

This directory contains two examples to show the content of a github repository prepared to run the conformance suite using Github actions.

### kind

Shows how to test an ingress controller that do not requires cloud resources, like [ingress-nginx](https://github.com/kubernetes/ingress-nginx) using [kind](https://kind.sigs.k8s.io/docs/user/ingress/) to bootstrap a Kubernetes cluster

### gce

Shows how to test an ingress controller provided by a cloud vendor, Google Cloud in this case using [kube-up](https://kubernetes.io/docs/tasks/tools/install-kubectl/#verifying-kubectl-configuration) to test [ingress-gce](https://github.com/kubernetes/ingress-gce)
