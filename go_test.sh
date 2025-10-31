#!/bin/bash
#
# Run the go tests

go_test_verbose=0

err() {
  echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&2
}

log() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&1
}

usage() {
    cat <<EOF
Usage: $(basename "$0") [options]

Options:
  -v        Enable verbose mode when running Go tests.
            This passes the '-v' flag to 'go test', causing it to display
            the names and results of each test as they are run.

  -h, --help
            Show this help message and exit.

Examples:
  $(basename "$0") -v        Run tests in verbose mode
  $(basename "$0")           Run tests normally (quiet mode)
EOF
exit 0
}

function main() {
    log "Delete omp database..."
    docker exec omp-postgresql psql -U postgres -q -c 'DROP DATABASE omp;';

    log "Create omp database..."
    docker exec omp-postgresql psql -U postgres -q -c 'CREATE DATABASE omp;';

    log "Delete omp realm..."
    terraform destroy -auto-approve > /dev/null

    log "Recreate omp realm..."
    terraform apply -auto-approve > /dev/null

    log "Clean go test cache..."
    go clean -testcache

    log "Run tests..."

    if [[ "$go_test_verbose" -eq 1 ]]; then
        go test ./... -v
        return
    fi

    go test ./...
}

while getopts 'hv' flag; do
  case "${flag}" in
    v) go_test_verbose=1 ;;
    h) usage ;;
  esac
done

main
