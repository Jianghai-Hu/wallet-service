#!/usr/bin/env bash

# Source the common.sh script
# shellcheck source=./common.sh
. "$(git rev-parse --show-toplevel || echo ".")/scripts/common.sh"

cd "$PROJECT_DIR" || exit 1

# Mandatory tools
#-------------------------------------------------------------------------------
echo_info "Download golang dependencies"
go get ./...
#-------------------------------------------------------------------------------

cd "$WORKING_DIR" || exit 1