#!/usr/bin/env sh

set -eu

manifest_file=release.env
version_file=VERSION
version=
backend_image=
frontend_image=

usage() {
  cat <<'EOF'
Usage: ./scripts/update-release-manifest.sh \
  --version 0.1.0 \
  --backend-image ghcr.io/owner/bakasub-backend:sha-abcdef1 \
  --frontend-image ghcr.io/owner/bakasub-frontend:sha-1234567

The script resolves tag-based image references to digest-based references and
updates VERSION plus release.env.
EOF
}

strip_tag() {
  ref=$1
  without_digest=${ref%@*}
  last_segment=${without_digest##*/}

  case "$last_segment" in
    *:*) printf '%s\n' "${without_digest%:*}" ;;
    *) printf '%s\n' "$without_digest" ;;
  esac
}

resolve_image_ref() {
  ref=$1

  case "$ref" in
    *@sha256:*)
      printf '%s\n' "$ref"
      return 0
      ;;
  esac

  if ! command -v docker >/dev/null 2>&1; then
    echo "docker is required to resolve image digests" >&2
    exit 1
  fi

  output=$(docker buildx imagetools inspect "$ref")
  digest=$(printf '%s\n' "$output" | awk '$1 == "Digest:" { print $2; exit }')

  if [ -z "$digest" ]; then
    echo "failed to resolve digest for $ref" >&2
    exit 1
  fi

  printf '%s@%s\n' "$(strip_tag "$ref")" "$digest"
}

while [ "$#" -gt 0 ]; do
  case "$1" in
    --version)
      version=${2-}
      shift 2
      ;;
    --backend-image)
      backend_image=${2-}
      shift 2
      ;;
    --frontend-image)
      frontend_image=${2-}
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

if [ -z "$version" ] || [ -z "$backend_image" ] || [ -z "$frontend_image" ]; then
  usage >&2
  exit 1
fi

resolved_backend=$(resolve_image_ref "$backend_image")
resolved_frontend=$(resolve_image_ref "$frontend_image")

printf '%s\n' "$version" > "$version_file"

cat > "$manifest_file" <<EOF
BAKASUB_VERSION=$version

# Product release manifest.
# This file is generated and pinned to immutable image digests.
BACKEND_IMAGE=$resolved_backend
FRONTEND_IMAGE=$resolved_frontend
EOF

echo "Updated $version_file and $manifest_file"
echo "Backend:  $resolved_backend"
echo "Frontend: $resolved_frontend"