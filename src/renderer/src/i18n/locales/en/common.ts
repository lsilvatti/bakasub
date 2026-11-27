export const common = {
  appName: 'BakaSub',
  nav: {
    work: 'Work',
    manage: 'Manage',
    translate: 'Translate',
    extract: 'Extract',
    mux: 'Mux',
    queue: 'Queue',
    presets: 'Presets',
    models: 'Models',
    languages: 'Languages',
    logs: 'Logs',
    settings: 'Settings'
  },
  status: {
    running_one: '{{count}} running',
    running_other: '{{count}} running',
    queued_one: '{{count}} queued',
    queued_other: '{{count}} queued',
    awaiting_one: '{{count}} needs confirmation',
    awaiting_other: '{{count}} need confirmation',
    tools: {
      ok: 'Tools ready',
      missing: 'Tools missing',
      checking: 'Checking tools'
    },
    apiKey: {
      set: 'API key set',
      missing: 'No API key'
    },
    tmdb: { set: 'TMDB on', missing: 'TMDB off' },
    version: 'v{{version}}',
    update: 'update available',
    cooldown: 'Rate limited · {{seconds}} s',
    sessionCost: 'US$ {{value}} this session'
  },
  pages: {
    translate: { title: 'Translate a file' },
    extract: { title: 'Extract subtitles' },
    mux: { title: 'Mux subtitles' },
    queue: { title: 'Queue' },
    presets: { title: 'Presets' },
    models: { title: 'Models' },
    languages: { title: 'Languages' },
    logs: { title: 'Logs' },
    settings: { title: 'Settings' },
    notBuilt: 'This screen is not built yet.'
  },
  actions: {
    cancel: 'Cancel',
    close: 'Close',
    save: 'Save',
    retry: 'Retry',
    reload: 'Reload BakaSub',
    back: 'Back'
  },
  boundary: {
    title: 'Something went wrong on this screen',
    body: 'The error was saved to the log. Try again, or reload BakaSub if it keeps happening.'
  },
  guides: {
    'openrouter-key': 'How to get an OpenRouter key',
    'openrouter-limits': 'Limits, waits and retries',
    'tmdb-key': 'How to get a TMDB key',
    'opensubtitles-key': 'How to get an OpenSubtitles API key',
    'choosing-models': 'How to choose a model',
    openExternal: 'Opens in your browser'
  },
  wait: {
    'rate-limit': 'Waiting for the rate limit',
    cooldown: 'Paused after a rate limit',
    backoff: 'Retrying',
    in: '{{label}} · {{seconds}} s'
  },
  keyForm: {
    saved: 'Key saved',
    notSet: 'No key yet',
    replace: 'Replace key',
    remove: 'Remove key',
    save: 'Check and save',
    checking: 'Checking…',
    howTo: 'How to get one',
    noKeyring:
      'Your system has no keyring, so the key is stored without system protection. Anyone with access to your user account could read it.'
  },
  models: {
    search: 'Search models',
    empty: 'No model matches.',
    offline: 'Showing the saved model list: OpenRouter could not be reached.',
    loading: 'Loading models from OpenRouter…',
    loadFailed: 'Could not load the model list.',
    favorite: 'Add to favorites',
    unfavorite: 'Remove from favorites',
    makeDefault: 'Make default',
    default: 'Default',
    free: 'Free',
    dynamicPrice: 'Price varies',
    perMillion: '${{input}} in · ${{output}} out per 1M tokens',
    context: '{{k}}k context',
    moderated: 'Moderated',
    moderatedHint: 'The provider filters content; some lines may stay untranslated.',
    structured: 'Reliable JSON',
    structuredHint: 'Supports structured output, so fewer retries.',
    sort: {
      label: 'Sort',
      catalog: 'OpenRouter order',
      'price-asc': 'Cheapest first',
      'price-desc': 'Most expensive first',
      'context-desc': 'Largest context first'
    },
    price: {
      label: 'Price per 1M tokens (input and output)',
      min: 'Minimum price',
      max: 'Maximum price'
    },
    family: 'Provider',
    allFamilies: 'All providers',
    filter: {
      freeOnly: 'Free only',
      jsonOnly: 'Reliable JSON',
      hideModerated: 'Not moderated',
      favoritesOnly: 'Favorites only'
    },
    shown: '{{shown}} of {{total}}',
    favoritesOnly: 'Favorites only',
    favoritesCount_one: '{{count}} favorite',
    favoritesCount_other: '{{count}} favorites'
  },
  languages: {
    search: 'Search languages',
    empty: 'No language matches.',
    makeDefault: 'Make default',
    default: 'Default',
    selectedCount_one: '{{count}} language',
    selectedCount_other: '{{count}} languages'
  },
  welcome: {
    title: 'Set up BakaSub',
    stepOf: 'Step {{n}} of {{total}}',
    next: 'Continue',
    skip: 'Skip for now',
    finish: 'Start translating',
    steps: {
      key: 'OpenRouter key',
      models: 'Models',
      languages: 'Languages',
      tools: 'Video tools',
      tmdb: 'TMDB (optional)'
    },
    key: {
      heading: 'Connect your OpenRouter account',
      body: 'BakaSub translates through OpenRouter. Paste your API key; it is checked with OpenRouter and stored encrypted on this computer.',
      label: 'OpenRouter API key'
    },
    models: {
      heading: 'Pick the models you want to use',
      body: 'Star the models to offer when translating and choose the default. You can change them later in Models.',
      needOne: 'Star at least one model to continue.'
    },
    languages: {
      heading: 'Which languages do you translate into?',
      body: 'Choose your target languages and the default. You can change them later in Languages.',
      needOne: 'Choose at least one language to continue.'
    },
    tools: {
      heading: 'Tools for video files',
      body: 'Extracting subtitles from .mkv files and putting translations back needs mkvtoolnix. BakaSub downloads it for you (you will be asked to confirm). Translating subtitle files works without it.',
      ready: 'All tools are installed.',
      install: 'Install tools',
      installing: 'Installing…',
      declined: 'Installation cancelled. You can install the tools later in Settings.',
      offline: 'The tool list could not be updated; install may fail until you are online.',
      phase: {
        downloading: 'Downloading',
        extracting: 'Unpacking',
        probing: 'Checking',
        done: 'Installed'
      },
      state: {
        ok: 'Installed',
        missing: 'Not installed',
        outdated: 'Update available',
        unsupported: 'Not available for this system'
      }
    },
    tmdb: {
      heading: 'Recognize what you translate (optional)',
      body: 'With a free TMDB key, BakaSub finds the movie or show from the file name and gives the AI its synopsis, which improves names and tone. Without it, everything else works the same.',
      label: 'TMDB key',
      attribution: 'This product uses the TMDB API but is not endorsed or certified by TMDB.'
    }
  },
  translate: {
    steps: { file: 'File', track: 'Track', preset: 'Preset & context', output: 'Output' },
    stepsLabel: 'Steps',
    fileChanged: 'The file changed, so the next steps need another look.',
    summary: {
      title: 'Summary',
      file: 'File',
      track: 'Track',
      preset: 'Preset',
      model: 'Model',
      language: 'Into',
      media: 'Media',
      mediaOff: 'Not included',
      output: 'Output',
      empty: 'Choose a file to get started.',
      edit: 'Edit {{step}}'
    },
    existing: {
      title: 'There is already a subtitle in {{language}}',
      track: 'This video already has a {{language}} subtitle track ({{name}}). Translate anyway?',
      file: 'The file name says this subtitle is already in {{language}}. Translate anyway?',
      go: 'Translate anyway',
      inline: 'Already has a {{language}} subtitle'
    },
    next: 'Continue',
    back: 'Back',
    start: 'Estimate cost',
    starting: 'Preparing…',
    another: 'Translate another file',
    file: {
      heading: 'Choose what to translate',
      body: 'A video (.mkv) or a subtitle file (.srt, .ass, .ssa, .vtt).',
      drop: 'Drop a file here',
      or: 'or',
      browse: 'Browse files',
      unsupported: 'This file type is not supported. Use .mkv, .srt, .ass, .ssa or .vtt.',
      video: 'Video',
      subtitle: 'Subtitle file',
      change: 'Choose another file'
    },
    track: {
      heading: 'Pick the subtitle track',
      loading: 'Reading the tracks…',
      none: 'This video has no subtitle track that can be translated.',
      events_one: '{{count}} line',
      events_other: '{{count}} lines',
      default: 'Default',
      forced: 'Forced',
      imageBased: 'Image subtitles are not supported',
      noParser: 'Format not supported',
      fonts_one: 'The video has {{count}} font attached; it is kept for styled (.ass) subtitles.',
      fonts_other:
        'The video has {{count}} fonts attached; they are kept for styled (.ass) subtitles.',
      unknownLang: 'Unknown language'
    },
    external: {
      other: 'Use another subtitle',
      body: 'Translate a subtitle from outside the video instead of one of its tracks.',
      bodyNone:
        'Choose a subtitle file you already have, or search OpenSubtitles for one that matches this video.',
      file: 'Choose a subtitle file',
      search: 'Search OpenSubtitles',
      remove: 'Remove this subtitle',
      origin: { file: 'Your file', opensubtitles: 'OpenSubtitles' }
    },
    preset: {
      heading: 'How to translate',
      suggested: {
        anime: 'Suggested: Japanese animation on TMDB',
        movie: 'Suggested: a movie on TMDB',
        series: 'Suggested: a series on TMDB',
        documentary: 'Suggested: a documentary on TMDB'
      },
      suggestedByName: 'Suggested: the preset name matches the media',
      useSuggestion: 'Back to the suggestion ({{preset}})',
      noSuggestion: 'No TMDB work: the default preset from Settings applies.',
      preset: 'Preset',
      builtin: 'Built-in',
      temperature: 'Temperature',
      mediaType: 'Media type',
      glossary_one: '{{count}} glossary entry',
      glossary_other: '{{count}} glossary entries',
      noGlossary: 'No glossary',
      style: 'Style',
      model: 'Model',
      modelHint: 'Only your favorite models are listed. Manage them in Models.',
      language: 'Translate into',
      languageHint: 'Manage your languages in Languages.',
      notes: 'Notes for the translator (optional)',
      notesHint: 'Characters, names, tone, anything the AI should know. {{count}}/{{max}}',
      notesPlaceholder: 'e.g. Frieren is an elf mage; keep honorifics like -san.'
    },
    output: {
      heading: 'Where to save',
      folder: 'Output folder',
      nextToSource: 'Next to the source file',
      choose: 'Choose folder',
      reset: 'Use the source folder',
      setDefault: 'Make the translated subtitle the default track',
      keepOriginal: 'Keep the original subtitles in the video',
      includeOriginal: 'Also add the original subtitle to the video',
      mode: {
        label: 'What to save',
        mux: 'MKV with the subtitle embedded',
        muxHint: 'A new MKV with the translated track; the original video is not touched.',
        subtitle: 'Subtitle only (for players)',
        subtitleHint:
          'Just the subtitle file next to the video, no remux: faster and needs no extra disk space.',
        both: 'Both',
        bothHint: 'The new MKV and the separate subtitle file.'
      },
      playersNote:
        'Players that load external subtitles pick the file up when it has the same name as the video, e.g. Episode.mkv and Episode.pt-BR.ass in the same folder.',
      previewLabel: 'Files that will be written',
      preview_one: 'Will be saved as',
      preview_other: 'Will be saved as',
      previewExample: 'Example, from the first file',
      overwrite: 'Replace files that already exist',
      overwriteOff: 'If a file with that name exists, the new one is saved as "name (2)".',
      overwriteOn: 'A file with that name is replaced. The source file is never overwritten.',
      summary:
        'BakaSub first estimates the cost and asks you to confirm before anything is sent to the AI.'
    }
  },
  choose: {
    toMany: 'Select several files',
    toSingle: 'Translate a single file',
    toSingleExtract: 'Extract from a single file',
    toSingleMux: 'Mux a single video'
  },
  batch: {
    title: 'Translate several files',
    steps: {
      files: 'Files',
      tracks: 'Tracks',
      works: 'Works',
      preset: 'Preset & context',
      output: 'Output',
      review: 'Review cost',
      run: 'Translation'
    },
    review: {
      heading: 'Review the cost before translating',
      body: 'Nothing goes to the AI until you confirm.',
      lines_one: '{{n}} line to translate',
      lines_other: '{{n}} lines to translate',
      noEstimate: 'Estimate not available',
      nothing: 'All cached',
      dynamic: 'Dynamic price',
      approx: 'approximate (estimated tokenizer)',
      quota: "The free model's daily quota may not be enough for every file.",
      pending_one:
        '{{count}} file is still being prepared; it shows up here when its estimate is ready.',
      pending_other:
        '{{count}} files are still being prepared; they show up here when their estimates are ready.'
    },
    files: {
      heading: 'Choose the files',
      body: 'A folder (with or without subfolders) or several files: videos (.mkv) or subtitle files (.srt, .ass, .ssa, .vtt).',
      drop: 'Drop files here, or choose a folder',
      folder: 'Choose a folder',
      files: 'Choose files',
      recursive: 'Include subfolders',
      fromFolder: 'Folder {{path}}',
      fromFolderRecursive: 'Folder {{path}}, with subfolders',
      fromFiles_one: '{{count}} file',
      fromFiles_other: '{{count}} files',
      reading: 'Reading the files…',
      identifying: 'Recognizing the works on TMDB…',
      none: 'No video (.mkv) or subtitle file was found.',
      truncated: 'The folder has more than 2,000 files: only the first ones are listed.',
      rejected_one: '{{count}} file was refused:',
      rejected_other: '{{count}} files were refused:',
      both: 'This folder has videos and subtitle files. Translate:',
      videos: 'Videos ({{count}})',
      subtitles: 'Subtitle files ({{count}})'
    },
    status: {
      ready: '{{count}} ready',
      attention: '{{count}} need you',
      out: '{{count}} left out',
      works_one: '{{count}} work recognized',
      works_other: '{{count}} works recognized'
    },
    problem: {
      'no-subtitle-track': 'No subtitles',
      'no-supported-track': 'Only image subtitles (PGS, VobSub)',
      'probe-failed': 'The file could not be read',
      'no-track': 'Pick a track',
      'odd-track': 'Different track from the rest of the work',
      'no-episode': 'Episode not recognized'
    },
    issues: {
      title: 'Need you ({{count}})',
      leaveOut: 'Leave out',
      include: 'Include',
      out_one: '{{count}} left out',
      out_other: '{{count}} left out',
      byYou: 'Left out by you',
      previousOutput: 'Looks like an earlier translation (the source is in the batch)'
    },
    tracks: {
      heading: 'Check the subtitle tracks',
      body: 'Picked by your preferred source language. Change the track of a whole work or of a single file.',
      workTrack: 'Track for the whole work',
      pick: 'Pick a track',
      different_one: '{{count}} different',
      different_other: '{{count}} different',
      noLanguage: 'No language',
      image: 'image'
    },
    works: {
      heading: 'Check the works',
      body: 'The files were grouped by the work recognized in their names. Each work has its own context (synopsis, preset, notes). Move files, merge or split works when the grouping is wrong.',
      regroup: 'Group automatically again',
      name: 'Work name',
      manual: 'Manual',
      tmdb: 'Work on TMDB',
      series: 'Series',
      movie: 'Movie',
      notFound: 'Not found on TMDB',
      mergeInto: 'Merge with…',
      moveTo: 'Move to…',
      newWork: 'New work',
      season: 'Season',
      episode: 'Episode',
      files_one: '{{count}} file',
      files_other: '{{count}} files'
    },
    preset: {
      heading: 'How to translate each work',
      body: 'Each work has its own preset, suggested from TMDB. Model, language and notes follow the batch defaults unless you customize the work.',
      defaults: 'Batch model, language and notes',
      defaultsHint: 'Used by the works that are not customized',
      modelLang: 'Model, language and notes',
      useDefaults: 'Use the batch defaults',
      customize: 'Customize this work',
      customized: 'Customized',
      usesDefaults: 'Batch defaults',
      notes: 'Notes for this work (optional)',
      includeMedia: 'Send the title, genres and synopsis (and the episode’s) to the AI',
      withMedia: 'with synopsis',
      withoutMedia: 'without synopsis',
      noMedia: 'No TMDB work: only the notes go as context. Pick the work in Works.'
    },
    output: {
      heading: 'Where to save',
      mirror: 'The subfolders of {{root}} are recreated inside the output folder.',
      confirm: 'Show the total cost and wait for my confirmation before translating',
      confirmOff:
        'Each file starts as soon as its cost is estimated (Settings may still ask above a limit).',
      review: 'Review'
    },
    summary: {
      title: 'Summary',
      source: 'Source',
      files: 'Files',
      tracks: 'Tracks',
      tracksValue: 'Per work, by your preferred language',
      works: 'Works',
      preset: 'Preset',
      customized_one: '{{preset}} · {{count}} work customized',
      customized_other: '{{preset}} · {{count}} works customized',
      model: 'Model',
      language: 'Into',
      output: 'Output',
      write: 'Write',
      nextToEach: 'Next to each file',
      included: '{{count}} included',
      includedOut: '{{count}} included · {{out}} out',
      none: '—',
      empty: 'Choose a folder or files to get started.'
    },
    start_one: 'Translate {{count}} file',
    start_other: 'Translate {{count}} files',
    starting: 'Creating…',
    detail: {
      progress: '{{done}} of {{total}} finished',
      awaiting_one: '{{count}} file is waiting for your confirmation.',
      awaiting_other: '{{count}} files are waiting for your confirmation.',
      cost: 'Estimated total: US$ {{low}} – {{high}}',
      costUnknown: 'Some files use a model with a dynamic price: the total is not known.',
      confirmAll: 'Confirm all',
      cancelAll: 'Cancel all',
      openQueue: 'Open in the queue',
      another: 'Translate other files',
      files: 'Files',
      phase: {
        review: 'Preparing the files and estimating the cost…',
        awaiting: 'Waiting for you to confirm the cost.',
        run: 'Translating…',
        done: 'Batch finished.'
      },
      count: {
        estimating_one: '{{count}} being prepared',
        estimating_other: '{{count}} being prepared',
        awaiting_one: '{{count}} waiting for you',
        awaiting_other: '{{count}} waiting for you',
        translating_one: '{{count}} translating',
        translating_other: '{{count}} translating',
        failed_one: '{{count}} failed or cancelled',
        failed_other: '{{count}} failed or cancelled'
      }
    }
  },
  jobs: {
    step: {
      probe: 'Reading the video',
      extract: 'Extracting the subtitle',
      preflight: 'Estimating the cost',
      translate: 'Translating',
      validate: 'Checking the result',
      write: 'Saving the subtitle',
      mux: 'Putting the subtitle in the video'
    },
    status: {
      queued: 'Waiting in the queue',
      running: 'Running',
      'awaiting-confirmation': 'Waiting for your confirmation',
      completed: 'Done',
      'completed-with-warnings': 'Done, with warnings',
      failed: 'Failed',
      cancelled: 'Cancelled',
      interrupted: 'Interrupted when BakaSub closed'
    },
    lines: '{{done}} of {{total}} lines',
    cached: '{{count}} from the translation cache',
    chunk: '{{completed}} of {{count}} parts',
    inFlight_one: '{{count}} in progress',
    inFlight_other: '{{count}} in progress',
    timeline: {
      label: 'Translation progress across the episode',
      translated: 'Translated',
      cached: 'From cache',
      active: 'In progress',
      failed: 'Not translated',
      lanes: '{{count}} stretches in parallel'
    },
    cost: 'US$ {{value}} so far',
    confirm: 'Translate',
    cancel: 'Cancel job',
    cancelHint: 'Some providers keep generating (and billing) the request already in progress.',
    retry: 'Retry',
    memory: {
      option: 'Ignore saved translations',
      off: 'Lines already in the translation cache are reused at no cost.',
      on: 'Every line goes to the AI again. The new translations replace the saved ones, so the cache is fixed for next time.',
      skipped: 'Saved translations ignored'
    },
    openQueue: 'Open in Queue',
    logs: 'Show logs',
    subtitleFile: 'Subtitle saved to',
    outputs: 'Saved to',
    finalCost: 'Cost: US$ {{value}}',
    costIncomplete: 'Cost of interrupted requests could not be fully determined.',
    costCompare: {
      heading: 'Estimated vs actual cost',
      estimated: 'Estimated',
      actual: 'Actual',
      actualPartial: 'at least US$ {{value}}',
      notEstimated: 'Not estimated',
      nothingSent: 'nothing was sent to the AI',
      within: 'Within the estimated range',
      below: 'Below the estimate by {{pct}}%',
      belowNoPct: 'Below the estimate',
      above: 'Above the estimate by {{pct}}%',
      aboveNoPct: 'Above the estimate',
      free: 'Free model: nothing was billed.',
      unavailable:
        'This model has a variable price, so the preflight only estimated tokens: there is nothing to compare.',
      unfinished:
        'The job did not finish: the actual cost covers only what ran until it stopped, so there is no comparison.',
      partial:
        'Partial value: the cost of some interrupted requests was not reported, so there is no comparison.',
      cached_one: 'The estimate already counted {{count}} line from saved translations.',
      cached_other: 'The estimate already counted {{count}} lines from saved translations.',
      cacheSkipped: 'The estimate ignored saved translations, as this job did.',
      approximate: 'Approximate estimate: this model counts tokens differently from BakaSub.'
    },
    untranslated_one: '{{count}} line stayed in the original language',
    untranslated_other: '{{count}} lines stayed in the original language',
    noExcerpt: 'Line #{{id}} (text not recorded for this job)',
    warnings: 'Warnings',
    remove: 'Remove',
    kind: {
      extract: 'Extract',
      translate: 'Translate',
      mux: 'Mux',
      pipeline: 'Video translation'
    },
    batch: 'Batch',
    setup: {
      title: 'Translation setup',
      model: 'Model',
      preset: 'Preset',
      languages: 'Languages',
      autoSource: 'Detected automatically',
      notRecorded: 'Not recorded'
    },
    attempt: 'attempt {{n}}',
    created: 'Added {{when}}',
    finished: 'Finished {{when}}',
    hint: {
      dailyLimit: 'Free models reset at {{time}} (00:00 UTC). Retry then, or pick a paid model.',
      credits: 'Add credits on OpenRouter',
      privacy: 'Open OpenRouter privacy settings',
      sidecar: 'Install tools',
      interrupted:
        'BakaSub closed while this job ran. Retrying continues from the step that stopped.'
    }
  },
  queue: {
    new: 'New translation',
    clearDone: 'Clear finished',
    clearDoneHint:
      'Removes completed and cancelled jobs from the queue. Failed and interrupted ones stay.',
    cleared_one: '{{count}} job removed',
    cleared_other: '{{count}} jobs removed',
    filters: { all: 'All', active: 'Active', attention: 'Needs you', done: 'Done' },
    empty: {
      all: 'No jobs yet. Translate a file to see it here.',
      active: 'Nothing is running.',
      attention: 'Nothing needs your attention.',
      done: 'No finished jobs.'
    },
    select: 'Select a job to see its details.',
    pill: {
      queued: 'Queued',
      running: 'Running',
      'awaiting-confirmation': 'Confirm',
      completed: 'Done',
      'completed-with-warnings': 'Warnings',
      failed: 'Failed',
      cancelled: 'Cancelled',
      interrupted: 'Interrupted'
    },
    list: 'Jobs'
  },
  preflight: {
    heading: 'Review before translating',
    lines: 'Lines',
    total: '{{count}} total',
    cachedLines: '{{count}} already translated (cache)',
    toTranslate: '{{count}} to translate',
    copies_one: '{{count}} repeat (translated once)',
    copies_other: '{{count}} repeats (translated once)',
    parts_one: '{{count}} request',
    parts_other: '{{count}} requests',
    tokens: 'Tokens',
    tokensIn: '{{low}}–{{high}} in',
    tokensOut: '{{low}}–{{high}} out',
    cost: 'Estimated cost',
    costRange: 'US$ {{low}} – {{high}}',
    free: 'Free',
    dynamic: 'This model has a variable price; only tokens are estimated.',
    approximate: 'Approximate: this model counts tokens differently from BakaSub.',
    model: 'Model {{model}}, prices from {{date}}',
    nothing: 'Everything is already in the translation cache: nothing will be sent to the AI.',
    cacheSkipped: 'cache ignored',
    quotaShort:
      'This needs about {{requests}} free-model requests and you have {{remaining}} left today. The job may stop when the daily limit is reached; what is translated by then is kept.',
    quotaUnknown: 'Free models have a daily limit and BakaSub could not check how much is left.',
    quotaOk: '{{remaining}} of {{limit}} free requests left today'
  },
  promptSummary: {
    heading: 'What the AI receives',
    blocks: 'Blocks in the system prompt',
    glossary_one: 'Glossary ({{count}} term)',
    glossary_other: 'Glossary ({{count}} terms)',
    noGlossary: 'No glossary',
    style: 'Preset style',
    tmdb: 'TMDB details',
    episode: 'Episode synopsis',
    notes: 'Your notes',
    noMedia: 'No media context',
    unavailable: 'The TMDB details you chose could not be fetched and will not be sent: {{reason}}',
    tokens: 'System prompt: {{value}} tokens, repeated in every request',
    show: 'Show the exact system prompt',
    missing: 'This job was created before BakaSub recorded the prompt.'
  },
  termConsistency: {
    heading_one: '{{count}} term translated in different ways',
    heading_other: '{{count}} terms translated in different ways',
    hint: 'Checked locally, without the AI. Pick the right rendering of each term, or “Do not add” to leave it out, and save them in a glossary: the next jobs will use them. Terms already in the glossary are not shown.',
    builtin:
      '“{{name}}” is a built-in preset and cannot be edited. Save the chosen terms in a new preset of your own:',
    newName: 'New preset name',
    newDetails:
      'Media type “{{type}}” and the style instructions of “{{name}}”, plus its glossary.',
    create_one: 'Create preset with {{count}} term',
    create_other: 'Create preset with {{count}} terms',
    created_one: 'Preset “{{name}}” created with {{count}} term. Choose it in your next jobs.',
    created_other: 'Preset “{{name}}” created with {{count}} terms. Choose it in your next jobs.',
    target: 'Glossary of “{{name}}”',
    noneChosen: 'No rendering chosen.',
    add_one: 'Add {{count}} term to the glossary',
    add_other: 'Add {{count}} terms to the glossary',
    added_one: '{{count}} term added to “{{name}}”',
    added_other: '{{count}} terms added to “{{name}}”',
    presetGone: 'The preset of this job (“{{name}}”) no longer exists.',
    lines_one: 'in {{count}} line',
    lines_other: 'in {{count}} lines',
    unmatched_one: '{{count}} line rephrased it',
    unmatched_other: '{{count}} lines rephrased it',
    glossaryIgnored: 'Glossary says “{{target}}”',
    kept: 'kept as in the original',
    choose: 'Use “{{target}}” for “{{term}}”',
    skip: 'Do not add to the glossary',
    skipTerm: 'Do not add “{{term}}” to the glossary',
    dismiss: 'Dismiss',
    ignoredHeading_one: 'The AI did not follow the glossary for {{count}} term',
    ignoredHeading_other: 'The AI did not follow the glossary for {{count}} terms'
  },
  fileCheck: {
    heading: 'File check',
    encHeading: 'Encoding',
    enc: 'Detected as {{name}}',
    bom: 'with BOM',
    conf: {
      certain: 'Certain',
      high: 'High confidence',
      medium: 'Medium confidence',
      low: 'Low confidence'
    },
    lowWarn:
      "BakaSub is not sure about this file's encoding ({{name}}, score {{score}}). Accented letters may come out garbled. Open the file in a text editor, save it as UTF-8 and choose it again. The encoding cannot be chosen here yet.",
    medWarn:
      'The encoding ({{name}}) was guessed with medium confidence. Check that accented letters look right in the lines below.',
    eol: {
      lf: 'LF',
      crlf: 'CRLF',
      cr: 'CR',
      mixed: 'mixed line endings'
    },
    diagHeading: 'Parser warnings',
    diagNone: 'No parser warnings.',
    diagSummary: '{{errors}} errors, {{warns}} warnings, {{infos}} notes',
    show: 'Show list',
    hide: 'Hide list',
    fixAuto: 'fixed on export',
    fixReport: 'report only',
    col: 'col {{column}}',
    truncated: 'Only the first {{count}} are listed.',
    unknown: 'Unrecognized problem ({{code}})',
    loading: 'Checking the file…',
    failed: 'Could not check this file: {{message}}',
    sev: {
      info: 'Note',
      warn: 'Warning',
      error: 'Error'
    },
    covHeading: 'Coverage',
    covBar: '{{translatable}} of {{total}} lines go to the AI',
    covTranslatable: 'To translate: {{count}}',
    covSkipped: 'Nothing to translate: {{count}}',
    covUntranslated: 'Will stay as is: {{count}}',
    code: {
      SRT_INDEX_MISSING: 'Cue number missing',
      SRT_INDEX_DUPLICATE: 'Repeated cue number',
      SRT_INDEX_GAP: 'Gap in the cue numbering',
      SRT_TIME_DOT_SEPARATOR: 'Timestamp uses a dot instead of a comma',
      SRT_TIME_PADDING: 'Timestamp is not zero-padded',
      SRT_ARROW_SPACING: 'Odd spacing around the “-->” arrow',
      SRT_MISSING_BLANK_LINE: 'Blank line between cues is missing',
      SRT_EOL_MIXED: 'Mixed line endings',
      SRT_TRAILING_WS: 'Trailing whitespace',
      SRT_END_BEFORE_START: 'A cue ends before it starts',
      SRT_OVERLAP: 'A cue overlaps the next one',
      SRT_EMPTY_CUE: 'Empty cue',
      SRT_TIMING_EXTRAS: 'Extra settings after the timing',
      SRT_UNPARSEABLE_BLOCK: 'A block could not be read',
      SRT_AMBIGUOUS_INDEX: 'The cue number is ambiguous',
      ASS_EOL_MIXED: 'Mixed line endings',
      ASS_TIME_FORMAT: 'Timestamp in an unusual format',
      ASS_END_BEFORE_START: 'A cue ends before it starts',
      ASS_FORMAT_MISSING: 'The “Format:” line is missing',
      ASS_UNPARSEABLE_EVENT: 'An event line could not be read',
      VTT_HEADER_MISSING: 'The WEBVTT header is missing',
      VTT_TIME_FORMAT: 'Timestamp in an unusual format',
      VTT_MISSING_BLANK_LINE: 'Blank line between cues is missing',
      VTT_ARROW_IN_TEXT: 'The text contains “-->”',
      VTT_EOL_MIXED: 'Mixed line endings',
      VTT_END_BEFORE_START: 'A cue ends before it starts',
      VTT_EMPTY_CUE: 'Empty cue',
      VTT_UNPARSEABLE_BLOCK: 'A block could not be read',
      SEG_MARKER_COLLISION: 'Text looks like an internal marker'
    }
  },
  coverage: {
    heading: 'Not translated',
    summaryUntranslated_one: '{{count}} line will stay in the original language',
    summaryUntranslated_other: '{{count}} lines will stay in the original language',
    summarySkipped_one: '{{count}} line has nothing to translate',
    summarySkipped_other: '{{count}} lines have nothing to translate',
    showSkipped: 'Show lines with nothing to translate',
    hideSkipped: 'Hide them',
    line: 'Line {{line}}',
    truncated: 'Only the first entries are listed.',
    kept: 'Kept as is',
    omitted: 'Left out',
    reason: {
      SKIPPED_NO_LETTERS: 'No words (music notes, dots or only formatting)',
      SKIPPED_EMPTY_BLOCK: 'Empty block',
      SKIPPED_DRAWING: 'Vector drawing',
      SKIPPED_FRAGMENT: 'Karaoke fragment (one letter/syllable per event)',
      UNTRANSLATED_ORPHAN_TEXT: 'Text without timing after a blank line',
      UNTRANSLATED_BAD_TIMING: 'The timing line could not be read',
      UNTRANSLATED_MARKER_COLLISION: 'Text looks like an internal marker',
      UNTRANSLATED_KARAOKE: 'Karaoke line (translating would break the syllable timing)',
      UNTRANSLATED_BAD_EVENT: 'Dialogue line that could not be read',
      UNTRANSLATED_VALIDATION_FAILED: 'The model kept breaking the formatting',
      UNTRANSLATED_CONTENT_FILTERED: "Blocked by the provider's moderation"
    }
  },
  opensubtitles: {
    title: 'Search OpenSubtitles',
    description: 'Subtitles for this video, to use as the source of the translation.',
    needsKey:
      'OpenSubtitles needs your own API key. It is free; add it once and it stays saved in Settings.',
    query: 'Title',
    season: 'Season',
    episode: 'Episode',
    languages: 'Subtitle languages',
    anyLanguage: 'Any language',
    search: 'Search',
    searching: 'Searching…',
    noResults: 'No subtitles found. Try another title or another language.',
    downloads_one: '{{count}} download',
    downloads_other: '{{count}} downloads',
    badge: {
      hashMatch: 'Made for this file',
      trusted: 'Trusted uploader',
      hearingImpaired: 'Hearing impaired',
      machine: 'Machine translated'
    },
    use: 'Use',
    useNamed: 'Use {{name}}',
    downloading: 'Downloading…',
    syncNote:
      'Subtitles marked "Made for this file" match this exact video. The others may be out of sync; BakaSub does not adjust timing.',
    quotaNote:
      'Downloads are limited per day. Connect your OpenSubtitles account in Settings for a bigger limit.',
    sync: {
      short: {
        'fps-mismatch': '{{subtitle}} fps ≠ video {{video}} fps',
        'no-hash-match': 'Not verified for this file'
      },
      long: {
        'fps-mismatch':
          'This subtitle was timed for {{subtitle}} fps, but the video runs at {{video}} fps: the lines will drift out of sync, and translating will not fix it.',
        'no-hash-match':
          'OpenSubtitles could not confirm this subtitle was made for this exact file, so it may be out of sync. Check it before paying for the translation.'
      }
    },
    remaining_one: '{{count}} OpenSubtitles download left today',
    remaining_other: '{{count}} OpenSubtitles downloads left today'
  },
  media: {
    heading: 'What is this?',
    title: 'Title',
    searchPlaceholder: 'Search by title',
    searching: 'Searching TMDB for “{{query}}”…',
    noResults: 'No results for “{{query}}”. Try the original title or add the year.',
    identifying: 'Looking up “{{query}}” from the file name…',
    notFound:
      '“{{query}}” was not recognized from the file name. Search for the title above, or continue without media context.',
    offline: 'TMDB could not be reached. You can continue without media context.',
    noKey:
      'Add a free TMDB key so BakaSub can recognize the movie or show and give the AI its synopsis.',
    addKey: 'Add TMDB key',
    recognized: 'Recognized from the file name',
    series: 'Series',
    movie: 'Movie',
    runtime: '{{minutes}} min',
    movieWithEpisode:
      'The file name looks like an episode ({{code}}), but this title is a movie. Check that the right title was chosen.',
    season: 'Season',
    episode: 'Episode',
    fromFile: 'From the file name',
    changed: 'Changed (file name says {{code}})',
    notInFile: 'Not in the file name',
    absolute:
      'The file name uses absolute numbering (episode {{n}}). Enter the season and episode.',
    enterEpisode: 'Enter the season and episode to load its synopsis.',
    noEpisode: 'TMDB has no episode {{code}} for this title.',
    noEpisodeSynopsis: 'No synopsis for {{code}} on TMDB.',
    include: 'Include these details in the translation context',
    attribution: 'Metadata and images from TMDB'
  },
  settings: {
    range: 'Use a value from {{min}} to {{max}}.',
    tabs: {
      general: 'General',
      keys: 'Keys',
      cache: 'Translation cache',
      tools: 'Video tools',
      advanced: 'Advanced',
      about: 'About'
    },
    general: {
      language: 'Language',
      languageHint: 'Language of BakaSub and its guides.',
      system: 'Same as the system',
      defaultPreset: 'Default preset',
      defaultPresetHint: 'Preset selected when you start a translation.',
      concurrency: 'Jobs at the same time',
      concurrencyHint: 'More jobs finish sooner but share the same OpenRouter limits.',
      autoConfirm: 'Translate without asking below',
      autoConfirmHint:
        'Skip the cost confirmation when the highest estimate is below this amount. Empty: always ask.',
      always: 'Always ask',
      sourceLangs: 'Preferred source languages',
      sourceLangsHint: 'Used to pick the subtitle track of a video automatically, in this order.',
      verbose: 'Detailed logs',
      verboseHint: 'Writes more detail to the log. Useful when reporting a problem.',
      notifications: 'System notifications',
      notificationsHint:
        'Notify when something needs you while BakaSub is in the background. Never while its window is in front.',
      notifyJob: 'Jobs',
      notifyJobHint: 'A file finished, failed or waits for the cost to be confirmed.',
      notifyBatch: 'Batches',
      notifyBatchHint: 'One notification per batch, when it finishes or waits for confirmation.',
      notifyUpdate: 'Updates',
      notifyUpdateHint: 'A new version of BakaSub is available.',
      add: 'Add language',
      remove: 'Remove {{name}}',
      up: 'Move {{name}} up'
    },
    keys: {
      openrouter: 'OpenRouter',
      openrouterHint: 'Used for every translation. Checked with OpenRouter before saving.',
      tmdb: 'TMDB (optional)',
      tmdbHint: 'Recognizes movies and shows to give the AI their synopsis.',
      opensubtitles: 'OpenSubtitles (optional)',
      opensubtitlesHint: 'Finds subtitles for videos that have none, so they can be translated.',
      opensubtitlesLabel: 'OpenSubtitles API key',
      opensubtitlesAccount: 'Account (optional)',
      opensubtitlesAccountHint:
        'Signing in raises the daily download limit. The password is stored encrypted on this computer.',
      opensubtitlesAccountHintNoKeyring: 'Signing in raises the daily download limit.',
      opensubtitlesNoKeyring:
        'Your system has no keyring, so the password is stored without system protection. Anyone with access to your user account could read it.',
      username: 'Username',
      usernameHint: 'Your OpenSubtitles username, not your email.',
      password: 'Password',
      connect: 'Sign in',
      connecting: 'Signing in…',
      connectedAs: 'Signed in as {{username}}',
      disconnect: 'Sign out'
    },
    cache: {
      body: 'Every translated line is kept, so translating the same file again (or a similar one) is instant and free.',
      size: 'Size on disk',
      entries_one: '{{n}} line',
      entries_other: '{{n}} lines',
      oldest: 'Least recently used on {{date}}',
      empty: 'The cache is empty.',
      prune: 'Remove lines not used for',
      days_one: 'more than {{count}} day',
      days_other: 'more than {{count}} days',
      all: 'Clear everything',
      apply: 'Remove',
      busy: 'Wait for the running jobs to finish before cleaning the cache.',
      done_one: 'Removed {{count}} line.',
      done_other: 'Removed {{count}} lines.'
    },
    tools: {
      body: 'mkvtoolnix is needed to read subtitle tracks from videos and to put translations back.',
      check: 'Check again',
      required: 'needs {{version}}'
    },
    advanced: {
      body: 'How BakaSub paces and retries requests to OpenRouter. The defaults follow OpenRouter limits; change them only if you have a reason.',
      quota: 'Free models today',
      quotaValue: '{{used}} of {{limit}} requests used',
      quotaUnknown: 'Could not check the free-model quota.',
      tierAuto: 'Detected: {{tier}}',
      tier: {
        'under-10': 'under US$ 10 of credits bought (50 requests/day)',
        'over-10': 'US$ 10 or more bought (1000 requests/day)',
        unknown: 'unknown (assuming 50 requests/day)'
      },
      freeTier: 'Free-model tier',
      freeTierHint:
        'Detected from your OpenRouter account. Set it by hand only if the detection is wrong.',
      tierOptions: {
        auto: 'Automatic',
        'under-10': 'Under US$ 10 bought',
        'over-10': 'US$ 10 or more bought'
      },
      freeRpm: 'Requests per minute (free models)',
      freeRpmHint:
        'Lower it if you use the same OpenRouter account in other apps at the same time.',
      paidRpm: 'Requests per minute (paid models)',
      paidRpmHint: 'Set a limit if a provider keeps asking you to slow down. Empty: no limit.',
      noLimit: 'No limit',
      maxThrottled: 'Attempts when asked to slow down',
      maxTransient: 'Attempts after errors',
      baseDelay: 'First wait between attempts',
      maxDelay: 'Longest wait between attempts',
      chunkRounds: 'Extra rounds for a failed part',
      chunkRoundsHint: 'Waits 1 minute, then 2, before retrying that part. 0 fails right away.',
      seconds: 's',
      parallel: 'Parallel translation',
      parallelBody:
        'The file is split into stretches translated at the same time. Each stretch runs in order, with the context of the previous lines; cuts land on pauses between lines. Cost stays the same.',
      recommended: 'Recommended: {{min}} to {{max}}',
      perJob: 'Parallel stretches per file',
      perJobHint: '1 translates in sequence, as before.',
      perJobLow: 'Sequential translation: much slower.',
      perJobHigh: 'More cuts in the file: stretches may read less consistently.',
      paidMaxInFlight: 'Simultaneous requests (paid models)',
      paidMaxInFlightHint: 'Across all jobs. Lowers itself when a provider asks to slow down.',
      paidMaxInFlightLow: 'Limits parallelism when more than one job is running.',
      paidMaxInFlightHigh:
        'More likely that providers ask you to slow down and that OpenRouter runs out of reserved balance.',
      freeMaxInFlight: 'Simultaneous requests (free models)',
      freeMaxInFlightHint: 'Across all jobs. OpenRouter allows 20 per minute on free models.',
      freeMaxInFlightLow: 'Sequential translation: much slower.',
      freeMaxInFlightHigh: 'Will hit the 20-per-minute limit and cause waits.',
      starved:
        '{{jobs}} jobs × {{perJob}} stretches = {{wanted}} requests, but the limit is {{limit}}: some stretches will wait for a slot.',
      restore: 'Restore defaults',
      restored: 'Defaults restored.'
    },
    about: {
      version: 'Version {{version}}',
      platform: '{{platform}} {{arch}}, Electron {{electron}}',
      upToDate: 'You have the latest version.',
      update: 'Version {{latest}} is available.',
      offline: 'Could not check for updates.',
      updates: 'Updates',
      checking: 'Checking for updates…',
      download: 'Download',
      downloading: 'Downloading version {{latest}}…',
      downloadProgress: 'Download progress',
      cancelDownload: 'Cancel download',
      ready: 'Version {{latest}} is downloaded and ready to install.',
      install: 'Restart and update',
      retry: 'Try again',
      failed: 'The update failed. Nothing was changed.',
      unsupported:
        'This build cannot update itself (development build or not an AppImage/installer). Download the new version from the release page.',
      devBuild:
        'Development build: in-app updates are off, so it never touches the installed BakaSub.',
      openRelease: 'Open the release page',
      showNotes: 'Show release notes',
      hideNotes: 'Hide release notes',
      autoCheck: 'Check for updates automatically',
      autoCheckHint:
        'Looks for a new version a few seconds after BakaSub starts. It never downloads or installs anything by itself.',
      dialogTitle: 'A new version of BakaSub is available',
      dialogVersions: 'Version {{current}} → {{latest}}',
      dialogPublished: 'Published on {{date}}',
      dialogSize: 'Download size: {{size}}',
      dialogNotes: 'Release notes',
      dialogHint:
        'Nothing is downloaded until you choose to. The download and the restart happen in Settings → About.',
      dialogLater: 'Later',
      dialogUpdate: 'See details and update',
      jobsTitle: 'Jobs are running',
      jobsBody:
        'Installing the update restarts BakaSub and cancels the jobs in progress. They are kept and resume the next time the app starts.',
      jobsBack: 'Back',
      jobsConfirm: 'Cancel jobs and update',
      author: 'Made by {{name}}',
      checkUpdate: 'Check for updates',
      repository: 'Repository',
      coffee: 'Buy me a coffee',
      data: 'Data folder',
      tmdb: 'This product uses the TMDB API but is not endorsed or certified by TMDB. Movie and show metadata and images come from The Movie Database (TMDB).',
      openrouter: 'Translations are made by the models you choose, through OpenRouter.'
    }
  },
  modelsPage: {
    keepOne: 'Keep at least one favorite model: it is what BakaSub translates with.',
    body: 'Star the models you want to offer when translating, and pick the default. Prices are per 1M tokens, from OpenRouter.',
    refresh: 'Refresh list',
    updated: 'List from {{date}}',
    missingTitle: 'No longer on OpenRouter',
    missingBody:
      'These favorites left the OpenRouter catalog. Translating with them will fail; remove them or pick another default.',
    remove: 'Remove',
    quota: '{{remaining}} of {{limit}} free-model requests left today'
  },
  languagesPage: {
    keepOne: 'Keep at least one language to translate into.',
    body: 'The languages you translate into. They are offered when translating; the default is selected first.'
  },
  common: { resize: 'Resize the list' },
  help: {
    label: 'What is this?',
    preset:
      'A preset bundles how to translate: style instructions, glossary and temperature. Built-in presets fit common media; make your own in Presets.',
    model:
      'The AI model that translates. Price, speed and quality vary a lot; see “How to choose a model” in Models.',
    language: 'The language the subtitle is translated into.',
    notes:
      'Sent to the AI with every part of the file. Short facts work best: who the characters are, how they speak, words to keep.',
    temperature:
      'How free the translation is. Low (0.1–0.3) stays literal and consistent; around 0.5 reads more natural; above 1 it starts to improvise.',
    mediaType:
      'What kind of media this preset is for. It is a label for you and is also given to the AI as context.',
    prompt:
      'Instructions about tone and register (formal, casual, keep jokes, adapt idioms…). The fixed rules that protect timing and formatting are added by BakaSub.',
    glossary:
      'Terms the AI must always translate the same way (names, places, techniques). Use a rule for things like “keep honorifics”. Built-in presets have no glossary: to add one, duplicate an existing preset or create a new one.',
    include:
      'Gives the AI the title, genres and synopsis (and the episode synopsis) so names and tone come out right. Nothing is sent to TMDB except the title.',
    concurrency:
      'How many jobs run at once. Each running job sends its own requests, so more jobs reach the OpenRouter limits sooner.',
    autoConfirm:
      'When the highest estimated cost is below this amount, the job starts without asking. The estimate is always shown in the queue.',
    sourceLangs:
      'When a video has several subtitle tracks, BakaSub picks the first one in this list as the source.',
    track:
      'The subtitle track that will be translated. Image-based tracks (PGS, VobSub) cannot be translated.'
  },
  presetNames: { movie: 'Movie', series: 'Series', anime: 'Anime', documentary: 'Documentary' },
  mediaTypes: { movie: 'Movie', series: 'Series', anime: 'Anime', documentary: 'Documentary' },
  presets: {
    copyOf: '{{name}} (copy)',
    copyOfN: '{{name}} (copy {{n}})',
    body: 'How the AI should translate: tone, a glossary of names and terms, and how creative it can be.',
    new: 'New preset',
    import: 'Import',
    exportPlain: 'Export',
    export_one: 'Export {{count}} preset',
    export_other: 'Export {{count}} presets',
    exported_one: 'Exported {{count}} preset.',
    exported_other: 'Exported {{count}} presets.',
    builtin: 'Built-in',
    custom: 'Your presets',
    noCustom: 'No presets of your own yet. Create one, or duplicate a built-in preset.',
    select: 'Select to export',
    readOnly: 'Built-in presets cannot be changed. Duplicate it to make your own version.',
    copiedFrom: 'Copied from {{name}}',
    duplicate: 'Duplicate',
    delete: 'Delete',
    deleteTitle: 'Delete “{{name}}”?',
    deleteBody: 'Jobs that already used it keep their settings.',
    save: 'Save preset',
    saved: 'Preset saved.',
    discard: 'Discard changes',
    unsaved: 'Unsaved changes',
    field: {
      name: 'Name',
      mediaType: 'Media type',
      mediaTypeHint: 'Free text, e.g. anime, sitcom, documentary.',
      temperature: 'Temperature',
      temperatureHint: 'Lower stays closer to the original; higher is freer. 0 to 2.',
      prompt: 'Style instructions',
      promptHint: 'Tone, register, what to keep or adapt. {{count}}/{{max}}',
      promptLanguageTip:
        'Tip: write these instructions in English, even when translating into another language. Models follow English instructions more reliably; the target language is set by the job.'
    },
    error: {
      required: 'Required.',
      tooLong: 'Too long.',
      range: 'Use a value from 0 to 2.'
    },
    glossary: {
      title: 'Glossary',
      hint: 'Names and terms the AI must follow. Each entry needs a translation or a rule.',
      source: 'Term',
      target: 'Translation',
      rule: 'Rule',
      targetPlaceholder: 'Same term to keep it',
      rulePlaceholder: 'e.g. keep the honorific',
      add: 'Add entry',
      remove: 'Remove entry {{n}}',
      empty: 'No entries.',
      tooMany: 'At most 500 entries.',
      rowError: {
        source: 'Enter the term.',
        targetOrRule: 'Add a translation or a rule.',
        tooLong: 'Too long.'
      }
    },
    importDialog: {
      title: 'Import presets',
      body: 'Choose a BakaSub presets file (.json). Imported presets are always added as new ones.',
      conflict: 'When a name already exists',
      rename: 'Import with a new name',
      skip: 'Skip it',
      choose: 'Choose file',
      imported_one: 'Imported {{count}} preset.',
      imported_other: 'Imported {{count}} presets.',
      skipped: 'Not imported',
      entry: 'Entry {{n}}',
      nameTaken: 'name already exists',
      done: 'Done'
    }
  },
  files: {
    videoOnly: 'Choose a video file (.mkv).'
  },
  extractPage: {
    body: 'Save subtitle tracks from a video as separate files, without translating.',
    video: 'Video',
    tracks: 'Tracks to save',
    output: 'Output folder',
    start: 'Extract',
    another: 'Extract from another video',
    naming: {
      label: 'File names',
      player: 'Player-compatible',
      playerHint: 'Episode.en.ass: players load it next to the video on their own.',
      track: 'With the track id',
      trackHint: 'Episode.3.en.ass: always unique, but players do not pick it up by itself.',
      ambiguous_one:
        '{{count}} file has several tracks in the same language: they get a number (Episode.en.1.ass, Episode.en.2.ass), which players may not load on their own.',
      ambiguous_other:
        '{{count}} files have several tracks in the same language: they get a number (Episode.en.1.ass, Episode.en.2.ass), which players may not load on their own.'
    },
    needTrack: 'Choose at least one track.'
  },
  extractBatch: {
    title: 'Extract from several files',
    steps: { files: 'Files', tracks: 'Tracks', output: 'Output' },
    files: {
      heading: 'Choose the videos',
      body: 'A folder (with or without subfolders) or several videos (.mkv).',
      none: 'No video (.mkv) was found.'
    },
    tracks: {
      heading: 'Choose the tracks to save',
      body: 'By default every text track is saved (image subtitles cannot be extracted). Change a kind of track for all files, or fine-tune each file.',
      kinds: 'Kinds of track',
      perFile: 'Files',
      count: '{{selected}} of {{total}} tracks'
    },
    output: {
      heading: 'Where to save',
      next: 'Without an output folder, each file is saved next to its video.',
      review: '{{tracks}} tracks from {{files}} files will be saved.'
    },
    error: {
      item: 'The batch was not created: fix {{file}}.',
      fix: 'Go to the tracks'
    },
    start: 'Extract {{count}} files',
    starting: 'Starting…',
    run: {
      running: 'Extracting…',
      done: 'Batch finished.',
      progress: '{{done}} of {{total}} finished',
      another: 'Extract other files',
      count: {
        queued_one: '{{count}} waiting',
        queued_other: '{{count}} waiting',
        running_one: '{{count}} extracting',
        running_other: '{{count}} extracting',
        done_one: '{{count}} done',
        done_other: '{{count}} done',
        failed_one: '{{count}} failed or cancelled',
        failed_other: '{{count}} failed or cancelled'
      }
    }
  },
  muxBatch: {
    title: 'Mux several videos',
    steps: { files: 'Files', pairing: 'Pairing', output: 'Output' },
    files: {
      heading: 'Choose the videos and subtitles',
      body: 'A folder (with or without subfolders) or several files, with the videos (.mkv) and their subtitle files (.srt, .ass, .ssa, .vtt) together. BakaSub pairs them by name or by season and episode.',
      none: 'No video (.mkv) was found.'
    },
    status: {
      paired_one: '{{count}} video with subtitles',
      paired_other: '{{count}} videos with subtitles',
      unmatchedSubs_one: '{{count}} subtitle needs a video',
      unmatchedSubs_other: '{{count}} subtitles need a video',
      unmatchedVideos_one: '{{count}} video without subtitles',
      unmatchedVideos_other: '{{count}} videos without subtitles'
    },
    pairing: {
      heading: 'Check the pairing',
      body: 'Each video gets its subtitles. Change the language, the track name and the default track, move a subtitle to another video or remove it.',
      unresolved_one: '{{count}} subtitle without a video',
      unresolved_other: '{{count}} subtitles without a video',
      unresolvedBody:
        'BakaSub could not tell which video these belong to (no match, or more than one). Choose the video, or ignore the subtitle.',
      assignTo: 'Video for {{name}}',
      moveTo: 'Move {{name}} to another video',
      addTo: 'Add a subtitle to {{name}}',
      addSubtitle: 'Add a subtitle…',
      pickVideo: 'Choose the video',
      ignore: 'Ignore',
      restore: 'Restore',
      videos_one: '{{count}} video will be created',
      videos_other: '{{count}} videos will be created',
      noneIncluded: 'No video with subtitles to create.',
      count_one: '{{count}} subtitle',
      count_other: '{{count}} subtitles',
      max: 'Over 32 subtitles',
      byEpisode: 'By episode',
      manual: 'Chosen by you',
      undetermined: 'Undetermined',
      out_one: '{{count}} item left out',
      out_other: '{{count}} items left out',
      noSubtitle: 'no subtitles',
      blocked_one: 'Solve {{count}} subtitle to continue.',
      blocked_other: 'Solve {{count}} subtitles to continue.'
    },
    output: {
      heading: 'Where to save',
      next: 'Without an output folder, each new video is saved next to its original.',
      collision_one:
        '{{count}} name repeats in the output folder; a number is added to the later ones:',
      collision_other:
        '{{count}} names repeat in the output folder; a number is added to the later ones:',
      space: 'Each video is copied in full: keep free space of about the size of all the videos.',
      review: '{{videos}} videos with {{subs}} subtitles will be created.'
    },
    error: {
      item: 'The batch was not created: fix {{file}}.',
      fix: 'Go to the pairing'
    },
    start: 'Create {{count}} videos',
    starting: 'Starting…',
    run: {
      running: 'Muxing…',
      done: 'Batch finished.',
      progress: '{{done}} of {{total}} finished',
      another: 'Mux other videos',
      count: {
        queued_one: '{{count}} waiting',
        queued_other: '{{count}} waiting',
        running_one: '{{count}} muxing',
        running_other: '{{count}} muxing',
        done_one: '{{count}} done',
        done_other: '{{count}} done',
        failed_one: '{{count}} failed or cancelled',
        failed_other: '{{count}} failed or cancelled'
      }
    }
  },
  muxPage: {
    body: 'Put subtitle files into a video. The video is copied; the original is never changed.',
    video: 'Video',
    subtitles: 'Subtitles to add',
    add: 'Add subtitle files',
    empty: 'No subtitles yet.',
    language: 'Language',
    name: 'Track name',
    namePlaceholder: 'Defaults to the language name',
    isDefault: 'Default track',
    remove: 'Remove {{name}}',
    guessed: 'From the file name',
    needLanguage: 'Choose the language.',
    options: 'Options',
    setDefault: 'Make the chosen default the only default subtitle',
    keepOriginal: 'Keep the subtitles already in the video',
    start: 'Create video',
    another: 'Mux another video',
    max: 'At most 32 subtitles.'
  },
  logsPage: {
    body: 'What BakaSub did, newest at the bottom. Keys and other secrets are never written here.',
    live: 'Live',
    level: 'Level',
    levels: {
      debug: 'Everything',
      info: 'Info and above',
      warn: 'Warnings and errors',
      error: 'Errors only'
    },
    scope: 'Area',
    allScopes: 'All areas',
    job: 'Job {{id}}',
    filteredJob: 'Only job: {{name}}',
    clearJob: 'Show all jobs',
    openFolder: 'Open log folder',
    export: 'Export logs',
    exported: 'Logs exported.',
    empty: 'No log lines match the filters.',
    jobEmpty: 'This job has no log lines left. They may have been rotated out or deleted.',
    loading: 'Reading the log…',
    loadFailed: 'Could not read the log.',
    details: 'Details',
    count: '{{shown}} of {{total}} lines'
  }
}
