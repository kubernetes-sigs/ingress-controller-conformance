module sigs.k8s.io/ingress-controller-conformance

go 1.14

require (
	github.com/cucumber/gherkin-go/v11 v11.0.0
	github.com/cucumber/godog v0.9.1-0.20200517063737-7568b291e4e1
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/iancoleman/orderedmap v0.0.0-20190318233801-ac98e3ecb4b0
	github.com/spf13/cobra v1.0.0
	golang.org/x/text v0.3.2
	golang.org/x/tools v0.0.0-20190920225731-5eefd052ad72
	k8s.io/api v0.18.3
	k8s.io/apimachinery v0.18.3
	k8s.io/cli-runtime v0.18.3
	k8s.io/client-go v0.18.3
	k8s.io/klog v1.0.0
	k8s.io/kubectl v0.18.3
)
