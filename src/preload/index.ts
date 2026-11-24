import { contextBridge, ipcRenderer, webUtils, type IpcRendererEvent } from 'electron'
import {
  EVENT_CHANNELS,
  INVOKE_CHANNELS,
  type BakaSubApi,
  type BatchMuxPlanRequest,
  type BatchPlanRequest,
  type EventChannel,
  type InvokeChannel,
  type IpcInvokeContract,
  type IpcResult
} from '../shared/ipc-contract'

const invokeAllowlist: ReadonlySet<string> = new Set(INVOKE_CHANNELS)
const eventAllowlist: ReadonlySet<string> = new Set(EVENT_CHANNELS)

function invoke<C extends InvokeChannel>(
  channel: C,
  ...args: IpcInvokeContract[C]['req'] extends void ? [] : [IpcInvokeContract[C]['req']]
): Promise<IpcResult<IpcInvokeContract[C]['res']>> {
  if (!invokeAllowlist.has(channel)) {
    return Promise.reject(new Error(`IPC channel not allowed: ${channel}`))
  }
  return ipcRenderer.invoke(channel, ...args)
}

// Nunca expor `ipcRenderer` ou o `event` cru ao renderer.
const api: BakaSubApi = {
  app: {
    getInfo: () => invoke('app:get-info'),
    checkUpdate: () => invoke('app:check-update'),
    getUpdateState: () => invoke('app:get-update-state'),
    downloadUpdate: () => invoke('app:download-update'),
    cancelUpdateDownload: () => invoke('app:cancel-update-download'),
    installUpdate: (opts) => invoke('app:install-update', opts ?? {}),
    openExternal: (url) => invoke('app:open-external', { url }),
    onboardingFinished: () => invoke('app:onboarding-finished')
  },
  sidecars: {
    check: () => invoke('sidecars:check'),
    install: (tools) => invoke('sidecars:install', { tools }),
    cancel: () => invoke('sidecars:cancel')
  },
  memory: {
    stats: () => invoke('memory:stats'),
    prune: (olderThanDays) => invoke('memory:prune', { olderThanDays })
  },
  settings: {
    setApiKey: (apiKey) => invoke('settings:set-api-key', { apiKey }),
    clearApiKey: () => invoke('settings:clear-api-key'),
    apiKeyStatus: () => invoke('settings:api-key-status'),
    setTmdbKey: (key) => invoke('settings:set-tmdb-key', { key }),
    clearTmdbKey: () => invoke('settings:clear-tmdb-key'),
    tmdbKeyStatus: () => invoke('settings:tmdb-key-status'),
    setOpenSubtitlesKey: (key) => invoke('settings:set-opensubtitles-key', { key }),
    clearOpenSubtitlesKey: () => invoke('settings:clear-opensubtitles-key'),
    openSubtitlesStatus: () => invoke('settings:opensubtitles-status'),
    openSubtitlesLogin: (username, password) =>
      invoke('settings:opensubtitles-login', { username, password }),
    openSubtitlesLogout: () => invoke('settings:opensubtitles-logout'),
    get: () => invoke('settings:get'),
    update: (patch) => invoke('settings:update', patch),
    reset: (keys) => invoke('settings:reset', { keys })
  },
  logs: {
    tail: (query) => invoke('logs:tail', query ?? {}),
    watch: (on) => invoke('logs:watch', { on }),
    openFolder: () => invoke('logs:open-folder'),
    export: () => invoke('logs:export'),
    report: (report) => invoke('logs:report', report)
  },
  presets: {
    list: () => invoke('presets:list'),
    create: (preset) => invoke('presets:create', preset),
    update: (id, preset) => invoke('presets:update', { id, preset }),
    duplicate: (id, name) => invoke('presets:duplicate', { id, name }),
    delete: (id) => invoke('presets:delete', { id }),
    export: (ids) => invoke('presets:export', { ids }),
    import: (onConflict) => invoke('presets:import', { onConflict })
  },
  models: {
    list: (force) => invoke('models:list', { force }),
    freeQuota: () => invoke('models:free-quota')
  },
  media: {
    probe: (path) => invoke('media:probe', { path }),
    identify: (name) => invoke('media:identify', { name }),
    search: (query) => invoke('media:search', { query }),
    lookup: (ref, episode) => invoke('media:lookup', { ref, ...episode })
  },
  dialog: {
    openFiles: (kind, multiple) => invoke('dialog:open-files', { kind, multiple }),
    chooseOutputDir: () => invoke('dialog:choose-output-dir'),
    openFolder: () => invoke('dialog:open-folder')
  },
  batches: {
    plan: ((request: BatchPlanRequest | BatchMuxPlanRequest) =>
      invoke('batches:plan', request)) as BakaSubApi['batches']['plan'],
    create: (request) => invoke('batches:create', request),
    list: () => invoke('batches:list'),
    get: (id) => invoke('batches:get', { id }),
    jobs: (id) => invoke('batches:jobs', { id }),
    cancel: (id) => invoke('batches:cancel', { id }),
    remove: (id) => invoke('batches:remove', { id }),
    confirm: (id) => invoke('batches:confirm', { id })
  },
  jobs: {
    create: (request) => invoke('jobs:create', request),
    list: () => invoke('jobs:list'),
    get: (id) => invoke('jobs:get', { id }),
    cancel: (id) => invoke('jobs:cancel', { id }),
    retry: (id, memory) => invoke('jobs:retry', memory ? { id, memory } : { id }),
    remove: (id) => invoke('jobs:remove', { id }),
    clearCompleted: () => invoke('jobs:clear-completed'),
    confirm: (id) => invoke('jobs:confirm', { id }),
    setMemory: (id, memory) => invoke('jobs:set-memory', { id, memory }),
    dismissConsistency: (id) => invoke('jobs:dismiss-consistency', { id }),
    preflight: (id) => invoke('jobs:preflight', { id })
  },
  subtitles: {
    search: (request) => invoke('subtitles:search', request),
    download: (fileId) => invoke('subtitles:download', { fileId })
  },
  files: {
    inspect: (path) => invoke('files:inspect', { path }),
    pathFor: (file) => webUtils.getPathForFile(file)
  },
  events: {
    on(channel: EventChannel, listener) {
      if (!eventAllowlist.has(channel)) throw new Error(`IPC event not allowed: ${channel}`)
      // eslint-disable-next-line @typescript-eslint/no-explicit-any
      const wrapped = (_e: IpcRendererEvent, payload: any): void => listener(payload)
      ipcRenderer.on(channel, wrapped)
      return () => ipcRenderer.removeListener(channel, wrapped)
    }
  }
}

if (!process.contextIsolated) {
  throw new Error('contextIsolation must be enabled')
}
contextBridge.exposeInMainWorld('api', api)
