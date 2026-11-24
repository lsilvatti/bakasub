import type { AppInfoService } from '../services/app-info.service'
import type { AppUpdater, AutoCheckHandle } from '../services/app-updater'
import type { BatchService } from '../batch/batch-service'
import type { JobManager } from '../jobs/job-manager'
import type { FreeQuota } from '../llm/free-quota'
import type { ModelCatalog } from '../llm/model-catalog'
import type { MkvToolnix } from '../media/mkvtoolnix-service'
import type { MediaContextService } from '../media-context/media-context.service'
import type { TmdbClient } from '../media-context/tmdb-client'
import type { JobMemoryService } from '../memory/job-memory.service'
import type { PresetService } from '../presets/preset.service'
import type { PreferencesService } from '../settings/preferences.service'
import type { SecretVault } from '../settings/secret-vault'
import type { SidecarManager } from '../sidecars/sidecar-manager'
import { registerAppHandlers } from './handlers/app.handler'
import { registerBatchHandlers } from './handlers/batches.handler'
import { registerFileHandlers } from './handlers/files.handler'
import { registerJobHandlers } from './handlers/jobs.handler'
import { registerLogHandlers } from './handlers/logs.handler'
import { registerMediaHandlers } from './handlers/media.handler'
import { registerMediaContextHandlers } from './handlers/media-context.handler'
import { registerMemoryHandlers } from './handlers/memory.handler'
import { registerPresetHandlers } from './handlers/presets.handler'
import { registerSettingsHandlers } from './handlers/settings.handler'
import { registerSidecarHandlers } from './handlers/sidecar.handler'
import { registerSubtitleHandlers } from './handlers/subtitles.handler'

export interface IpcDependencies {
  appInfo: AppInfoService
  appUpdate: AppUpdater
  autoCheck: Pick<AutoCheckHandle, 'onboardingFinished'>
  sidecars: SidecarManager
  memory: JobMemoryService
  vault: SecretVault
  preferences: PreferencesService
  presets: PresetService
  jobs: JobManager
  mkv: MkvToolnix
  models: ModelCatalog
  batches: BatchService
  tmdb: TmdbClient
  mediaContext: MediaContextService
  freeQuota: FreeQuota
  validateApiKey: (key: string) => Promise<unknown>
  subtitleSources: Parameters<typeof registerSubtitleHandlers>[0]['subtitleSources']
}

/** Composition root dos handlers IPC. Um arquivo por domínio em ./handlers. */
export function registerIpc(deps: IpcDependencies): void {
  registerAppHandlers(deps)
  registerSidecarHandlers(deps)
  registerMemoryHandlers(deps)
  registerSettingsHandlers(deps)
  registerPresetHandlers(deps)
  registerJobHandlers(deps)
  registerMediaHandlers(deps)
  registerFileHandlers()
  registerBatchHandlers(deps)
  registerLogHandlers()
  registerMediaContextHandlers(deps)
  registerSubtitleHandlers(deps)
}
