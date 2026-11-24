import { APP_LANGUAGES, type AppLanguage } from '@shared/ipc-contract'
import { resolveAppLanguage } from '@renderer/i18n/language'

/**
 * Guias embutidos no app (o repositório é privado, então não dá para linkar o GitHub).
 * Um arquivo por guia e por idioma: `guides/<lang>/<id>.md`. Inglês é a fonte e o fallback.
 */
const files = import.meta.glob<string>('./*/*.md', {
  query: '?raw',
  import: 'default',
  eager: true
})

export const GUIDE_IDS = [
  'choosing-models',
  'openrouter-key',
  'opensubtitles-key',
  'openrouter-limits',
  'tmdb-key'
] as const
export type GuideId = (typeof GUIDE_IDS)[number]

/** Guides exist in every app language (the test enforces it). */
export const GUIDE_LANGS = APP_LANGUAGES
export type GuideLang = AppLanguage
const FALLBACK: GuideLang = 'en'

const path = (lang: GuideLang, id: GuideId): string => `./${lang}/${id}.md`

/** Idioma do app -> idioma do guia (mesma regra da interface). */
export const resolveGuideLang = resolveAppLanguage

export interface Guide {
  readonly id: GuideId
  /** Idioma realmente servido (pode ser o fallback). */
  readonly lang: GuideLang
  readonly markdown: string
}

export function getGuide(id: GuideId, appLang: string | undefined): Guide {
  const lang = resolveGuideLang(appLang)
  const markdown = files[path(lang, id)]
  if (markdown !== undefined) return { id, lang, markdown }
  return { id, lang: FALLBACK, markdown: files[path(FALLBACK, id)] ?? '' }
}

/**
 * Link entre guias (`./openrouter-limits.md`, como no repositório) -> id, para o visualizador
 * navegar dentro do app. Qualquer outro link é externo (`app:open-external`) ou ignorado.
 */
export function guideIdFromHref(href: string): GuideId | null {
  const m = /^(?:\.\/)?([a-z0-9-]+)\.md(?:#.*)?$/.exec(href)
  return m && (GUIDE_IDS as readonly string[]).includes(m[1]) ? (m[1] as GuideId) : null
}

/** Para testes: arquivos `.md` presentes, por caminho relativo. */
export const guideFiles = (): readonly string[] => Object.keys(files)
