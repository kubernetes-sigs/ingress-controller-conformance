# ingress-conformance

## Running

At the moment, `ingress-conformance` does not apply modifications to the running resources in the target Kubernetes cluster.

You must manually install and setup you environment targeted by the `kubectl config current-context`:
1. Apply the backing ingress-controller implementation. Samples are available under `examples/`.
1. Apply all, or a subset, of the ingress and service resources found under `deployments`.

### Context

```
./ingress-conformance context
Targetting Kubernetes cluster under active context 'docker-desktop'
The target Kubernetes cluster is running verion v1.14.6
```

### Verify

Execute a series of assertions on the deployed ingress-controller.
```
./ingress-conformance verify 
Running all verifications...
Running host-rules verifications...
        1) Assertion failed: Expected the responding host would be 'foo.bar.com' but was ''
  Check failed: host-rules
Running path-rules verifications...
  Check passed: path-rules
Running single-service verifications...
  Check passed: single-service
2 checks passed! 1 failures
```

### Help

```
./ingress-conformance help

Kubernetes ingress controller conformance test suite in Go.
  Complete documentation is available at https://github.com/datawire/ingress-conformance

Usage:
  ingress-conformance [command]

Available Commands:
  context     Print the current Kubernetes context and server version
  help        Help about any command
  verify      Run all Ingress verifications for conformance

Flags:
  -h, --help   help for ingress-conformance

Use "ingress-conformance [command] --help" for more information about a command.
```

## Building

### Local build

#### Build and push the `ingress-conformance-echo` docker image:
```console
docker build -f build/package/Dockerfile . -t agervais/ingress-conformance-echo:latest
docker push agervais/ingress-conformance-echo:latest
```

#### Build the `ingress-conformance` CLI and the `echo-server`:
```console
go build
go build -o echoserver tools/echoserver.go
```

### Release build

```console
export VERSION=0.0.1 # Or git tag
go build -ldflags "-X github.com/datawire/ingress-conformance/cmd.VERSION=$VERSION"
```