/**
 * Contrato IPC único (fonte da verdade) compartilhado entre main, preload e renderer.
 * Não importar zod/electron aqui: o preload roda em sandbox e só pode `require('electron')`.
 */

export interface AppInfo {
  name: string
  /** Raw `app.getVersion()` (semver); format it for display with `isDev`. */
  version: string
  /** Unpackaged build (`electron-vite dev`): own userData, AppUserModelId and no auto-update. */
  isDev: boolean
  platform: NodeJS.Platform
  arch: string
  electron: string
  userDataPath: string
}

/**
 * Estados do auto-update (ver docs/release.md). `unsupported`: dev, `linux-unpacked`, AppImage
 * sem `APPIMAGE` ou plataforma sem updater; use `fallbackUrl` para a página da release.
 */
export type AppUpdateState =
  | 'idle'
  | 'checking'
  | 'available'
  | 'not-available'
  | 'downloading'
  | 'downloaded'
  | 'error'
  | 'unsupported'

export interface AppUpdateProgress {
  /** 0–100. */
  percent: number
  transferred: number
  total: number
  bytesPerSecond: number
}

/** Retrato completo do updater; também é o payload do evento push `app:update-state`. */
export interface AppUpdateStatus {
  state: AppUpdateState
  /** Versão em execução. */
  current: string
  /** Versão publicada mais recente (a partir de `available`). */
  latest?: string
  /** ISO 8601 da release. */
  publishedAt?: string
  /** Notas da release (HTML vindo do GitHub, não confiável: sanitize antes de renderizar). */
  notes?: string
  /** Bytes do download. */
  size?: number
  /** Só em `downloading`. */
  progress?: AppUpdateProgress
  /** Só em `error`. */
  error?: IpcError
  /** `error` causado por rede/HTTP: o app segue normalmente. */
  offline?: boolean
  /** Só em `unsupported`: página da release (passa em `isAllowedExternalUrl`). */
  fallbackUrl?: string
}

export interface AppUpdateInfo {
  current: string
  latest: string
  updateAvailable: boolean
  /** Data da release (ISO 8601); vazio quando desconhecida. */
  publishedAt: string
  /** true quando a checagem falhou por rede/HTTP. */
  offline: boolean
  notes?: string
  size?: number
  state?: AppUpdateState
  fallbackUrl?: string
  error?: IpcError
}

export type SidecarState = 'ok' | 'missing' | 'outdated' | 'unsupported'

export interface SidecarToolStatus {
  name: string
  requiredVersion: string
  installedVersion: string | null
  state: SidecarState
  /** Caminho absoluto do executável quando instalado e íntegro. */
  path: string | null
}

export interface SidecarReport {
  checkedAt: string
  manifestUpdatedAt: string
  offline: boolean
  installDir: string
  tools: SidecarToolStatus[]
  /** true quando todas as ferramentas estão `ok`. */
  ready: boolean
}

/** Unidade (artefato) que falhou; as demais seguem instaladas. */
export interface SidecarInstallFailure {
  tools: string[]
  code: string
  message: string
}

export type SidecarInstallResult =
  | { outcome: 'installed'; tools: string[]; failures?: SidecarInstallFailure[] }
  | { outcome: 'declined' }
  | { outcome: 'cancelled' }

export type SidecarPhase = 'downloading' | 'extracting' | 'probing' | 'done'

export interface SidecarProgress {
  /** Ferramentas servidas pelo mesmo artefato, ex.: "mkvextract+mkvmerge". */
  archive: string
  phase: SidecarPhase
  received: number
  total: number | null
}

export interface MemoryStats {
  rows: number
  /** Bytes em disco da tabela e índices do Job Memory. */
  bytes: number
  /** ISO 8601 do `last_used_at` mais antigo; `null` com o cache vazio. */
  oldestUsedAt: string | null
}

/** "Não usado há mais de N dias", ou limpar tudo. */
export type MemoryPruneScope = 1 | 3 | 7 | 30 | 'all'

export interface MemoryPruneResult {
  deleted: number
  /** Estado depois do prune (o tamanho em disco já reflete o `incremental_vacuum`). */
  stats: MemoryStats
}

/**
 * Configurações que o renderer pode ler e alterar (`settings:get`/`settings:update`).
 * Chaves internas (ex.: `models.cache`) nunca aparecem aqui.
 */
export interface PublicSettings {
  'translation.defaultModel': string | null
  'translation.defaultTargetLang': string | null
  'translation.defaultPresetId': string
  'translation.targetLangs': string[]
  'models.favorites': string[]
  'jobs.concurrency': number
  'preflight.autoConfirmBelowUsd': number | null
  'tracks.preferredSourceLangs': string[]
  'logs.level': 'info' | 'debug'
  /** Checagem silenciosa de atualização ~30 s após o boot. Nunca baixa nem instala sozinha. */
  'updates.autoCheck': boolean
  /** Notificações nativas do SO (só com a janela sem foco). Geral: desligado cala as categorias. */
  'notifications.enabled': boolean
  /** Job avulso concluído / com avisos / falhou / aguardando confirmação. */
  'notifications.job': boolean
  /** Batch concluído ou aguardando confirmação (uma notificação agregada por batch). */
  'notifications.batch': boolean
  /** Nova versão disponível. */
  'notifications.update': boolean
  /** `null` = idioma do sistema (`navigator.language`). */
  'ui.language': AppLanguage | null
  /** Último modo de saída escolhido ao traduzir um vídeo; a UI começa por ele. */
  'output.mode': OutputMode
  /** Última nomeação escolhida ao extrair legendas. */
  'output.extractNaming': ExtractNaming
  /** Avançado (Settings → Advanced). Ver docs/rate-limits.md. */
  'llm.rateLimit.freeRpm': number
  'llm.rateLimit.paidRpm': number | null
  'llm.rateLimit.freeTier': FreeTierSetting
  'llm.retry.maxThrottled': number
  'llm.retry.maxTransient': number
  'llm.retry.baseDelayMs': number
  'llm.retry.maxDelayMs': number
  'llm.retry.chunkRounds': number
  /** Lanes paralelas por arquivo (1 = sequencial). */
  'llm.concurrency.perJob': number
  /** Teto global de requisições em voo (todos os jobs), por faixa. */
  'llm.concurrency.paidMaxInFlight': number
  'llm.concurrency.freeMaxInFlight': number
}

/**
 * Faixas recomendadas das settings de concorrência. Fora delas a UI avisa (sem bloquear); os
 * limites rígidos ficam no schema do main.
 */
export const CONCURRENCY_RECOMMENDED = {
  'llm.concurrency.perJob': { min: 2, max: 4 },
  'llm.concurrency.paidMaxInFlight': { min: 4, max: 12 },
  'llm.concurrency.freeMaxInFlight': { min: 2, max: 3 }
} as const satisfies Partial<Record<keyof PublicSettings, { min: number; max: number }>>

/** Idiomas da interface (os guias seguem o mesmo idioma). */
export const APP_LANGUAGES = ['en', 'pt-BR'] as const
export type AppLanguage = (typeof APP_LANGUAGES)[number]

export type FreeTierSetting = 'auto' | 'under-10' | 'over-10'

