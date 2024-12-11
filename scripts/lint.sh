#!/usr/bin/env bash

# Source the common.sh script
# shellcheck source=./common.sh
. "$(git rev-parse --show-toplevel || echo ".")/scripts/common.sh"

PATH=./bin:$PATH # don't export this

if ! has golangci-lint || ! ./bin/golangci-lint --version | grep -q 1.61.0; then
    echo_info "Install golangci-lint for static code analysis (via curl)"
    # install into ./bin/
    # because different project might need different golang version,
    # and thus, need to use different linter version
    curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s v1.61.0
fi

CI_CONFIG_PATH=$CONFIG_PATH

if [[ -z "$CI_CONFIG_PATH" ]]; then
    CI_CONFIG_PATH="$PROJECT_DIR"/.golangci.yml
fi

lint_fast() {
    # Run only fast linter
    golangci-lint run --fast --timeout 120m0s "$@"
}

lint_all() {
    # Run all the linter on all packages
    golangci-lint run --timeout 120m0s
}

case "$1" in
all)
    lint_all
    ;;
*)
    shift
    lint_fast "$@"
    ;;
esac