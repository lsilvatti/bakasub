/** Códigos de erro de domínio (`AppError`), estáveis: a UI traduz por código, nunca pela mensagem. */
export const APP_ERROR_CODES = [
  'NETWORK',
  'INVALID_MANIFEST',
  'CHECKSUM_MISMATCH',
  'EXTRACT_FAILED',
  'PROBE_FAILED',
  'PROCESS_TIMEOUT',
  'UNSUPPORTED_PLATFORM',
  'SIDECAR_BUSY',
  'SIDECAR_MISSING',
  'UNKNOWN_TOOL',
  'CANCELLED',
  'MEMORY_BUSY',
  'LLM_AUTH',
  'LLM_CREDITS',
  'LLM_BAD_REQUEST',
  'LLM_UNAVAILABLE',
  'LLM_DAILY_LIMIT',
  'LLM_UNKNOWN_MODEL',
  'LLM_TOO_MANY_FAILURES',
  /** O raciocínio do modelo esgotou `max_tokens` mesmo com o orçamento ampliado. */
  'LLM_REASONING_BUDGET',
  'LLM_NO_ENDPOINT',
  'INVALID_PATH',
  'DISK_FULL',
  'PRESET_NOT_FOUND',
  'PRESET_IN_USE',
  'PRESET_NAME_TAKEN',
  'PRESET_READ_ONLY',
  'MODEL_NOT_SET',
  'INVALID_PRESET_FILE',
  'INVALID_INPUT',
  'JOB_NOT_FOUND',
  'JOB_INVALID_STATE',
  'BATCH_NOT_FOUND',
  'MEDIA_TOOL_FAILED',
  'INVALID_SUBTITLE',
  'SETTINGS_CONFLICT',
  'TARGET_LANG_NOT_SET',
  'TMDB_AUTH',
  'TMDB_KEY_MISSING',
  'OPENSUBTITLES_AUTH',
  'OPENSUBTITLES_KEY_MISSING',
  'OPENSUBTITLES_QUOTA',
  'URL_NOT_ALLOWED',
  'UPDATE_UNSUPPORTED',
  'UPDATE_NOT_READY',
  'UPDATE_JOBS_ACTIVE'
] as const
export type AppErrorCode = (typeof APP_ERROR_CODES)[number]

/** Códigos que só a camada IPC produz (`createHandler`). */
export const IPC_ERROR_CODES = ['FORBIDDEN', 'INTERNAL'] as const

/** Tudo que pode chegar em `IpcError.code`. */
export const ERROR_CODES = [...APP_ERROR_CODES, ...IPC_ERROR_CODES] as const
export type ErrorCode = (typeof ERROR_CODES)[number]
