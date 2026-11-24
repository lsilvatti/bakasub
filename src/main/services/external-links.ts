/**
 * The only sites the renderer may open in the browser (key guides and sign-up pages). Exact
 * hosts, HTTPS, default port, no credentials: a compromised renderer cannot use
 * `app:open-external` to launch arbitrary URLs or protocol handlers.
 */
const ALLOWED_HOSTS = new Set([
  'www.themoviedb.org',
  'themoviedb.org',
  'openrouter.ai',
  'www.opensubtitles.com',
  'opensubtitles.com',
  // Settings → About: project repository and support page.
  'github.com',
  'buymeacoffee.com',
  'www.buymeacoffee.com'
])

export function isAllowedExternalUrl(raw: string): boolean {
  let url: URL
  try {
    url = new URL(raw)
  } catch {
    return false
  }
  return (
    url.protocol === 'https:' &&
    url.port === '' &&
    url.username === '' &&
    url.password === '' &&
    ALLOWED_HOSTS.has(url.hostname)
  )
}
