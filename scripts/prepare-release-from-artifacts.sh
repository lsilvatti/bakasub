#!/usr/bin/env sh

set -eu

backend_artifact=backend-image-release.env
frontend_artifact=frontend-image-release.env
manifest_file=release.env
version_file=VERSION
version=

usage() {
  cat <<'EOF'
Usage: sh scripts/prepare-release-from-artifacts.sh \
  --version 0.1.0 \
  --backend-artifact path/to/backend-image-release.env \
  --frontend-artifact path/to/frontend-image-release.env

The script reads IMAGE_REF_DIGEST from the backend and frontend release artifact
files and updates VERSION plus release.env.
EOF
}

read_env_value() {
  file=$1
  key=$2

  if [ ! -f "$file" ]; then
    echo "artifact file not found: $file" >&2
    exit 1
  fi

  value=$(grep "^${key}=" "$file" | cut -d= -f2- | head -n 1 || true)

  if [ -z "$value" ]; then
    echo "$key not found in $file" >&2
    exit 1
  fi

  printf '%s\n' "$value"
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version)
      version=${2-}
      shift 2
      ;;
    --backend-artifact)
      backend_artifact=${2-}
      shift 2
      ;;
    --frontend-artifact)
      frontend_artifact=${2-}
      shift 2
      ;;
    --manifest-file)
      manifest_file=${2-}
      shift 2
      ;;
    --version-file)
      version_file=${2-}
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "unknown argument: $1" >&2
      usage >&2
      exit 1
      ;;
  esac
done

if [ -z "$version" ]; then
  usage >&2
  exit 1
fi

backend_image=$(read_env_value "$backend_artifact" IMAGE_REF_DIGEST)
frontend_image=$(read_env_value "$frontend_artifact" IMAGE_REF_DIGEST)

sh scripts/update-release-manifest.sh \
  --version "$version" \
  --backend-image "$backend_image" \
  --frontend-image "$frontend_image" \
  --manifest-file "$manifest_file" \
  --version-file "$version_file"