/** Cota diária de modelos `:free` no OpenRouter (zera 00:00 UTC). */
export interface FreeQuotaStatus {
  /** `unknown`: `GET /key` não respondeu (sem key, offline); assume-se 50/dia. */
  tier: 'under-10' | 'over-10' | 'unknown'
  /** `override`: fixado em Settings (`llm.rateLimit.freeTier`). */
  source: 'auto' | 'override'
  limit: number
  used: number | null
  remaining: number | null
}

export type LogLevel = 'debug' | 'info' | 'warn' | 'error'

/** Uma linha do log (arquivo JSON por linha), já sem segredos. */
export interface LogEntry {
  /** ISO 8601 (linhas antigas, em texto puro, trazem o horário como estava). */
  time: string
  level: LogLevel
  /** `app`, `db`, `ipc`, `sidecar`, `llm`, `jobs`, `renderer`… */
  scope: string
  msg: string
  /** Presente em tudo que foi logado durante a execução de um job. */
  jobId?: string
  /** Erro (código + stack), stderr de ferramentas, corpo de resposta HTTP… */
  detail?: string
}

export interface LogTailQuery {
  /** 1..2000, padrão 500. */
  limit?: number
  minLevel?: LogLevel
  scopes?: string[]
  jobId?: string
}

export interface RendererErrorReport {
  level: 'warn' | 'error'
  message: string
  stack?: string
  /** `boundary` (React), `window` (onerror), `promise` (unhandledrejection). */
  source: 'boundary' | 'window' | 'promise'
}

/** Contexto de mídia (TMDB). `series` = TV no TMDB. */
export type MediaKind = 'movie' | 'series'

export interface MediaRef {
  provider: 'tmdb'
  id: number
  kind: MediaKind
}

/** URLs de arte são `bakasub-media://…`: o Main baixa e guarda em cache; o renderer nunca vai à rede. */
export interface MediaSearchResult extends MediaRef {
  title: string
  originalTitle?: string
  year?: number
  posterUrl?: string
}

export interface MediaDetails extends MediaSearchResult {
  genres: string[]
  /** Ids de gênero do TMDB (16 = Animation, 99 = Documentary): estáveis, ao contrário dos nomes. */
  genreIds?: number[]
  /** ISO 639-1 do TMDB (`ja`). */
  originalLanguage?: string
  /** ISO 3166-1 (`JP`); séries trazem `origin_country`, filmes `production_countries`. */
  originCountries?: string[]
  /** Minutos (filme) ou duração típica do episódio (série). */
  runtime?: number
  synopsis: string
  backdropUrl?: string
}

export interface EpisodeDetails {
  season: number
  episode: number
  title: string
  synopsis?: string
}

/** O que o nome do arquivo diz sobre a mídia (parser de release). */
export interface ReleaseNameInfo {
  title: string
  year?: number
  season?: number
  episode?: number
  /** Numeração absoluta de anime (`Show - 28`): o usuário converte em temporada/episódio. */
  absoluteEpisode?: number
}

/**
 * Resultado da identificação automática. `matched` só quando o título (e o ano, se o nome tiver)
 * batem; senão `not-found` e o usuário busca. `no-key`: sem chave do TMDB (recurso opcional).
 */
export interface MediaIdentification {
  release: ReleaseNameInfo
  status: 'matched' | 'not-found' | 'no-key' | 'offline'
  match: MediaDetails | null
  episode: EpisodeDetails | null
}

export interface MediaContextLookup {
  details: MediaDetails
  episode: EpisodeDetails | null
}

export interface TmdbKeyStatus {
  present: boolean
}

export interface OpenSubtitlesStatus {
  keyPresent: boolean
  username: string | null
}

export interface SubtitleSearchRequest {
  videoPath: string
  query?: string
  languages: string[]
  season?: number
  episode?: number
}

export interface SubtitleSearchResult {
  fileId: number
  release: string
  language: string
  downloads: number
  hashMatch: boolean
  hearingImpaired: boolean
  machineTranslated: boolean
  trusted: boolean
  fps?: number
  uploader?: string
}

export interface SubtitleSearchResponse {
  query: string
  season?: number
  episode?: number
  results: SubtitleSearchResult[]
}

export interface SubtitleDownload {
  path: string
  name: string
  format: SubtitleTextFormat
  cached: boolean
  remaining?: number
  resetAt?: string
}

export interface ApiKeyStatus {
  present: boolean
  /** `os`: protegida pelo cofre do sistema; `none`: sem keyring disponível (a UI deve avisar). */
  protection: 'os' | 'none'
}

export interface GlossaryEntry {
  source: string
  target?: string
  rule?: string
}

/** Rótulo livre ("Anime", "Sitcom"…); só informativo, não entra no prompt. */
export type MediaType = string

/**
 * Preset de tradução: só ESTILO (prompt, temperatura, glossário, tipo de mídia). Modelo e idioma
 * de destino são escolhidos por job (padrões nas configurações). Os embutidos (`builtin: true`,
 * ids `builtin:*`) são imutáveis; os personalizados são criados do zero ou copiados (`copiedFrom`).
 */
export interface Preset {
  id: string
  builtin: boolean
  name: string
  mediaType: MediaType
  temperature: number
  /** Orientação de estilo para o tipo de mídia (vai no System Prompt, antes das regras fixas). */
  prompt: string
  glossary: GlossaryEntry[]
  copiedFrom: string | null
  /** `null` nos embutidos. */
  createdAt: string | null
}

export type PresetInput = Pick<Preset, 'name' | 'mediaType' | 'temperature' | 'prompt' | 'glossary'>

/** Formatos de legenda em texto que o app lê, traduz e grava (a extensão do arquivo é o id). */
export type SubtitleTextFormat = 'srt' | 'ass' | 'ssa' | 'vtt'

export type JobKind = 'extract' | 'translate' | 'mux' | 'pipeline'
export type JobStatus =
  | 'queued'
  | 'running'
  /** Reservado para o preflight (inativo nesta versão). */
  | 'awaiting-confirmation'
  | 'completed'
  | 'completed-with-warnings'
  | 'failed'
  | 'cancelled'
  | 'interrupted'
export type JobStep = 'probe' | 'extract' | 'preflight' | 'translate' | 'validate' | 'write' | 'mux'

export interface JobMediaRequest {
  ref: MediaRef
  season?: number
  episode?: number
}

export interface MuxOptions {
  /** A faixa nova vira `default` e as legendas originais perdem o `default`. Padrão `true`. */
  setDefault?: boolean
  /** Mantém as legendas originais no MKV novo. Padrão `true`. */
  keepOriginalSubtitles?: boolean
}

/**
 * Saída do `pipeline`: `mux` (MKV novo, padrão), `subtitle` (só a legenda ao lado do vídeo, o MKV
 * não é tocado) ou `both`.
 */
export type OutputMode = 'mux' | 'subtitle' | 'both'

/**
 * Nome da legenda extraída: `player` = `<base>.<lang>.<ext>` (players carregam sozinhos; várias
 * faixas do mesmo idioma viram `<base>.<lang>.<n>.<ext>`), `track` = `<base>.<trackId>.<lang>.<ext>`.
 */
export type ExtractNaming = 'player' | 'track'

