#!/usr/bin/env bash

# Source the common.sh script
# shellcheck source=./common.sh
. "$(git rev-parse --show-toplevel || echo ".")/scripts/common.sh"

cd "$PROJECT_DIR" || exit 1

echo_info "Remove all log files"
find . -path './log/*' -print | xargs rm -vrf
find . -path './log' -type d -empty -print | xargs rm -vrf

echo_info "Remove binary artifacts"
rm -vrf ./bin

cd "$WORKING_DIR" || exit 1
