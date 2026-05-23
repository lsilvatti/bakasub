#!/usr/bin/env sh

set -eu

release_env=${RELEASE_ENV:-release.env}

compose() {
  docker compose --env-file "$release_env" "$@"
}

compose up -d

frontend_port=$(compose port frontend 80 | awk -F: 'END { print $NF }')

if [ -z "$frontend_port" ]; then
  echo "BakaSub started, but the published frontend port could not be resolved." >&2
  exit 1
fi

echo "BakaSub is available at http://127.0.0.1:${frontend_port}"