/** Opções de saída comuns às legendas gravadas (extract/translate/pipeline). */
export interface WriteOptions {
  /** Colisão de nome: `true` substitui o arquivo existente; padrão `false` (` (2)`, ` (3)`…). */
  overwrite?: boolean
}

/**
 * Job Memory num job de tradução: `use` (padrão) consulta e grava o cache; `skip` ("ignorar
 * traduções salvas") trata todas as linhas como misses, mas continua gravando: as traduções novas
 * sobrescrevem as antigas, então o cache é corrigido, não só contornado.
 */
export type JobMemoryMode = 'use' | 'skip'

/** Pedido do renderer. Caminhos vêm do diálogo/drag & drop e são revalidados no Main. */
export type JobCreateRequest =
  | {
      kind: 'extract'
      source: string
      trackIds: number[]
      outputDir?: string
      options?: WriteOptions & {
        /** Padrão: `player` com uma só faixa, `track` com várias. */
        naming?: ExtractNaming
      }
    }
  | {
      kind: 'translate'
      source: string
      /** Padrão: `translation.defaultPresetId`. */
      presetId?: string
      /** Id OpenRouter. Padrão: `translation.defaultModel`. */
      model?: string
      /** Notas/contexto do usuário (sinopse, personagens…), até 4000 caracteres; vai no `<media>`. */
      notes?: string
      /**
       * Contexto do TMDB escolhido na UI (só quando "Include" está marcado). O Main busca os
       * detalhes e os congela no job.
       */
      media?: JobMediaRequest
      /** Padrão: `translation.defaultTargetLang`. */
      targetLang?: string
      sourceLang?: string
      outputDir?: string
      /** Para em `awaiting-confirmation` com a estimativa de custo. Padrão `true`. */
      options?: WriteOptions
      confirmBeforeTranslate?: boolean
      /** Padrão `use`. Congelado no spec; o retry herda (ou troca, ver `jobs:retry`). */
      memory?: JobMemoryMode
    }
  | {
      kind: 'mux'
      source: string
      subtitles: { path: string; language: string; name?: string; default?: boolean }[]
      outputDir?: string
      options?: MuxOptions
    }
  | {
      kind: 'pipeline'
      source: string
      /** Exactly one of `trackId` and `subtitle` (a subtitle from outside the video). */
      trackId?: number
      subtitle?: { path: string; language?: string }
      /** Padrão: `translation.defaultPresetId`. */
      presetId?: string
      /** Id OpenRouter. Padrão: `translation.defaultModel`. */
      model?: string
      /** Notas/contexto do usuário (sinopse, personagens…), até 4000 caracteres; vai no `<media>`. */
      notes?: string
      /**
       * Contexto do TMDB escolhido na UI (só quando "Include" está marcado). O Main busca os
       * detalhes e os congela no job.
       */
      media?: JobMediaRequest
      /** Padrão: `translation.defaultTargetLang`. */
      targetLang?: string
      sourceLang?: string
      outputDir?: string
      options?: MuxOptions &
        WriteOptions & {
          /** Padrão `mux`. */
          outputMode?: OutputMode
          /** @deprecated Use `outputMode`: `true` = `both`. Ignorado se `outputMode` vier. */
          alsoWriteSubtitle?: boolean
          includeSourceSubtitle?: boolean
        }
      /** Para em `awaiting-confirmation` com a estimativa de custo. Padrão `true`. */
      confirmBeforeTranslate?: boolean
      /** Padrão `use`. Congelado no spec; o retry herda (ou troca, ver `jobs:retry`). */
      memory?: JobMemoryMode
    }

/** `[from, to)` em posições da linha do tempo dos traduzíveis (`timeIndex`). */
export type TimelineRange = readonly [from: number, to: number]

/**
 * Onde do episódio a tradução está. Faixas compactadas; posições vão de 0 a `total`.
 * `active` = chunks em voo; `lanes` = segmento de cada lane (cortes = fronteiras de paralelismo).
 */
export interface TranslationTimeline {
  total: number
  cached: TimelineRange[]
  translated: TimelineRange[]
  failed: TimelineRange[]
  active: TimelineRange[]
  lanes: TimelineRange[]
}

export interface JobProgress {
  step: JobStep | null
  /** 0..1 dentro da etapa, quando conhecido. */
  stepFraction?: number
  /** Etapa de tradução. */
  translation?: {
    done: number
    total: number
    cached: number
    /** Chunks concluídos, total e em voo (as lanes rodam em paralelo). */
    chunk?: { completed: number; count: number; inFlight: number }
    /** Mapa da linha do tempo para a barra de progresso. */
    timeline?: TranslationTimeline
  }
  cost?: number
  /**
   * Parado esperando o OpenRouter: `rate-limit` (teto por minuto), `cooldown` (429 recente na
   * faixa; vale para todos os jobs), `backoff` (retry após erro). Ausente = trabalhando.
   */
  waiting?: { reason: 'rate-limit' | 'cooldown' | 'backoff'; untilMs: number }
}

/**
 * Fala que ficou no idioma original. `excerpt`/`location` vêm do documento de origem (mesmo
 * formato de `CoverageEntry`); ausentes em jobs gravados antes de existirem.
 */
export interface FailedLine {
  cueId: number
  reason: 'UNTRANSLATED_VALIDATION_FAILED' | 'UNTRANSLATED_CONTENT_FILTERED'
  /** ≤ 120 caracteres, controles visíveis. */
  excerpt?: string
  location?: { line: number; endLine: number; startMs?: number; endMs?: number }
}

/** Uma fala onde o termo aparece: origem e tradução em texto corrido (≤ 120 caracteres). */
export interface TermOccurrence {
  cueId: number
  line: number
  startMs: number
  source: string
  translation: string
}

/** Uma forma de traduzir o termo e quantas falas a usam. */
export interface TermVariant {
  /** Como aparece na tradução (a grafia mais frequente). */
  text: string
  cues: number
  /** O termo foi mantido como no original. */
  kept: boolean
  /** Até 3 falas de exemplo. */
  examples: TermOccurrence[]
}

/** Termo recorrente da origem traduzido de mais de uma forma no mesmo arquivo. */
export interface TermVariation {
  /** Grafia mais frequente na origem (ex.: "Bro"). */
  term: string
  /** Falas traduzidas em que o termo aparece. */
  cues: number
  /** Ordenadas pela mais usada. Sempre ≥ 2. */
  variants: TermVariant[]
  /** Falas com o termo em que nenhuma das variantes foi encontrada (reformulação, omissão). */
  unmatched: number
  /** O glossário do preset já fixava este termo (`target`) e mesmo assim houve variação. */
  glossaryTarget?: string
}

/**
 * Passada local de consistência após a tradução (sem LLM): termos recorrentes da origem (nomes,
 * apelidos, termos capitalizados) e as formas diferentes como foram traduzidos.
 */
export interface TermConsistencyReport {
  /** Falas comparadas (traduzidas de fato; as que ficaram iguais à origem não contam). */
  analyzedCues: number
  /** Termos recorrentes examinados. */
  candidates: number
  /** Só os termos com variação, do mais frequente ao menos. */
  terms: TermVariation[]
  /** Preset do job: onde fixar a variante escolhida. */
  preset: { id: string; name: string }
}

