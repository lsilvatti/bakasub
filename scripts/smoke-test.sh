#!/usr/bin/env sh

set -eu

release_env=${RELEASE_ENV:-release.env}
succeeded=0
frontend_host_port=

compose() {
  docker compose --env-file "$release_env" "$@"
}

resolve_frontend_port() {
  compose port frontend 80 | awk -F: 'END { print $NF }'
}

cleanup() {
  if [ "$succeeded" -ne 1 ]; then
    compose logs --no-color || true
  fi
  compose down -v --remove-orphans || true
}

trap cleanup EXIT INT TERM

compose up -d
frontend_host_port=$(resolve_frontend_port)

if [ -z "$frontend_host_port" ]; then
  echo "Could not determine the published frontend port" >&2
  exit 1
fi

wait_for() {
  url=$1
  label=$2
  attempts=${3:-60}

  while [ "$attempts" -gt 0 ]; do
    if curl -fsS "$url" >/dev/null 2>&1; then
      return 0
    fi

    attempts=$((attempts - 1))
    sleep 2
  done

  echo "Timed out waiting for $label at $url" >&2
  return 1
}

wait_for "http://127.0.0.1:${frontend_host_port}/" "frontend"
wait_for "http://127.0.0.1:${frontend_host_port}/api/v1/health" "backend API through frontend proxy"

succeeded=1
echo "Smoke test passed"