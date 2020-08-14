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

set -o errexit
set -o nounset
set -o pipefail

# Shutdown the tests gracefully then save the results
shutdown () {
    TEST_SUITE_PID=$(pgrep ingress-controller-conformance)
    echo "sending TERM to ${TEST_SUITE_PID}"
    kill -s TERM "${TEST_SUITE_PID}"

    # Kind of a hack to wait for this pid to finish.
    # Since it's not a child of this shell we cannot use wait.
    tail --pid "${TEST_SUITE_PID}" -f /dev/null
}

# We get the TERM from kubernetes and handle it gracefully
trap shutdown TERM

mkdir -p "${RESULTS_DIR}"

set -x
/ingress-controller-conformance \
    --format=cucumber \
    --ingress-class="${INGRESS_CLASS}" \
    --output-directory="${RESULTS_DIR}" \
    --wait-time-for-ingress-status="${WAIT_FOR_STATUS_TIMEOUT}" \
    --test.timeout="${TEST_TIMEOUT}"
ret=$?
set -x
exit ${ret}