export interface JobResult {
  outputs: string[]
  /** Legenda gravada (`translate`/`pipeline` com saída `subtitle`/`both`); também está em `outputs`. */
  subtitle?: string
  warnings: string[]
  usage?: {
    requests: number
    promptTokens: number
    completionTokens: number
    cachedTokens: number
    cost: number
    costIncomplete: boolean
  }
  /** `skipped`: o job ignorou as traduções salvas (`memory: 'skip'`); `hits` é 0. */
  cache?: { hits: number; misses: number; skipped?: boolean }
  failedLines?: FailedLine[]
  /** Termos traduzidos de formas diferentes (ausente em jobs anteriores a esta passada). */
  consistency?: TermConsistencyReport
  /** O usuário dispensou as sugestões de `consistency` neste job (`jobs:dismiss-consistency`). */
  consistencyDismissed?: boolean
  /** Estimativa do preflight (a cobertura completa vem de `jobs:preflight`). */
  preflight?: PreflightEstimate
}

/** Mídia congelada no job, para a fila (título, arte, episódio). Sem sinopse. */
export interface JobMediaSummary {
  title: string
  kind: MediaKind
  year?: number
  posterUrl?: string
  backdropUrl?: string
  episode?: { season: number; episode: number; title: string }
}

/**
 * How a `translate`/`pipeline` job was set up, frozen in its spec at creation: editing or deleting
 * the preset later does not change it. Fields missing = not recorded (spec persisted before them).
 */
export interface JobTranslationSummary {
  /** OpenRouter model id actually used. */
  model?: string
  /** Preset as it was named when the job was created. */
  preset?: { id: string; name: string }
  targetLang?: string
  /** Absent = detected from the file/track when the job ran (or left to the model). */
  sourceLang?: string
}

export interface Job {
  id: string
  kind: JobKind
  batchId: string | null
  status: JobStatus
  /** Etapa atual (running) ou onde parou (failed/cancelled/interrupted). */
  step: JobStep | null
  /** Nome do arquivo de origem (para listar sem expor o caminho completo à toa). */
  source: string
  progress: JobProgress
  /** Só com contexto de mídia incluído. */
  media?: JobMediaSummary
  /** Só `translate`/`pipeline`: modelo, preset e idiomas congelados no spec. */
  translation?: JobTranslationSummary
  /** Só `translate`/`pipeline`: a escolha congelada no spec (jobs antigos = `use`). */
  memory?: JobMemoryMode
  result: JobResult | null
  error: { code: string; message: string } | null
  attempt: number
  createdAt: string
  updatedAt: string
  startedAt: string | null
  finishedAt: string | null
}

export interface JobProgressEvent {
  id: string
  status: JobStatus
  progress: JobProgress
}

export interface SubtitleTrackInfo {
  id: number
  codecId: string
  /** `null` = não traduzível/extraível nesta versão (ver `unsupported`). */
  format: SubtitleTextFormat | null
  unsupported?: 'image-based' | 'no-parser'
  language?: string
  name?: string
  default: boolean
  forced: boolean
  eventCount?: number
}

export interface MediaProbe {
  subtitles: SubtitleTrackInfo[]
  attachments: number
}

export interface ModelSummary {
  id: string
  name: string
  contextLength: number | null
  maxCompletionTokens: number | null
  isModerated: boolean
  /** USD por token; `null` = preço dinâmico. */
  pricing: { prompt: number; completion: number } | null
  structuredOutputs: boolean
  responseFormat: boolean
  /** Família do tokenizer (GPT, Claude…): só GPT tem contagem local exata. */
  tokenizer?: string
}

export interface ModelList {
  models: ModelSummary[]
  fetchedAt: string
  offline: boolean
}

export type OpenDialogKind = 'media' | 'subtitle' | 'any'

export type BatchKind = 'extract' | 'translate' | 'pipeline' | 'mux'

export type BatchSource =
  { type: 'folder'; path: string; recursive?: boolean } | { type: 'files'; paths: string[] }

export interface BatchPlanRequest {
  kind: Exclude<BatchKind, 'mux'>
  source: BatchSource
}

/** Plano do batch `mux`: pasta ou arquivos soltos com MKVs e legendas misturados. */
export interface BatchMuxPlanRequest {
  kind: 'mux'
  source: BatchSource
}

export interface BatchPlanItem {
  path: string
  /** Relativo à pasta do batch (ou só o nome, para arquivos avulsos). */
  relPath: string
  /** Faixas do MKV (`extract`/`pipeline`). */
  tracks?: SubtitleTrackInfo[]
  /** Seleção automática: várias no `extract`, uma no `pipeline`, nenhuma no `translate`. */
  selected: number[]
  /** Arquivo sem o que processar (o usuário pode removê-lo do plano ou escolher outra faixa). */
  skipped?: 'no-subtitle-track' | 'no-supported-track' | 'probe-failed'
  error?: string
}

export interface BatchPlan {
  kind: BatchKind
  /** Pasta de origem (espelhamento) ou `null`. */
  root: string | null
  items: BatchPlanItem[]
  /** Caminhos recusados pelo validador (arquivos avulsos). */
  rejected: { path: string; reason: string }[]
  /** Pasta com mais arquivos que o limite. */
  truncated: boolean
}

export interface BatchMuxPlanFile {
  path: string
  /** Relativo à pasta do batch (ou só o nome, para arquivos avulsos). */
  relPath: string
}

export interface BatchMuxPlanSubtitle extends BatchMuxPlanFile {
  /** Do sufixo do nome (`.pt-BR.srt`); `und` sem sufixo. Editável no `create`. */
  language: string
}

export interface BatchMuxPlanItem {
  video: BatchMuxPlanFile
  /** 1..32, em ordem natural. `matchedBy`: nome-base igual (certo) ou temporada/episódio. */
  subtitles: (BatchMuxPlanSubtitle & { matchedBy: 'name' | 'episode' })[]
}

/**
 * Plano de `kind: 'mux'`. Pareamento determinístico (docs/jobs.md): ambiguidade nunca é
 * resolvida sozinha, vai para `unmatched*` para o usuário decidir.
 */
export interface BatchMuxPlan {
  kind: 'mux'
  root: string | null
  /** Só vídeos com ao menos uma legenda. */
  items: BatchMuxPlanItem[]
  unmatchedVideos: BatchMuxPlanFile[]
  unmatchedSubtitles: BatchMuxPlanSubtitle[]
  rejected: { path: string; reason: string }[]
  truncated: boolean
}

/**
 * Obra de um batch (`translate`/`pipeline`): a unidade de contexto. Os arquivos da obra
 * compartilham a mídia do TMDB e, quando definidos, preset, modelo, idioma e notas, que
 * substituem os padrões do batch. Sem `media`, os jobs da obra não levam contexto de mídia.
 */
export interface BatchGroup {
  /** Rótulo da obra na UI. */
  name?: string
  media?: MediaRef
  presetId?: string
  model?: string
  targetLang?: string
  notes?: string
}

export interface BatchCreateItem {
  path: string
  trackIds?: number[]
  /** Índice em `groups`. */
  group?: number
  /** Episódio da obra (série), quando a obra tem `media`. */
  season?: number
  episode?: number
}

/** Item de um batch `mux`: um MKV e 1..32 legendas externas (validadas e lidas no job). */
export interface BatchMuxCreateItem {
  video: string
  subtitles: { path: string; language: string; name?: string; default?: boolean }[]
}

