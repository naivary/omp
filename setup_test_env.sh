#!/bin/bash
#
# Setup the test environment

err() {
  echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&2
}

log() {
    echo "[$(date +'%Y-%m-%dT%H:%M:%S%z')]: $*" >&1
}

# Report where docker is installed on the local machine
# Arguments:
#   None
function is_docker_installed() {
    local code=$(docker version)
    if [ $? -eq 0]; then
        return 0
    fi
    return 1
}

function is_go_installed() {
    local code=$(go version)
    if [ $? -eq 0]; then
        return 0
    fi
    return 1
}

# is_env_running is reporting whether the test environment is already up and
# running
function is_env_running() {
    local code=$(docker inspect omp_postgresql)
    if [ $? -eq 0]; then
        return 0
    fi
    return 1
}

function main() {
    if [ ! is_docker_installed ]; then
        err "docker engine is not isntalled"
    fi
    if [ ! is_go_installed ]; then
        err "golang is not isntalled"
    fi

    if [ is_env_running ]; then
        docker compose down
    fi

    docker compose up -d

    while true; do
        log "Waiting for keycloak to become ready..."
        curl -k -s -o /dev/null https://localhost:9000/health/ready
        if [ $? -eq 0 ]; then
            break
        fi
        sleep 2
    done

    terraform apply -auto-approve
}

main
