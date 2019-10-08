# ingress-conformance

## Running

```
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

### Local build

```console
go build
```

### Release build

```console
export VERSION=0.0.1 # Or git tag
go build -ldflags "-X github.com/datawire/ingress-conformance/cmd.VERSION=$VERSION"
```