export interface BatchCreateRequest {
  kind: BatchKind
  name?: string
  /** `plan.root`: revalidado; base do espelhamento quando há `outputDir`. */
  root: string | null
  /** `BatchMuxCreateItem` se e somente se `kind: 'mux'`. */
  items: (BatchCreateItem | BatchMuxCreateItem)[]
  /** Só `translate`/`pipeline` (`mux` com `groups` é `INVALID_INPUT`). */
  groups?: BatchGroup[]
  /** Padrões do batch (cada um cai no padrão das configurações quando omitido). */
  presetId?: string
  model?: string
  notes?: string
  targetLang?: string
  sourceLang?: string
  /** Sem `outputDir`: saída ao lado de cada arquivo. Com `outputDir` e `root`: estrutura espelhada. */
  outputDir?: string
  options?: MuxOptions &
    WriteOptions & {
      /** `pipeline`. Padrão `mux`. */
      outputMode?: OutputMode
      /** @deprecated Use `outputMode`. */
      alsoWriteSubtitle?: boolean
      /** `extract`. Padrão: `player` com uma só faixa por item, `track` com várias. */
      extractNaming?: ExtractNaming
    }
  confirmBeforeTranslate?: boolean
  /** `translate`/`pipeline`: vale para todos os jobs do batch. Padrão `use`. */
  memory?: JobMemoryMode
}

export interface Batch {
  id: string
  name: string
  kind: BatchKind
  root: string | null
  createdAt: string
  total: number
  counts: Partial<Record<JobStatus, number>>
  /** `active` enquanto algum job não é terminal. */
  state: 'active' | 'finished'
  /** Soma das estimativas dos jobs em `awaiting-confirmation` (`cost: null` se algum preço é dinâmico). */
  preflight?: { awaiting: number; toTranslate: number; cost: Range | null; approximate: boolean }
}

export interface Range {
  low: number
  high: number
}

/** Por que um bloco não vai ao LLM (`skipped` benigno, `untranslated` fica no original). */
export type CoverageReason =
  | 'SKIPPED_NO_LETTERS'
  | 'SKIPPED_EMPTY_BLOCK'
  | 'SKIPPED_DRAWING'
  | 'SKIPPED_FRAGMENT'
  | 'UNTRANSLATED_ORPHAN_TEXT'
  | 'UNTRANSLATED_BAD_TIMING'
  | 'UNTRANSLATED_MARKER_COLLISION'
  | 'UNTRANSLATED_KARAOKE'
  | 'UNTRANSLATED_BAD_EVENT'
  | 'UNTRANSLATED_VALIDATION_FAILED'
  | 'UNTRANSLATED_CONTENT_FILTERED'

/** Espelho serializável do `CoverageReport` do parser (o Main checa a compatibilidade no tipo). */
export interface CoverageEntry {
  ref: { kind: 'cue'; cueId: number } | { kind: 'block'; blockIndex: number }
  status: 'skipped' | 'untranslated'
  severity: 'info' | 'warn'
  reason: CoverageReason
  /** Para a UI localizar a mensagem; o Main nunca manda texto. */
  params: Readonly<Record<string, string | number>>
  location: { line: number; endLine: number; startMs?: number; endMs?: number }
  /** ≤ 120 caracteres, controles visíveis. */
  excerpt: string
  outcome: 'kept-original' | 'omitted'
  /** Código do diagnóstico do parser que causou o bloco, quando houver. */
  cause?: string
  actions: readonly ('edit-raw' | 'ignore')[]
}

export interface CoverageReport {
  totals: {
    cues: number
    blocks: number
    translatable: number
    skipped: number
    untranslated: number
  }
  byReason: Readonly<Partial<Record<CoverageReason, number>>>
  /** `untranslated` primeiro, depois `skipped`; ordem do documento. */
  entries: readonly CoverageEntry[]
  /** `entries` cortado (totais completos). */
  truncated: boolean
}

/** Severidade de um diagnóstico do parser. */
export type DiagnosticSeverity = 'info' | 'warn' | 'error'

/** Espelho serializável do `Diagnostic` do parser (o Main checa a compatibilidade no tipo). */
export interface ParserDiagnostic {
  /** Código estável (ex.: `SRT_OVERLAP`); a UI localiza a mensagem. */
  code: string
  severity: DiagnosticSeverity
  /** `auto`: a exportação reparada corrige sozinha; `report`: só informa. */
  fix: 'auto' | 'report'
  /** 1-based. */
  line: number
  /** 1-based. */
  column: number
  ref?: { kind: 'cue'; cueId: number } | { kind: 'block'; blockIndex: number }
  /** Para a UI localizar a mensagem; o Main nunca manda texto. */
  params?: Readonly<Record<string, string | number>>
}

/**
 * Confiança da detecção de encoding. `certain`: decisão determinística (BOM ou UTF-8 estrito);
 * os demais derivam da pontuação do detector (`score`): `high` ≥ 0.8, `medium` ≥ 0.5, senão `low`.
 */
export type EncodingConfidence = 'certain' | 'high' | 'medium' | 'low'

export interface DetectedEncoding {
  /** Ex.: `utf-8`, `utf-16le`, `windows-1252`. */
  name: string
  confidence: EncodingConfidence
  /** 0..1 do detector; `null` quando `confidence` é `certain`. */
  score: number | null
  /** O arquivo tinha BOM. */
  bom: boolean
  /** Estilo de quebra de linha do arquivo. */
  eol: 'lf' | 'crlf' | 'cr' | 'mixed'
}

export interface SubtitleInspection {
  file: { name: string; sizeBytes: number }
  format: SubtitleTextFormat
  encoding: DetectedEncoding
  diagnostics: {
    /** Ordem do arquivo. Cortado em `maxItems` (500); os totais seguem completos. */
    items: ParserDiagnostic[]
    totals: Record<DiagnosticSeverity, number>
    truncated: boolean
  }
  /** O que não seria traduzido (sem falhas do LLM: nada foi enviado). */
  coverage: CoverageReport
}

/** Estimativa do preflight (antes de qualquer chamada ao LLM). */
export interface PreflightEstimate {
  model: {
    id: string
    /** USD por token; `null` = preço dinâmico. */
    pricing: { prompt: number; completion: number } | null
    catalogFetchedAt: string | null
  }
  /** Tokenizer local (o200k) diferente do do modelo: faixa alargada em ±20%. */
  approximate: boolean
  /**
   * `toTranslate`: falas distintas enviadas ao LLM. `copies`: falas com o mesmo context hash de
   * uma delas (letreiro quadro a quadro), resolvidas sem requisição. `total` inclui as duas.
   */
  lines: { total: number; cached: number; toTranslate: number; copies: number; chunks: number }
  /** O job ignora as traduções salvas: `cached` é 0 e o custo é cheio. Ausente em jobs antigos. */
  memory?: JobMemoryMode
  tokens: { input: Range; output: Range }
  /** USD; `null` quando o preço é dinâmico. */
  cost: Range | null
  /**
   * Só para modelos free: requisições previstas (chunks + margem de reparo) contra o que resta da
   * cota diária. `enough: null` = cota desconhecida. Não bloqueia; força a confirmação.
   */
  freeQuota?: FreeQuotaStatus & { requests: number; enough: boolean | null }
  /** Decisão tomada pelo job. */
  decision: 'auto-confirmed' | 'needs-confirmation' | 'confirmed'
}

