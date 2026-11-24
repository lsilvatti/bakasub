import { docIdFor, docsIndexFile, listDocFiles } from './paths'

/** Pages for people using the app; every other `docs/*.md` goes under "Development". */
const guideFiles = [docsIndexFile, 'reporting-issues.md']

/** Reading order of the developer docs; files not listed here are appended alphabetically. */
const developmentOrder = [
  'architecture.md',
  'development.md',
  'ipc.md',
  'parsing.md',
  'llm.md',
  'rate-limits.md',
  'jobs.md',
  'job-memory.md',
  'preflight.md',
  'presets.md',
  'media-context.md',
  'subtitle-sources.md',
  'sidecars.md',
  'database.md',
  'logging.md',
  'notifications.md',
  'renderer.md',
  'brand-tokens.md',
  'release.md',
  'site.md'
]

/** Starlight sidebar: two groups built from the files on disk, so a new doc can never be orphaned. */
export function docsSidebar() {
  const files = listDocFiles()
  const missing = [...guideFiles, ...developmentOrder].filter((name) => !files.includes(name))
  if (missing.length > 0) throw new Error(`[docs] sidebar lists files missing from docs/: ${missing.join(', ')}`)
  const rest = files.filter((name) => !guideFiles.includes(name) && !developmentOrder.includes(name))
  const item = (name: string) => ({ slug: docIdFor(name) })
  return [
    { label: 'Guide', items: guideFiles.map(item) },
    { label: 'Development', items: [...developmentOrder, ...rest].map(item) }
  ]
}
