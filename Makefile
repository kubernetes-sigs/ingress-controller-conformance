MKDIR_P := mkdir -p
RM_F := rm -rf

export GO111MODULE=on

PROGRAMS := \
	echoserver \
	ingress-controller-conformance

DEPLOYMENT_YAML := \
	$(wildcard deployments/*.yaml)

build: $(PROGRAMS) ## Build the conformance tool

.PHONY: echoserver
echoserver:
	go build -o $@ tools/echoserver.go

.PHONY: ingress-controller-conformance
ingress-controller-conformance: internal/pkg/assets/assets.go
	go build -o $@ .

internal/pkg/assets/assets.go: $(DEPLOYMENT_YAML)
	@$(MKDIR_P) $$(dirname $@)
	@./hack/go-bindata.sh -pkg assets -o $@ $^

.PHONY: clean
clean: ## Remove build artifacts
	$(RM_F) internal/pkg/assets/assets.go
	$(RM_F) $(PROGRAMS)

.PHONY: help
help: ## Display this help
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort
