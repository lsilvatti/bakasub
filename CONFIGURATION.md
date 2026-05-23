# Configuration Inventory

This file separates user configuration from infrastructure wiring in the Bakasub product.

## User Configuration

These values are persisted in PostgreSQL table `user_config` and must be edited through the Settings UI or `/api/v1/config`.

- `default_model`
- `default_preset`
- `default_language`
- `library_root_path`
- `remove_sdh_default`
- `video_timeout_minutes`
- `log_retention_days`
- `openrouter_api_key`
- `tmdb_access_token`
- `tmdb_metadata_enabled`
- `concurrent_translations`
- `max_retries`
- `base_retry_delay`
- favorite models in `favorite_models`

End users should not manage these values through `.env` files.

## Infrastructure Configuration

These values still exist outside the database because they describe deployment/runtime wiring rather than user preferences.

### Backend repository

- `DATABASE_URL`: PostgreSQL connection string for local development, CI, or deployment.

### Frontend repository

- `VITE_API_URL`: optional development/build override when the backend is not available through the default local proxy.
- `window.BAKASUB_RUNTIME_CONFIG.apiUrl`: optional hosting-time override injected before the app boots.

### Product repository

- `BACKEND_IMAGE` and `FRONTEND_IMAGE` in `release.env`: immutable release manifest values used to pin container digests.
- `IMAGE_PULL_POLICY`: optional runtime override for local smoke tests or controlled deployments.
- Compose-level `DATABASE_URL`: internal container wiring between backend and PostgreSQL, not a user preference.

## Rule of Thumb

If a value changes how a user translates, browses media, or configures providers, it belongs in the database.

If a value changes how services connect, build, or deploy, it stays in infrastructure config.