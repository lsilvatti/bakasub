import type { ErrorCode } from '@shared/error-codes'

/** One message per error code. The type makes a missing code a compile error. */
const codes = {
  NETWORK: 'Could not reach the server. Check your internet connection and try again.',
  INVALID_MANIFEST: 'The download list for the tools is invalid. Try again later.',
  CHECKSUM_MISMATCH: 'The downloaded file did not pass the integrity check and was discarded.',
  EXTRACT_FAILED: 'The downloaded tool could not be unpacked.',
  PROBE_FAILED: 'The installed tool did not start. Reinstall it from Settings.',
  PROCESS_TIMEOUT: 'An external tool took too long and was stopped.',
  UNSUPPORTED_PLATFORM: 'This tool is not available for your system.',
  SIDECAR_BUSY: 'The tools are already being installed.',
  SIDECAR_MISSING: 'mkvtoolnix is not installed. Install it to work with video files.',
  UNKNOWN_TOOL: 'Unknown tool.',
  CANCELLED: 'Cancelled.',
  MEMORY_BUSY: 'The translation cache is busy. Try again when the running jobs finish.',
  LLM_AUTH: 'OpenRouter rejected the API key. Check it in Settings.',
  LLM_CREDITS: 'Your OpenRouter credits ran out. Add credits and retry the job.',
  LLM_BAD_REQUEST: 'The model rejected the request. Try another model.',
  LLM_UNAVAILABLE: 'The model is unavailable right now. Retry later or pick another model.',
  LLM_DAILY_LIMIT:
    'Daily limit for free models reached. It resets at 00:00 UTC; what was already translated is kept, so retrying later continues from there.',
  LLM_UNKNOWN_MODEL: 'This model is no longer available on OpenRouter. Pick another one.',
  LLM_TOO_MANY_FAILURES: 'Too many lines failed to translate. Try another model or preset.',
  LLM_NO_ENDPOINT:
    'No OpenRouter provider accepted this model for your account. Check the provider restrictions in your OpenRouter privacy settings (allowed/ignored providers, ZDR) or pick another model.',
  LLM_REASONING_BUDGET:
    'The model spent its whole output budget on reasoning and returned no translation. Pick another model.',
  INVALID_PATH: 'This file or folder cannot be used. Choose it again.',
  DISK_FULL: 'Not enough free disk space for the output.',
  PRESET_NOT_FOUND: 'This preset no longer exists.',
  PRESET_IN_USE: 'This preset is used by a queued or running job.',
  PRESET_NAME_TAKEN: 'A preset with this name already exists.',
  PRESET_READ_ONLY: 'Built-in presets cannot be changed. Duplicate it to edit a copy.',
  MODEL_NOT_SET: 'Choose a model first: star one in Models and make it the default.',
  INVALID_PRESET_FILE: 'This is not a valid BakaSub presets file.',
  INVALID_INPUT: 'Some of the information sent was invalid.',
  JOB_NOT_FOUND: 'This job no longer exists.',
  JOB_INVALID_STATE: 'This action is not available for the job right now.',
  BATCH_NOT_FOUND: 'This batch no longer exists.',
  MEDIA_TOOL_FAILED: 'mkvtoolnix could not process this file.',
  INVALID_SUBTITLE: 'This subtitle file could not be read.',
  SETTINGS_CONFLICT:
    'The default model and language must stay in your lists. Pick another default before removing it.',
  TARGET_LANG_NOT_SET:
    'Choose a target language first: add one in Languages and make it the default.',
  TMDB_AUTH: 'TMDB rejected the key. Check it in Settings.',
  TMDB_KEY_MISSING: 'Add a TMDB key in Settings to look up titles. It is optional.',
  OPENSUBTITLES_AUTH:
    'OpenSubtitles rejected the API key or the account. Check them in Settings (the account takes your username, not your email).',
  OPENSUBTITLES_KEY_MISSING: 'Add an OpenSubtitles API key in Settings to search for subtitles.',
  OPENSUBTITLES_QUOTA:
    'The OpenSubtitles download limit for today was reached. Try again tomorrow, or sign in to your account in Settings for a bigger limit.',
  URL_NOT_ALLOWED: 'This link cannot be opened from BakaSub.',
  UPDATE_UNSUPPORTED: 'Updates are only installed on an AppImage or the Windows installer.',
  UPDATE_NOT_READY: 'No update is ready to install.',
  UPDATE_JOBS_ACTIVE: 'Jobs are running. Confirm to cancel them and install the update.',
  FORBIDDEN: 'The request was blocked.',
  INTERNAL: 'Something went wrong inside BakaSub. Details are in the logs.'
} satisfies Record<ErrorCode, string>

export const errors = {
  ...codes,
  unknown: 'Something went wrong ({{code}}). Details are in the logs.'
}