/**
 * O que entrou no System Prompt de um job de tradução (idêntico em todos os chunks). Gravado
 * pelo preflight, antes de qualquer chamada ao LLM.
 */
export interface PromptSummary {
  /** Texto exato enviado como mensagem `system`. */
  system: string
  /** Tokens do System Prompt (o200k; aproximado para modelos de outra família). */
  systemTokens: number
  glossaryEntries: number
  /** O preset tem orientação de estilo (`<style>`). */
  style: boolean
  /** Bloco `<media>`: o que foi incluído. */
  media: {
    tmdb: boolean
    episode: boolean
    notes: boolean
    /** O usuário pediu o contexto do TMDB, mas ele não pôde ser obtido: o motivo. */
    unavailable?: string
  }
}

export interface PresetImportResult {
  /** `null` se o diálogo foi cancelado. */
  path: string | null
  imported: Preset[]
  /** Entradas não importadas: inválidas (campo com problema) ou `name-taken` com `onConflict: 'skip'`. */
  skipped: { index: number; name?: string; reason: string }[]
}

/** `code` é um `ErrorCode` (ver `error-codes.ts`); tipado como string porque o Main é a fonte. */
export type IpcError = { code: string; message: string }
export type IpcResult<T> = { ok: true; data: T } | { ok: false; error: IpcError }

/** Canais request/response (ipcRenderer.invoke <-> ipcMain.handle). */
export interface IpcInvokeContract {
  'app:get-info': { req: void; res: AppInfo }
  /** Nunca rejeita por falha de rede: devolve `offline: true`. Sem download nem instalação. */
  'app:check-update': { req: void; res: AppUpdateInfo }
  /** Snapshot atual do updater, sem rede (o renderer lê ao montar; depois segue `app:update-state`). */
  'app:get-update-state': { req: void; res: AppUpdateStatus }
  /**
   * Inicia o download (só em `available`, senão `UPDATE_NOT_READY`; `UPDATE_UNSUPPORTED` fora do
   * AppImage/instalador). Retorna na hora com `downloading`; progresso e fim vêm por `app:update-state`.
   */
  'app:download-update': { req: void; res: AppUpdateStatus }
  'app:cancel-update-download': { req: void; res: { cancelled: boolean } }
  /**
   * Só em `downloaded` (senão `UPDATE_NOT_READY`). Com job executando e sem `confirmJobs`,
   * `UPDATE_JOBS_ACTIVE`. Com confirmação: aborta jobs, encerra sidecars/DB e reinicia instalando.
   */
  'app:install-update': { req: { confirmJobs?: boolean }; res: { installing: true } }
  /** Só HTTPS para hosts da allowlist do Main (TMDB, OpenRouter, GitHub, Buy Me a Coffee); senão `URL_NOT_ALLOWED`. */
  'app:open-external': { req: { url: string }; res: { opened: true } }
  /**
   * O wizard de boas-vindas terminou ("Concluir"): libera a checagem automática de update da
   * sessão (se `updates.autoCheck` e a regra de onboarding valerem). Idempotente.
   */
  'app:onboarding-finished': { req: void; res: { acknowledged: true } }
  'sidecars:check': { req: void; res: SidecarReport }
  /** O Main exibe diálogo nativo de consentimento; sem aceite, nada é baixado. */
  'sidecars:install': { req: { tools: string[] }; res: SidecarInstallResult }
  'sidecars:cancel': { req: void; res: { cancelled: boolean } }
  'memory:stats': { req: void; res: MemoryStats }
  /** Recusado com `MEMORY_BUSY` enquanto houver job ativo. */
  'memory:prune': { req: { olderThanDays: MemoryPruneScope }; res: MemoryPruneResult }
  /** Somente escrita: a chave nunca volta ao renderer. Conferida em `GET /key` antes (`LLM_AUTH`). */
  'settings:set-api-key': { req: { apiKey: string }; res: ApiKeyStatus }
  'settings:clear-api-key': { req: void; res: ApiKeyStatus }
  /** Validada no TMDB (`GET /configuration`) antes de ir para o vault: `TMDB_AUTH` se recusada. */
  'settings:set-tmdb-key': { req: { key: string }; res: TmdbKeyStatus }
  'settings:clear-tmdb-key': { req: void; res: TmdbKeyStatus }
  'settings:tmdb-key-status': { req: void; res: TmdbKeyStatus }
  'settings:set-opensubtitles-key': { req: { key: string }; res: OpenSubtitlesStatus }
  'settings:clear-opensubtitles-key': { req: void; res: OpenSubtitlesStatus }
  'settings:opensubtitles-status': { req: void; res: OpenSubtitlesStatus }
  'settings:opensubtitles-login': {
    req: { username: string; password: string }
    res: OpenSubtitlesStatus
  }
  'settings:opensubtitles-logout': { req: void; res: OpenSubtitlesStatus }
  /** `searchId` (UUID) makes the search cancellable through `subtitles:cancel-search`. */
  'subtitles:search': {
    req: SubtitleSearchRequest & { searchId?: string }
    res: SubtitleSearchResponse
  }
  /** Aborts the sender's search in flight; it then fails with `CANCELLED`. `false` if none. */
  'subtitles:cancel-search': { req: { searchId: string }; res: { cancelled: boolean } }
  'subtitles:download': { req: { fileId: number }; res: SubtitleDownload }
  /** Só o nome do arquivo é usado (nada é lido do disco). */
  'media:identify': { req: { name: string }; res: MediaIdentification }
  'media:search': { req: { query: string }; res: MediaSearchResult[] }
  'media:lookup': {
    req: { ref: MediaRef; season?: number; episode?: number }
    res: MediaContextLookup
  }
  'settings:api-key-status': { req: void; res: ApiKeyStatus }
  'settings:get': { req: void; res: PublicSettings }
  /** Patch parcial, validado por chave; aplica as regras de favoritos/idiomas (`SETTINGS_CONFLICT`). */
  'settings:update': { req: Partial<PublicSettings>; res: PublicSettings }
  /** Volta as chaves ao padrão de fábrica ("Restore defaults"), com as mesmas regras do update. */
  'settings:reset': { req: { keys: (keyof PublicSettings)[] }; res: PublicSettings }
  /** Últimas entradas (mais antigas primeiro), de todos os arquivos rotacionados. */
  'logs:tail': { req: LogTailQuery; res: LogEntry[] }
  /** Liga/desliga o envio de `logs:entry` (só enquanto a tela de logs está aberta). */
  'logs:watch': { req: { on: boolean }; res: { watching: boolean } }
  'logs:open-folder': { req: void; res: { opened: boolean } }
  /** Diálogo nativo para salvar; home vira `~` e segredos são removidos. `null` = cancelado. */
  'logs:export': { req: void; res: { path: string | null } }
  /** Erro do renderer gravado no log (escopo `renderer`), com limite de taxa. */
  'logs:report': { req: RendererErrorReport; res: { accepted: boolean } }
  /** Embutidos primeiro, depois os personalizados. */
  'presets:list': { req: void; res: Preset[] }
  'presets:create': { req: PresetInput; res: Preset }
  /** `PRESET_READ_ONLY` nos embutidos. */
  'presets:update': { req: { id: string; preset: PresetInput }; res: Preset }
  /** Cópia editável (de embutido ou personalizado); nome padrão "<nome> (copy)". */
  'presets:duplicate': { req: { id: string; name?: string }; res: Preset }
  /** `PRESET_READ_ONLY` nos embutidos; `PRESET_IN_USE` enquanto um job não terminal o usa. */
  'presets:delete': { req: { id: string }; res: { deleted: true } }
  /** Diálogo nativo de salvar; `path: null` se cancelado. Embutidos também podem ser exportados. */
  'presets:export': { req: { ids: string[] }; res: { path: string | null; count: number } }
  /** Diálogo nativo de abrir (.json); sempre cria presets personalizados novos. */
  'presets:import': {
    req: { onConflict?: 'rename' | 'skip' }
    res: PresetImportResult
  }
  'models:list': { req: { force?: boolean }; res: ModelList }
  /** Cota diária de modelos `:free` (`GET /key`, cache de 30 s; override em Settings). */
  'models:free-quota': { req: void; res: FreeQuotaStatus }
  /** Faixas de legenda de um MKV (caminho revalidado no Main). */
  'media:probe': { req: { path: string }; res: MediaProbe }
  /** Diálogo nativo; caminhos já validados. `[]` se cancelado. */
  'dialog:open-files': {
    req: { kind: OpenDialogKind; multiple?: boolean }
    res: { paths: string[] }
  }
  'dialog:choose-output-dir': { req: void; res: { path: string | null } }
  /** Responde já com o job enfileirado; a execução é assíncrona (`jobs:progress`). */
  'jobs:create': { req: JobCreateRequest; res: Job }
  'jobs:list': { req: void; res: Job[] }
  'jobs:get': { req: { id: string }; res: Job }
  'jobs:cancel': { req: { id: string }; res: Job }
  /**
   * `memory` ausente = a escolha original do job. Trocar refaz o preflight e a tradução (os
   * artefatos dessas etapas são descartados); só `translate`/`pipeline`.
   */
  'jobs:retry': { req: { id: string; memory?: JobMemoryMode }; res: Job }
  'jobs:remove': { req: { id: string }; res: { removed: true } }
  /** Remove todos os concluídos/cancelados (falhas e interrompidos ficam). */
  'jobs:clear-completed': { req: void; res: { removed: string[] } }
  /** `awaiting-confirmation` -> fila; retoma na tradução, sem reextrair nem reestimar. */
  'jobs:confirm': { req: { id: string }; res: Job }
  /**
   * Troca o `memory` de um job em `awaiting-confirmation`: refaz só o preflight (volta à fila e
   * para de novo para confirmação, mesmo abaixo do limite de confirmação automática).
   */
  'jobs:set-memory': { req: { id: string; memory: JobMemoryMode }; res: Job }
  /** Esconde as sugestões de consistência de termos deste job (job terminal; não toca no preset). */
  'jobs:dismiss-consistency': { req: { id: string }; res: Job }
  /**
   * Inspeciona uma legenda (.srt/.ass/.ssa/.vtt) sem custo de LLM e sem criar job: diagnósticos
   * do parser, cobertura e encoding detectado. Caminho revalidado no Main.
   */
  'files:inspect': { req: { path: string }; res: SubtitleInspection }
  /** Estimativa + cobertura completa (o que não será traduzido, onde e por quê). */
  'jobs:preflight': {
    req: { id: string }
    res: { estimate: PreflightEstimate; coverage: CoverageReport; prompt?: PromptSummary } | null
  }
  /** Pasta de entrada de um batch (legível). `null` se cancelado. */
  'dialog:open-folder': { req: void; res: { path: string | null } }
  /** Varre, faz o probe e pré-seleciona faixas; não cria nada. */
  'batches:plan': { req: BatchPlanRequest | BatchMuxPlanRequest; res: BatchPlan | BatchMuxPlan }
  /** Cria o batch e todos os jobs numa transação (um item inválido recusa tudo). */
  'batches:create': { req: BatchCreateRequest; res: Batch }
  'batches:list': { req: void; res: Batch[] }
  'batches:get': { req: { id: string }; res: Batch }
  'batches:jobs': { req: { id: string }; res: Job[] }
  'batches:cancel': { req: { id: string }; res: Batch }
  /** Só batches sem jobs ativos; remove também os jobs. */
  'batches:remove': { req: { id: string }; res: { removed: true } }
  /** Confirma todos os jobs do batch em `awaiting-confirmation`. */
  'batches:confirm': { req: { id: string }; res: Batch }
}

