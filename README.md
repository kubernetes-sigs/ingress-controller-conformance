# ingress-conformance

## Running

At the moment, `ingress-conformance` does not apply modifications to the running resources in the target Kubernetes cluster.

You must manually install and setup you environment targeted by the `kubectl config current-context`:
1. Apply the backing ingress-controller implementation. Samples are available under `examples/`.
1. Apply all, or a subset, of the ingress and service resources found under `deployments`. Files under the deployments folder match implemented check names.

### Context

```
./ingress-conformance context
Targetting Kubernetes cluster under active context 'docker-desktop'
The target Kubernetes cluster is running verion v1.14.6
```

### List

List, in a human-readable form, all Ingress verifications
```
./ingress-conformance list
- Ingress with host rule should send traffic to the correct backend service (host-rules)
- Ingress with path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request path (path-rules-foo)
- Ingress with path rule without a trailing slash should send traffic to the correct backend service, and preserve the original request including sub-paths (path-rules-foo-trailing)
- Ingress with path rule with a trailing slash should send traffic to the correct backend service, and preserve the original request path (path-rules-bar)
- Ingress with path rule with a trailing slash should send traffic to the correct backend service, and preserve the original request including sub-paths and double '/' (path-rules-bar-subpath)
- Ingress with no rules should send traffic to the correct backend service (single-service)
```

### Verify

Execute a series of assertions on the deployed ingress-controller.
```
./ingress-conformance verify
Running 'all' verifications...
Running 'host-rules' verifications...
Running 'path-rules' verifications...
Running 'path-rules-foo' verifications...
        1) Assertion failed: Expected the request path would be '/foo' but was '/'
  Check failed: path-rules-foo
Running 'path-rules-foo-trailing' verifications...
        1) Assertion failed: Expected the request path would be '/foo/' but was '//'
  Check failed: path-rules-foo-trailing
Running 'path-rules-bar' verifications...
        1) Assertion failed: Expected the request path would be '/bar/' but was '/'
  Check failed: path-rules-bar
Running 'path-rules-bar-subpath' verifications...
        1) Assertion failed: Expected the request path would be '/bar//bershop' but was '//bershop'
  Check failed: path-rules-bar-subpath
Running 'single-service' verifications...
--- Verification completed ---
3 checks passed! 4 failures!
in 1.777148914s
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
  list        List all Ingress verifications
  verify      Run Ingress verifications for conformance

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