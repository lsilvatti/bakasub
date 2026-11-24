/**
 * English strings, the reference dictionary. Other locales (`pt.ts`) must provide the same keys:
 * `Dictionary` is derived from this object so a missing key is a type error.
 */
export const en = {
  'meta.title': 'BakaSub: AI subtitle translator for MKV',
  'meta.description':
    'Extract, translate with AI and re-mux the subtitles of your .mkv videos, from a desktop app, on your own machine.',
  'nav.home': 'BakaSub home',
  'nav.docs': 'Docs',
  'nav.docsLabel': 'Docs',
  'nav.github': 'GitHub',
  'nav.githubLabel': 'BakaSub on GitHub',
  'nav.primary': 'Main',
  'nav.language': 'Language',
  'hero.kicker': 'AI subtitle translator',
  // U+2011 non-breaking hyphen: keeps "re‑mux" on one line at narrow widths.
  'hero.titleStrong': 'Extract, translate with AI and re\u2011mux',
  'hero.titleRest': 'the subtitles of your .mkv videos.',
  'hero.lead':
    'A desktop app that runs on your own machine. Bring your OpenRouter key; your files never leave your computer.',
  'download.windows': 'Download for Windows (.exe)',
  'download.linux': 'Download for Linux (.AppImage)',
  'download.others': 'Other versions',
  'download.mobile': 'BakaSub is a desktop app for Windows and Linux: open this page on your computer to download it.',
  'download.version': 'Version {version}',
  'hero.source': 'View source',
  'hero.platforms': 'Windows x64 · Linux x64 · Free and open source',
  'screenshots.title': 'A look inside',
  'screenshots.lead': 'From the first file to the final mux, every step stays in front of you.',
  'carousel.label': 'BakaSub screenshots',
  'carousel.controls': 'Choose a screenshot',
  'carousel.prev': 'Previous screenshot',
  'carousel.next': 'Next screenshot',
  // Placeholders: {index}, {total}.
  'carousel.slide': '{index} of {total}',
  'carousel.goTo': 'Show screenshot {index}',
  'carousel.placeholder': 'Screenshot coming soon.',
  'screenshots.wizard.caption': 'Translation wizard',
  'screenshots.wizard.alt': 'The translation wizard picking a subtitle track from an .mkv file.',
  'screenshots.cost.caption': 'Cost estimate',
  'screenshots.cost.alt': 'Preflight screen with token count, cache hit rate and the estimated cost.',
  'screenshots.queue.caption': 'Job queue',
  'screenshots.queue.alt': 'The job queue with translations running, waiting and finished.',
  'screenshots.tmdb.caption': 'TMDB context',
  'screenshots.tmdb.alt': 'Media context fetched from TMDB, with the backdrop and synopsis of the episode.',
  'screenshots.presets.caption': 'Presets & glossary',
  'screenshots.presets.alt': 'A translation preset with its glossary rules.',
  'footer.license': 'Free and open source software, licensed under',
  'footer.docs': 'Documentation',
  'footer.releases': 'All releases',
  'footer.coffee': 'Buy me a coffee',
  'footer.source': 'Source code',
  'features.kicker': 'What it does',
  'features.title': 'Built around the subtitle file, not just the text',
  'features.lead':
    'Translates .srt, .ass, .ssa and .vtt with any model from OpenRouter, and handles the MKV around them.',
  'features.cost.title': 'Shows the cost before you pay',
  'features.cost.body':
    'A token count and a price estimate from live OpenRouter pricing, plus a list of what will not be translated.',
  'features.memory.title': 'Remembers what it translated',
  'features.memory.body':
    'Translated lines are kept in a local cache, so reprocessing a file, or another one with the same lines, is nearly free.',
  'features.tmdb.title': 'Knows what you are watching',
  'features.tmdb.body':
    'It reads the release name and can pull the synopsis and episode details from TMDB to improve the translation.',
  'features.presets.title': 'Keeps your terms consistent',
  'features.presets.body':
    'Presets for Movie, Series, Anime, Documentary or your own, and a glossary for names, honorifics and jargon.',
  'features.batch.title': 'Works in bulk',
  'features.batch.body': 'Translate, extract or mux a whole folder or season at once.',
  'features.tags.title': 'Leaves the formatting alone',
  'features.tags.body':
    'Only the text goes to the model: timestamps, styles and tags come back exactly where they were.',
  'features.privacy.title': 'Stays on your computer',
  'features.privacy.body': 'No account, no telemetry. Your key and your files stay on your machine.',
  'pipeline.kicker': 'How it works',
  'pipeline.title': 'Three steps, one job',
  'pipeline.lead':
    'Use each step on its own or run them all in one go. A subtitle file alone goes straight to translation.',
  'pipeline.extract.title': 'Extract',
  'pipeline.extract.body': 'Pick a text subtitle track from the .mkv and pull it out with mkvtoolnix.',
  'pipeline.translate.title': 'Translate',
  'pipeline.translate.body':
    'The text lines go to the OpenRouter model you chose, with neighboring lines as context. Lines already in the cache are not sent again.',
  'pipeline.remux.title': 'Re\u2011mux',
  'pipeline.remux.body':
    'Mux the translated track back into the video, keeping the other tracks and the attached fonts, or just save the subtitle file.',
  'pipeline.stepLabel': 'Step',
  'privacy.kicker': 'Privacy',
  'privacy.title': 'Your machine, your key, your files',
  'privacy.lead':
    'BakaSub is a desktop app, not a service. There is nothing to sign up for and nothing reports back to the author.',
  'privacy.account.title': 'No account',
  'privacy.account.body': 'Install it and bring your own OpenRouter key. You pay OpenRouter directly for what you translate.',
  'privacy.telemetry.title': 'No telemetry',
  'privacy.telemetry.body':
    'Logs stay on your machine. The exported log file hides your API key and replaces your home folder with ~.',
  'privacy.local.title': 'Key and files stay local',
  'privacy.local.body':
    'Your key is stored encrypted with the OS keyring when there is one. Settings, cache and logs live in your user data folder.',
  'privacy.sentTitle': 'What does leave your computer',
  'privacy.sent.openrouter':
    'OpenRouter receives the subtitle text to translate, your notes, the preset instructions and the glossary. Timestamps and tags are not sent.',
  'privacy.sent.tmdb': 'TMDB receives the title parsed from the file name, only if you add a TMDB key.',
  'privacy.sent.none': 'Nothing is sent until you start a job.',
  'privacy.guideLink': 'Read the full privacy notes in the user guide'
} as const

export type Dictionary = { [K in keyof typeof en]: string }
export type TranslationKey = keyof Dictionary
