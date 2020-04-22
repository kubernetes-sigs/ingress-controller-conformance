module github.com/kubernetes-sigs/ingress-controller-conformance

go 1.13

require (
	github.com/cucumber/godog v0.9.0
	github.com/cucumber/messages-go/v10 v10.0.3
	github.com/go-bindata/go-bindata v3.1.2+incompatible
	github.com/spf13/cobra v0.0.5
	k8s.io/apimachinery v0.18.0-alpha.2
	k8s.io/cli-runtime v0.18.0-alpha.2
	k8s.io/client-go v0.18.0-alpha.2
	k8s.io/kubectl v0.18.0-alpha.2
)
