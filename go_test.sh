#!/bin/bash
#
# Run the go tests

err() {
  echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&2
}

log() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&1
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
    go test ./...
}

main
