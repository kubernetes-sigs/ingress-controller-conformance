# Copyright 2020 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

MKDIR_P := mkdir -p
RM_F := rm -rf

export GO111MODULE=on

PROGRAMS := \
	ingress-conformance-tests

DEPLOYMENT_YAML := \
	$(wildcard deployments/*.yaml)

TAG ?= 0.0

REGISTRY ?= local

build: $(PROGRAMS) ## Build the conformance tool

.PHONY: build-image
build-image: ## Build the ingress conformance image
	docker build -t $(REGISTRY)/ingress-controller-conformance:$(TAG) .

.PHONY: publish-image
publish-image:
	docker push $(REGISTRY)/ingress-controller-conformance:$(TAG)

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
