module sigs.k8s.io/ingress-controller-conformance

go 1.14

require (
	github.com/cucumber/gherkin-go/v11 v11.0.0
	github.com/cucumber/godog v0.10.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0
	github.com/spf13/cobra v1.0.0
	golang.org/x/text v0.3.3
	golang.org/x/tools v0.0.0-20200602230032-c00d67ef29d0
	k8s.io/api v0.19.0-rc.0
	k8s.io/apimachinery v0.19.0-rc.0
	k8s.io/cli-runtime v0.19.0-rc.0
	k8s.io/client-go v0.19.0-rc.0
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.19.0-rc.0
)
