/**
 * Every query key in one place, so invalidation and event patches can't drift from the queries.
 * Hierarchical: invalidating `qk.jobs.all` also covers the list and every detail.
 */
export const qk = {
  app: {
    info: () => ['app', 'info'] as const,
    update: () => ['app', 'update'] as const
  },
  sidecars: {
    report: () => ['sidecars', 'report'] as const
  },
  settings: {
    all: () => ['settings', 'values'] as const,
    apiKey: () => ['settings', 'api-key'] as const,
    tmdbKey: () => ['settings', 'tmdb-key'] as const,
    openSubtitles: () => ['settings', 'opensubtitles'] as const
  },
  models: {
    all: ['models'] as const,
    list: () => ['models', 'list'] as const,
    freeQuota: () => ['models', 'free-quota'] as const
  },
  presets: {
    list: () => ['presets', 'list'] as const
  },
  media: {
    probe: (path: string) => ['media', 'probe', path] as const,
    identify: (name: string) => ['media', 'identify', name] as const,
    search: (q: string) => ['media', 'search', q.trim().toLowerCase()] as const,
    lookup: (kind: string, id: number, season: number, episode: number) =>
      ['media', 'lookup', kind, id, season, episode] as const
  },
  subtitles: {
    search: (request: unknown) => ['subtitles', 'search', request] as const
  },
  memory: {
    stats: () => ['memory', 'stats'] as const
  },
  batches: {
    plan: (kind: string, source: unknown) => ['batches', 'plan', kind, source] as const,
    /** `jobs`: the statuses of its jobs, so the counts and totals refetch as they change. */
    detail: (id: string, jobs: string) => ['batches', 'detail', id, jobs] as const
  },
  files: {
    inspect: (path: string) => ['files', 'inspect', path] as const
  },
  jobs: {
    all: ['jobs'] as const,
    list: () => ['jobs', 'list'] as const,
    detail: (id: string) => ['jobs', 'detail', id] as const,
    preflight: (id: string) => ['jobs', 'preflight', id] as const
  }
}
