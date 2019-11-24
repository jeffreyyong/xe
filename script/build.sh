#!/usr/bin/env bash

# Exit on first error
set -e

CUR_DIR="$(pwd)"

print_help() {
    echo "usage: build.sh fix|lint|test|help"
    echo "  fix     auto format Go source code and tidy go.sum for each sub module in this project"
    echo "  lint    lint each sub module in this project"
    echo "  test    run test in each sub module"
    echo "  help    print this message"
    echo ""
    echo "Make sure golangci-lint is installed, by running:"
    echo "    go get -u github.com/golangci/golangci-lint/cmd/golangci-lint"
    exit 1
}

list_dirs() {
    find . -name '*.go' -print0 | xargs -0 -n1 dirname | sort --unique | grep -v vendor | grep -v tmp
}

fix() {
    echo "Fixing imports ..."
    goimports -l -w $(list_dirs)
    echo "Tidying go.sum ..."
    for d in $(list_dirs); do
        echo "Tidying directory $d ..."
        cd "$d"
        go mod tidy
        cd "$CUR_DIR"
    done
}

lint() {
    for d in $(list_dirs); do
        echo "Linting directory $d ..."
        cd "$d"
        golangci-lint run
        cd "$CUR_DIR"
    done
}

local_run() {
    export GO111MODULE=on
    go run xe.go
}

test_all() {
    go test $(go list ./...) -count=1
}

cover() {
    go test -race -count=1 -coverprofile=coverage.out $(go list ./...)
    go tool cover -html=coverage.out -o coverage.html
    echo Now open the web browser - for example
    echo open coverage.html
}

case "$1" in
    help)
        print_help
    ;;
    local_run)
        local_run
    ;;
    fix)
        fix
    ;;
    lint)
        lint
    ;;
    test)
        test_all
    ;;
    cover)
        cover
    ;;
    *)
        print_help
    ;;
esac

exit 0
