module sigs.k8s.io/ingress-controller-conformance

go 1.15

require (
	github.com/cucumber/gherkin-go/v11 v11.0.0
	github.com/cucumber/godog v0.10.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0
	golang.org/x/tools v0.0.0-20191227053925-7b8e75db28f4
	k8s.io/api v0.19.2
	k8s.io/apimachinery v0.19.2
	k8s.io/cli-runtime v0.19.2
	k8s.io/client-go v0.19.2
	k8s.io/klog/v2 v2.3.0
	sigs.k8s.io/yaml v1.2.0
)