/** Para onde o clique numa notificação nativa leva a janela (evento `app:navigate`). */
export type NavigationTarget =
  { kind: 'job'; id: string } | { kind: 'batch'; id: string } | { kind: 'settings'; tab: 'about' }

/** Canais push (main -> renderer). Payload tipado por canal. */
export interface IpcEventContract {
  'sidecars:progress': SidecarProgress
  /** Cada transição do updater; `downloading` com throttle de 250 ms. */
  'app:update-state': AppUpdateStatus
  /** Throttle de 250 ms por job; transições de estado saem na hora. */
  'jobs:progress': JobProgressEvent
  /** Lotes a cada 250 ms, só enquanto `logs:watch` estiver ligado. */
  'logs:entry': LogEntry[]
  /** Clique numa notificação nativa: a janela já foi mostrada/focada; o renderer navega. */
  'app:navigate': NavigationTarget
}

export type InvokeChannel = keyof IpcInvokeContract
export type EventChannel = keyof IpcEventContract

/** Allowlists em runtime — o preload só expõe o que está listado aqui. */
export const INVOKE_CHANNELS = [
  'app:get-info',
  'app:check-update',
  'app:get-update-state',
  'app:download-update',
  'app:cancel-update-download',
  'app:install-update',
  'app:open-external',
  'app:onboarding-finished',
  'sidecars:check',
  'sidecars:install',
  'sidecars:cancel',
  'memory:stats',
  'memory:prune',
  'settings:set-api-key',
  'settings:clear-api-key',
  'settings:set-tmdb-key',
  'settings:clear-tmdb-key',
  'settings:tmdb-key-status',
  'settings:set-opensubtitles-key',
  'settings:clear-opensubtitles-key',
  'settings:opensubtitles-status',
  'settings:opensubtitles-login',
  'settings:opensubtitles-logout',
  'subtitles:search',
  'subtitles:cancel-search',
  'subtitles:download',
  'media:identify',
  'media:search',
  'media:lookup',
  'settings:api-key-status',
  'settings:get',
  'settings:update',
  'settings:reset',
  'logs:tail',
  'logs:watch',
  'logs:open-folder',
  'logs:export',
  'logs:report',
  'presets:list',
  'presets:create',
  'presets:update',
  'presets:duplicate',
  'presets:delete',
  'presets:export',
  'presets:import',
  'models:list',
  'models:free-quota',
  'media:probe',
  'dialog:open-files',
  'dialog:choose-output-dir',
  'jobs:create',
  'jobs:list',
  'jobs:get',
  'jobs:cancel',
  'jobs:retry',
  'jobs:set-memory',
  'jobs:remove',
  'jobs:clear-completed',
  'dialog:open-folder',
  'batches:plan',
  'batches:create',
  'batches:list',
  'batches:get',
  'batches:jobs',
  'batches:cancel',
  'batches:remove',
  'batches:confirm',
  'jobs:confirm',
  'jobs:dismiss-consistency',
  'jobs:preflight',
  'files:inspect'
] as const satisfies readonly InvokeChannel[]
export const EVENT_CHANNELS = [
  'sidecars:progress',
  'app:update-state',
  'jobs:progress',
  'logs:entry',
  'app:navigate'
] as const satisfies readonly EventChannel[]

