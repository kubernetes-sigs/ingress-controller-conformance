MKDIR_P := mkdir -p
RM_F := rm -rf

export GO111MODULE=on

PROGRAMS := \
	echoserver \
	ingress-controller-conformance \
	ingress-conformance-tests

DEPLOYMENT_YAML := \
	$(wildcard deployments/*.yaml)

build: $(PROGRAMS) ## Build the conformance tool

.PHONY: image
image:
	docker build -t ingress-controller-conformance .

.PHONY: echoserver
echoserver: check-go-version
	go build -o $@ tools/echoserver.go

.PHONY: ingress-controller-conformance
ingress-controller-conformance: check-go-version internal/pkg/assets/assets.go
	go build -o $@ .

internal/pkg/assets/assets.go: $(DEPLOYMENT_YAML)
	@$(MKDIR_P) $$(dirname $@)
	@./hack/go-bindata.sh -pkg assets -o $@ $^

.PHONY: ingress-conformance-tests
ingress-conformance-tests: check-go-version
	go test -c -o $@ conformance_test.go

.PHONY: clean
clean: ## Remove build artifacts
	$(RM_F) internal/pkg/assets/assets.go
	$(RM_F) $(PROGRAMS)

.PHONY: codegen
codegen: check-go-version ## Generate or update missing Go code defined in feature files
	@go run hack/codegen.go -update -conformance-path=test/conformance features

.PHONY: verify-codegen
verify-codegen: check-go-version ## Verify if generated Go code is in sync with feature files
	@go run hack/codegen.go -conformance-path=test/conformance features

.PHONY: help
help: ## Display this help
	@echo Targets:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z0-9._-]+:.*?## / {printf "  %-20s %s\n", $$1, $$2}' $(MAKEFILE_LIST) | sort

.PHONY: check-go-version
check-go-version:
	@hack/check-go-version.sh
