# ingress-controller-conformance

The goal of this project is to act as an executable specification in the form of a test suite, implementing a standard, actively maintained ingress-controller conformance specification.

The conformance test suite will both ensure consistency across implementations, as well as simplify the work needed for other implementations to conform to the specification. The test suite can also be viewed through human readable descriptions of what it is testing so that implementers can understand the tests without reading source code.

Currently, the `ingress-controller-conformance` supports the `Ingress` resource from [`networking.k8s.io/v1beta1` API](https://kubernetes.io/docs/concepts/services-networking/ingress/), with the desire to serve as a benchmark for the Ingress resource under the [`networking.k8s.io/v1` API](https://github.com/kubernetes/enhancements/blob/master/keps/sig-network/20190125-ingress-api-group.md) and the [Ingress/Service V2 evolution](https://kubernetes-sigs.github.io/service-apis/) proposals.

## Coverage

The current suite of implemented tests covers the following features of the `Ingress` resource:
- Plain text HTTP/1.1 requests
- Exact and wildcard Host rules
- Prefix path matches rules
- No rules default-backend delegation

Future tests should align with the Ingress resource specification, and support:
- HTTP requests without any matching ingress rules should result in a standard `404 Not Found` responses. Currently, all unmatched requests will fallback to the default-backend.
- Multiple exposed ingress addresses and ports.
- Ingress resources with backend services across different namespaces.
- HTTPS with SSL termination and SNI.
- Support different path match modes: "exact", "prefix", and "regex" path types.
- Other TCP protocols backend types.
- Load-balancing between multiple upstream instances.
- Assert commonly implemented extension points in various ingress-controllers such as gzip compression and keep-alive connections.

## Running

### CLI

At the moment, `ingress-controller-conformance` does not apply modifications to the running resources in the target Kubernetes cluster.

You must manually install and setup you environment targeted by the `kubectl config current-context`:
1. Apply the backing ingress-controller implementation. Samples are available under `examples/`.
1. Apply all, or a subset, of the ingress and service resources found under `deployments/`. Files under the `deployments` folder correspond to implemented check names.

#### Apply

The `ingress-controller-conformance` tool embeds copies of the Kubernetes resources that are used in the conformance checks.
The "apply" command applies these resources to your Kuberenetes cluster using your `kubectl` current-context.

```
$ ./ingress-controller-conformance apply
ingress.networking.k8s.io/host-rules created
service/host-rules created
deployment.apps/host-rules created
ingress.networking.k8s.io/path-rules created
service/path-rules-foo created
deployment.apps/path-rules-foo created
service/path-rules-bar created
deployment.apps/path-rules-bar created
ingress.networking.k8s.io/single-service created
service/single-service created
deployment.apps/single-service created
```

#### Context

```
$ ./ingress-controller-conformance context
Using active Kubernetes context 'docker-desktop'
The target Kubernetes cluster is running verion v1.14.6
  Supports Ingress kind APIVersion: 'extensions/v1beta1'
  Supports Ingress kind APIVersion: 'networking.k8s.io/v1beta1'
```

#### List

Lists, in a human-readable form, all Ingress verifications
```
$ ./ingress-controller-conformance list
- Ingress with host rule should send traffic to the correct backend service (host-rules)
- [...]
- [...]
```

#### Verify

Execute a series of assertions on the deployed ingress-controller, using your `kubectl` current-context.
```
$ ./ingress-controller-conformance verify
Running 'all' verifications...
Running 'host-rules' verifications...
Running 'path-rules' verifications...
Running 'path-rules-foo' verifications...
        1) Assertion failed: Expected the request path would be '/foo' but was '/'
  Check failed: path-rules-foo
Running 'path-rules-foo-trailing' verifications...
        1) Assertion failed: Expected the request path would be '/foo/' but was '//'
  Check failed: path-rules-foo-trailing
[...]
--- Verification completed ---
3 checks passed! 4 failures!
in 1.777148914s
```

#### Help

```
$ ./ingress-controller-conformance help

Kubernetes ingress controller conformance test suite in Go.
  Complete documentation is available at https://github.com/kubernetes-sigs/ingress-controller-conformance

Usage:
  ingress-controller-conformance [command]

Available Commands:
  apply       Apply Ingress conformance resources to the current cluster
  context     Print the current Kubernetes context, server version, and supported Ingress APIVersions
  help        Help about any command
  list        List all Ingress verifications
  verify      Run Ingress verifications for conformance

Flags:
  -h, --help   help for ingress-controller-conformance

Use "ingress-controller-conformance [command] --help" for more information about a command.
```

### ingress-conformance-echo

The `ingress-conformance-echo` binary is published as docker image of the same name. The purpose of this component is to handle backend-requests made through an Ingress interface and respond using data from the original request. This, in turn, allows to build assertions on the original HTTP request as it is relayed through the ingress-controller.

```
# You may chose to pass in a TEST_ID environment variable to assert which running process actually handles the downstream request.
$ TEST_ID=sample ./echoserver
Starting server, listening on port 3000
Reporting TestId 'sample'

# You can then send test HTTP requests on localhost:3000
$ curl localhost:3000
{"TestId":"sample","Path":"/","Host":"localhost:3000","Method":"GET","Proto":"HTTP/1.1","Headers":{"Accept":["*/*"],"User-Agent":["curl/7.54.0"]}}
```

---

## Building

### Local build

#### Build and push the `ingress-conformance-echo` docker image:
```console
$ docker build -f build/package/Dockerfile . -t agervais/ingress-conformance-echo:latest
$ docker push agervais/ingress-conformance-echo:latest
```

#### Build the `ingress-controller-conformance` CLI and the `echo-server`:
```console
$ make
go build -o echoserver tools/echoserver.go
go build -o ingress-controller-conformance .
```

### Release build

```console
$ export VERSION=0.0.1 # Or git tag
$ go build -ldflags "-X github.com/kubernetes-sigs/ingress-controller-conformance/cmd.VERSION=$VERSION"
```
