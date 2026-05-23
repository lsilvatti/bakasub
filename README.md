# BakaSub

This is the product repository for BakaSub. End users should only need this folder or GitHub repository to install, update, and roll back the application stack.

For the full release runbook, see [RELEASING.md](RELEASING.md).

## Responsibilities

- Orchestrate PostgreSQL, backend, and frontend with Docker Compose.
- Pin which backend and frontend images belong to a given product release.
- Hold the public product version and release tags.

The backend and frontend repositories are developer-facing repositories. They are responsible for local development, tests, and publishing their own Docker images.

## Local Setup

1. Copy `.env.example` to `.env`.
2. Set `VIDEO_DIR` to the absolute path of your local media directory.
3. Review `release.env`. Public releases should commit exact backend and frontend digest pins there. This local scaffold still carries placeholder digests until the first real release is prepared.
4. Start the stack:

```bash
docker compose --env-file release.env --env-file .env up -d
```

The frontend is exposed at `http://localhost:3000` and the backend API at `http://localhost:8080/api/v1`.

## Versioning

- The public product version lives in this repository.
- `release.env` is the release manifest for pinned container images.
- Release tags such as `v0.1.0` should be created here, not in the backend or frontend repositories.
- Backend and frontend image tags are implementation details pinned by this repository.
- Only immutable digest references should be pinned here.

## Upgrade Flow

1. Resolve the published backend and frontend image refs into digest pins and update the release manifest:

```bash
sh scripts/update-release-manifest.sh \
	--version 0.1.0 \
	--backend-image ghcr.io/lsilvatti/bakasub-backend:sha-abcdef1 \
	--frontend-image ghcr.io/lsilvatti/bakasub-frontend:sha-1234567
```

If you already downloaded the publish artifacts from the backend and frontend workflows, you can skip copying the refs manually:

```bash
sh scripts/prepare-release-from-artifacts.sh \
	--version 0.1.0 \
	--backend-artifact /path/to/backend-image-release.env \
	--frontend-artifact /path/to/frontend-image-release.env
```

2. Review the updated `VERSION` and `release.env` files.
3. Run `./scripts/smoke-test.sh`.
4. Commit the change in this repository.
5. Create a release tag in this repository.
6. Pull the new commit or tag and run `docker compose --env-file release.env --env-file .env up -d` again.

## Release Automation

If you prefer not to update the release manifest locally, use the `Prepare Release` workflow in this repository. It resolves the provided backend and frontend image refs to digests, updates `VERSION` plus `release.env`, and opens a pull request.

The easiest way to get those image refs is from the backend and frontend image publish workflows:

- Each publish run writes the digest pin and the `sha-*` tag into the workflow summary.
- Each publish run uploads an artifact named `backend-image-release` or `frontend-image-release` with `IMAGE_REF_DIGEST`, `IMAGE_SHA_TAG`, and related values.
- `scripts/prepare-release-from-artifacts.sh` can consume those downloaded artifact files directly and update `VERSION` plus `release.env`.
- You can pass either the digest pin or the `sha-*` tag into `Prepare Release`; the workflow will normalize the manifest to digest pins.

## Smoke Test

To validate the pinned release images locally, run:

```bash
./scripts/smoke-test.sh
```

The script boots the stack with `release.env` plus `.env.ci` by default, waits for the backend and frontend to answer, and then tears the stack down.

To validate your local runtime settings instead of the CI defaults, run:

```bash
RUNTIME_ENV=.env ./scripts/smoke-test.sh
```