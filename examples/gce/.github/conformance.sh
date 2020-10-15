#!/bin/bash

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

set -o errexit
set -o nounset
set -o pipefail

log() {
  echo "$(date -u +'%Y-%M-%dT%H:%M')" "$@"
}

if [ -z "${INGRESS_CLASS}" ]; then
  log "Environment variable INGRESS_CLASS must be set"
  exit 1
fi

if [ -z "${INGRESS_CONFORMANCE_IMAGE}" ]; then
  log "Environment variable INGRESS_CONFORMANCE_IMAGE must be set"
  exit 1
fi

log "Running... (can take some time)"

STARTTIME=`date +%s`;

sonobuoy run \
  --skip-preflight \
  --kube-conformance-image=${INGRESS_CONFORMANCE_IMAGE} \
  --plugin-env e2e.INGRESS_CLASS=${INGRESS_CLASS} \
  --plugin-env e2e.WAIT_FOR_STATUS_TIMEOUT=${WAIT_FOR_STATUS_TIMEOUT:-5m} \
  --plugin-env e2e.TEST_TIMEOUT=${TEST_TIMEOUT:-20m} \
  --wait

# retrieve the result file to local system
sonobuoy retrieve

mkdir -p /tmp/reports
tar zxpvf *_sonobuoy_*.tar.gz --wildcards "*-report.json"
mv plugins/e2e/results/global/* /tmp/reports

CURRTIME=`date +%s`;
CURRELAPSED=$(( CURRTIME - STARTTIME));

log "Conformance tests execution time was ${CURRELAPSED} seconds"
