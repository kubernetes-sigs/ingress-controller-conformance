#!/usr/bin/env bash

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

saveResults() {
    cd "${RESULTS_DIR}" || exit
    tar -czf results.tar.gz ./*
    # mark the done file as a termination notice.
    echo -n "${RESULTS_DIR}/results.tar.gz" > "${RESULTS_DIR}/done"
}

trap saveResults TERM

set -x
/ingress-controller-conformance \
    --format=cucumber \
    --ingress-class="${INGRESS_CLASS}" \
    --output-directory="${RESULTS_DIR}" \
    --wait-time-for-ingress-status="${WAIT_FOR_STATUS_TIMEOUT}" \
    --test.timeout="${TEST_TIMEOUT}"
ret=$?
#set -x
saveResults
#exit ${ret}
exit 0
