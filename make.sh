#!/bin/bash

if [[ $# -lt 1 ]]; then
  echo "Example: make.sh (build|clean|restart|start|status|stop)"
  exit 1
fi

set -ex

cd "$(dirname "$0")"

build() {
  docker build --tag lystor/pushgateway-admin .
}

clean() {
  docker-compose down --remove-orphans --volumes
}

push() {
  docker push lystor/pushgateway-admin
}

restart() {
  stop "$@"
  start "$@"
}

start() {
  docker-compose up --detach --remove-orphans "$@"
}

status() {
  docker-compose ps "$@"
}

stop() {
  docker-compose rm --force --stop --volumes "$@"
}

$*
