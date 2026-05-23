#!/usr/bin/env sh

set -eu

release_env=${RELEASE_ENV:-release.env}
runtime_env=${RUNTIME_ENV:-.env.ci}
succeeded=0

read_env_value() {
  key=$1
  file=$2

  if [ ! -f "$file" ]; then
    return 1
  fi

  grep "^${key}=" "$file" | tail -n 1 | cut -d= -f2-
}

backend_host_port=${BACKEND_HOST_PORT:-$(read_env_value BACKEND_HOST_PORT "$runtime_env" || printf '%s' 8080)}
frontend_host_port=${FRONTEND_HOST_PORT:-$(read_env_value FRONTEND_HOST_PORT "$runtime_env" || printf '%s' 3000)}

compose() {
  docker compose --env-file "$release_env" --env-file "$runtime_env" "$@"
}

cleanup() {
  if [ "$succeeded" -ne 1 ]; then
    compose logs --no-color || true
  fi
  compose down -v --remove-orphans || true
}

trap cleanup EXIT INT TERM

compose up -d

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

wait_for "http://127.0.0.1:${backend_host_port}/api/v1/health" "backend API"
wait_for "http://127.0.0.1:${frontend_host_port}/" "frontend"

succeeded=1
echo "Smoke test passed"