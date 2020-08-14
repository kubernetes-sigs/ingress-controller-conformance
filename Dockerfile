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

# Build
FROM golang:1.15 as builder

WORKDIR /go/src/sigs.k8s.io/ingress-controller-conformance
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY . /go/src/sigs.k8s.io/ingress-controller-conformance

# Build
RUN make ingress-controller-conformance

FROM k8s.gcr.io/debian-hyperkube-base-amd64:0.12.1

RUN clean-install bash procps

ENV RESULTS_DIR="/tmp/results"
ENV INGRESS_CLASS="conformance"
ENV WAIT_FOR_STATUS_TIMEOUT="5m"
ENV TEST_TIMEOUT="20m"

COPY --from=builder /go/src/sigs.k8s.io/ingress-controller-conformance/ingress-controller-conformance /

COPY features /features
COPY run_e2e.sh /

CMD [ "/run_e2e.sh" ]