/** Superfície tipada de `window.api`. */
export interface BakaSubApi {
  app: {
    getInfo(): Promise<IpcResult<AppInfo>>
    checkUpdate(): Promise<IpcResult<AppUpdateInfo>>
    getUpdateState(): Promise<IpcResult<AppUpdateStatus>>
    downloadUpdate(): Promise<IpcResult<AppUpdateStatus>>
    cancelUpdateDownload(): Promise<IpcResult<{ cancelled: boolean }>>
    installUpdate(opts?: { confirmJobs?: boolean }): Promise<IpcResult<{ installing: true }>>
    openExternal(url: string): Promise<IpcResult<{ opened: true }>>
    onboardingFinished(): Promise<IpcResult<{ acknowledged: true }>>
  }
  sidecars: {
    check(): Promise<IpcResult<SidecarReport>>
    install(tools: string[]): Promise<IpcResult<SidecarInstallResult>>
    cancel(): Promise<IpcResult<{ cancelled: boolean }>>
  }
  memory: {
    stats(): Promise<IpcResult<MemoryStats>>
    prune(olderThanDays: MemoryPruneScope): Promise<IpcResult<MemoryPruneResult>>
  }
  settings: {
    setApiKey(apiKey: string): Promise<IpcResult<ApiKeyStatus>>
    clearApiKey(): Promise<IpcResult<ApiKeyStatus>>
    apiKeyStatus(): Promise<IpcResult<ApiKeyStatus>>
    setTmdbKey(key: string): Promise<IpcResult<TmdbKeyStatus>>
    clearTmdbKey(): Promise<IpcResult<TmdbKeyStatus>>
    tmdbKeyStatus(): Promise<IpcResult<TmdbKeyStatus>>
    setOpenSubtitlesKey(key: string): Promise<IpcResult<OpenSubtitlesStatus>>
    clearOpenSubtitlesKey(): Promise<IpcResult<OpenSubtitlesStatus>>
    openSubtitlesStatus(): Promise<IpcResult<OpenSubtitlesStatus>>
    openSubtitlesLogin(username: string, password: string): Promise<IpcResult<OpenSubtitlesStatus>>
    openSubtitlesLogout(): Promise<IpcResult<OpenSubtitlesStatus>>
    get(): Promise<IpcResult<PublicSettings>>
    update(patch: Partial<PublicSettings>): Promise<IpcResult<PublicSettings>>
    reset(keys: (keyof PublicSettings)[]): Promise<IpcResult<PublicSettings>>
  }
  logs: {
    tail(query?: LogTailQuery): Promise<IpcResult<LogEntry[]>>
    watch(on: boolean): Promise<IpcResult<{ watching: boolean }>>
    openFolder(): Promise<IpcResult<{ opened: boolean }>>
    export(): Promise<IpcResult<{ path: string | null }>>
    report(report: RendererErrorReport): Promise<IpcResult<{ accepted: boolean }>>
  }
  presets: {
    list(): Promise<IpcResult<Preset[]>>
    create(preset: PresetInput): Promise<IpcResult<Preset>>
    update(id: string, preset: PresetInput): Promise<IpcResult<Preset>>
    duplicate(id: string, name?: string): Promise<IpcResult<Preset>>
    delete(id: string): Promise<IpcResult<{ deleted: true }>>
    export(ids: string[]): Promise<IpcResult<{ path: string | null; count: number }>>
    import(onConflict?: 'rename' | 'skip'): Promise<IpcResult<PresetImportResult>>
  }
  models: {
    list(force?: boolean): Promise<IpcResult<ModelList>>
    freeQuota(): Promise<IpcResult<FreeQuotaStatus>>
  }
  media: {
    probe(path: string): Promise<IpcResult<MediaProbe>>
    identify(name: string): Promise<IpcResult<MediaIdentification>>
    search(query: string): Promise<IpcResult<MediaSearchResult[]>>
    lookup(
      ref: MediaRef,
      episode?: { season: number; episode: number }
    ): Promise<IpcResult<MediaContextLookup>>
  }
  dialog: {
    openFiles(kind: OpenDialogKind, multiple?: boolean): Promise<IpcResult<{ paths: string[] }>>
    chooseOutputDir(): Promise<IpcResult<{ path: string | null }>>
    openFolder(): Promise<IpcResult<{ path: string | null }>>
  }
  batches: {
    plan(request: BatchMuxPlanRequest): Promise<IpcResult<BatchMuxPlan>>
    plan(request: BatchPlanRequest): Promise<IpcResult<BatchPlan>>
    create(request: BatchCreateRequest): Promise<IpcResult<Batch>>
    list(): Promise<IpcResult<Batch[]>>
    get(id: string): Promise<IpcResult<Batch>>
    jobs(id: string): Promise<IpcResult<Job[]>>
    cancel(id: string): Promise<IpcResult<Batch>>
    remove(id: string): Promise<IpcResult<{ removed: true }>>
    confirm(id: string): Promise<IpcResult<Batch>>
  }
  jobs: {
    create(request: JobCreateRequest): Promise<IpcResult<Job>>
    list(): Promise<IpcResult<Job[]>>
    get(id: string): Promise<IpcResult<Job>>
    cancel(id: string): Promise<IpcResult<Job>>
    retry(id: string, memory?: JobMemoryMode): Promise<IpcResult<Job>>
    remove(id: string): Promise<IpcResult<{ removed: true }>>
    clearCompleted(): Promise<IpcResult<{ removed: string[] }>>
    confirm(id: string): Promise<IpcResult<Job>>
    setMemory(id: string, memory: JobMemoryMode): Promise<IpcResult<Job>>
    dismissConsistency(id: string): Promise<IpcResult<Job>>
    preflight(
      id: string
    ): Promise<IpcResult<{ estimate: PreflightEstimate; coverage: CoverageReport } | null>>
  }
  subtitles: {
    search(
      request: SubtitleSearchRequest,
      searchId?: string
    ): Promise<IpcResult<SubtitleSearchResponse>>
    cancelSearch(searchId: string): Promise<IpcResult<{ cancelled: boolean }>>
    download(fileId: number): Promise<IpcResult<SubtitleDownload>>
  }
  files: {
    inspect(path: string): Promise<IpcResult<SubtitleInspection>>
    /** Drag & drop: caminho real de um `File` (`webUtils.getPathForFile`). Revalidado no Main. */
    pathFor(file: File): string
  }
  events: {
    on<C extends EventChannel>(
      channel: C,
      listener: (payload: IpcEventContract[C]) => void
    ): () => void
  }
}
