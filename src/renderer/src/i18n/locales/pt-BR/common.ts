import type { LocaleShape } from '../../locale-shape'
import type { common as en } from '../en/common'

export const common: LocaleShape<typeof en> = {
  appName: 'BakaSub',
  nav: {
    work: 'Trabalho',
    manage: 'Gerenciar',
    translate: 'Traduzir',
    extract: 'Extrair',
    mux: 'Mux',
    queue: 'Fila',
    presets: 'Presets',
    models: 'Modelos',
    languages: 'Idiomas',
    logs: 'Logs',
    settings: 'Configurações'
  },
  status: {
    running_one: '{{count}} em andamento',
    running_other: '{{count}} em andamento',
    queued_one: '{{count}} na fila',
    queued_other: '{{count}} na fila',
    awaiting_one: '{{count}} aguarda confirmação',
    awaiting_other: '{{count}} aguardam confirmação',
    tools: {
      ok: 'Ferramentas prontas',
      missing: 'Faltam ferramentas',
      checking: 'Verificando ferramentas'
    },
    apiKey: {
      set: 'API key salva',
      missing: 'Sem API key'
    },
    tmdb: { set: 'TMDB ativo', missing: 'TMDB desligado' },
    version: 'v{{version}}',
    update: 'atualização disponível',
    cooldown: 'Limite de requisições · {{seconds}} s',
    sessionCost: 'US$ {{value}} nesta sessão'
  },
  pages: {
    translate: { title: 'Traduzir um arquivo' },
    extract: { title: 'Extrair legendas' },
    mux: { title: 'Mux de legendas' },
    queue: { title: 'Fila' },
    presets: { title: 'Presets' },
    models: { title: 'Modelos' },
    languages: { title: 'Idiomas' },
    logs: { title: 'Logs' },
    settings: { title: 'Configurações' },
    notBuilt: 'Esta tela ainda não foi construída.'
  },
  actions: {
    cancel: 'Cancelar',
    close: 'Fechar',
    save: 'Salvar',
    retry: 'Tentar de novo',
    reload: 'Recarregar o BakaSub',
    back: 'Voltar'
  },
  boundary: {
    title: 'Algo deu errado nesta tela',
    body: 'O erro foi salvo no log. Tente de novo, ou recarregue o BakaSub se continuar acontecendo.'
  },
  guides: {
    'openrouter-key': 'Como conseguir uma key do OpenRouter',
    'openrouter-limits': 'Limites, esperas e novas tentativas',
    'tmdb-key': 'Como conseguir uma key do TMDB',
    'opensubtitles-key': 'Como conseguir uma API key do OpenSubtitles',
    'choosing-models': 'Como escolher um modelo',
    openExternal: 'Abre no navegador'
  },
  wait: {
    'rate-limit': 'Aguardando o limite de requisições',
    cooldown: 'Pausado após limite de requisições',
    backoff: 'Tentando de novo',
    in: '{{label}} · {{seconds}} s'
  },
  keyForm: {
    saved: 'Key salva',
    notSet: 'Nenhuma key ainda',
    replace: 'Trocar key',
    remove: 'Remover key',
    save: 'Conferir e salvar',
    checking: 'Conferindo…',
    howTo: 'Como conseguir',
    noKeyring:
      'Seu sistema não tem chaveiro, então a key fica salva sem proteção do sistema. Quem tiver acesso à sua conta de usuário pode lê-la.'
  },
  models: {
    search: 'Buscar modelos',
    empty: 'Nenhum modelo encontrado.',
    offline: 'Mostrando a lista salva: não foi possível acessar o OpenRouter.',
    loading: 'Carregando modelos do OpenRouter…',
    loadFailed: 'Não foi possível carregar a lista de modelos.',
    favorite: 'Adicionar aos favoritos',
    unfavorite: 'Remover dos favoritos',
    makeDefault: 'Tornar padrão',
    default: 'Padrão',
    free: 'Grátis',
    dynamicPrice: 'Preço variável',
    perMillion: '${{input}} entrada · ${{output}} saída por 1M tokens',
    context: 'contexto de {{k}}k',
    moderated: 'Moderado',
    moderatedHint: 'O provedor filtra conteúdo; algumas falas podem ficar sem tradução.',
    structured: 'JSON confiável',
    structuredHint: 'Suporta saída estruturada, então há menos novas tentativas.',
    sort: {
      label: 'Ordenar',
      catalog: 'Ordem do OpenRouter',
      'price-asc': 'Mais baratos primeiro',
      'price-desc': 'Mais caros primeiro',
      'context-desc': 'Maior contexto primeiro'
    },
    price: {
      label: 'Preço por 1M tokens (entrada e saída)',
      min: 'Preço mínimo',
      max: 'Preço máximo'
    },
    family: 'Provedor',
    allFamilies: 'Todos os provedores',
    filter: {
      freeOnly: 'Só gratuitos',
      jsonOnly: 'JSON confiável',
      hideModerated: 'Sem moderação',
      favoritesOnly: 'Só favoritos'
    },
    shown: '{{shown}} de {{total}}',
    favoritesOnly: 'Só favoritos',
    favoritesCount_one: '{{count}} favorito',
    favoritesCount_other: '{{count}} favoritos'
  },
  languages: {
    search: 'Buscar idiomas',
    empty: 'Nenhum idioma encontrado.',
    makeDefault: 'Tornar padrão',
    default: 'Padrão',
    selectedCount_one: '{{count}} idioma',
    selectedCount_other: '{{count}} idiomas'
  },
  welcome: {
    title: 'Configurar o BakaSub',
    stepOf: 'Passo {{n}} de {{total}}',
    next: 'Continuar',
    skip: 'Pular por enquanto',
    finish: 'Começar a traduzir',
    steps: {
      key: 'Key do OpenRouter',
      models: 'Modelos',
      languages: 'Idiomas',
      tools: 'Ferramentas de vídeo',
      tmdb: 'TMDB (opcional)'
    },
    key: {
      heading: 'Conecte sua conta do OpenRouter',
      body: 'O BakaSub traduz pelo OpenRouter. Cole sua API key; ela é conferida com o OpenRouter e guardada criptografada neste computador.',
      label: 'API key do OpenRouter'
    },
    models: {
      heading: 'Escolha os modelos que quer usar',
      body: 'Marque com estrela os modelos que devem aparecer ao traduzir e escolha o padrão. Dá para mudar depois em Modelos.',
      needOne: 'Marque pelo menos um modelo para continuar.'
    },
    languages: {
      heading: 'Para quais idiomas você traduz?',
      body: 'Escolha seus idiomas de destino e o padrão. Dá para mudar depois em Idiomas.',
      needOne: 'Escolha pelo menos um idioma para continuar.'
    },
    tools: {
      heading: 'Ferramentas para vídeos',
      body: 'Para extrair legendas de arquivos .mkv e colocar as traduções de volta é preciso o mkvtoolnix. O BakaSub baixa para você (vai pedir confirmação). Traduzir arquivos de legenda funciona sem ele.',
      ready: 'Todas as ferramentas estão instaladas.',
      install: 'Instalar ferramentas',
      installing: 'Instalando…',
      declined: 'Instalação cancelada. Você pode instalar depois em Configurações.',
      offline:
        'Não foi possível atualizar a lista de ferramentas; a instalação pode falhar até você estar online.',
      phase: {
        downloading: 'Baixando',
        extracting: 'Descompactando',
        probing: 'Verificando',
        done: 'Instalado'
      },
      state: {
        ok: 'Instalada',
        missing: 'Não instalada',
        outdated: 'Atualização disponível',
        unsupported: 'Indisponível para este sistema'
      }
    },
    tmdb: {
      heading: 'Reconhecer o que você traduz (opcional)',
      body: 'Com uma key gratuita do TMDB, o BakaSub encontra o filme ou a série pelo nome do arquivo e passa a sinopse para a IA, o que melhora nomes e tom. Sem ela, todo o resto funciona igual.',
      label: 'Key do TMDB',
      attribution: 'Este produto usa a API do TMDB, mas não é endossado nem certificado pelo TMDB.'
    }
  },
  translate: {
    steps: { file: 'Arquivo', track: 'Faixa', preset: 'Preset e contexto', output: 'Saída' },
    stepsLabel: 'Passos',
    fileChanged: 'O arquivo mudou, então os próximos passos precisam ser revistos.',
    summary: {
      title: 'Resumo',
      file: 'Arquivo',
      track: 'Faixa',
      preset: 'Preset',
      model: 'Modelo',
      language: 'Para',
      media: 'Mídia',
      mediaOff: 'Não incluída',
      output: 'Saída',
      empty: 'Escolha um arquivo para começar.',
      edit: 'Editar {{step}}'
    },
    existing: {
      title: 'Já existe uma legenda em {{language}}',
      track:
        'Este vídeo já tem uma faixa de legenda em {{language}} ({{name}}). Traduzir mesmo assim?',
      file: 'O nome do arquivo indica que esta legenda já está em {{language}}. Traduzir mesmo assim?',
      go: 'Traduzir mesmo assim',
      inline: 'Já tem legenda em {{language}}'
    },
    next: 'Continuar',
    back: 'Voltar',
    start: 'Estimar custo',
    starting: 'Preparando…',
    another: 'Traduzir outro arquivo',
    file: {
      heading: 'Escolha o que traduzir',
      body: 'Um vídeo (.mkv) ou um arquivo de legenda (.srt, .ass, .ssa, .vtt).',
      drop: 'Solte um arquivo aqui',
      or: 'ou',
      browse: 'Procurar arquivos',
      unsupported: 'Este tipo de arquivo não é aceito. Use .mkv, .srt, .ass, .ssa ou .vtt.',
      video: 'Vídeo',
      subtitle: 'Arquivo de legenda',
      change: 'Escolher outro arquivo'
    },
    track: {
      heading: 'Escolha a faixa de legenda',
      loading: 'Lendo as faixas…',
      none: 'Este vídeo não tem nenhuma faixa de legenda que possa ser traduzida.',
      events_one: '{{count}} fala',
      events_other: '{{count}} falas',
      default: 'Padrão',
      forced: 'Forçada',
      imageBased: 'Legendas em imagem não são aceitas',
      noParser: 'Formato não aceito',
      fonts_one:
        'O vídeo tem {{count}} fonte anexada; ela é mantida para legendas com estilo (.ass).',
      fonts_other:
        'O vídeo tem {{count}} fontes anexadas; elas são mantidas para legendas com estilo (.ass).',
      unknownLang: 'Idioma desconhecido'
    },
    external: {
      other: 'Usar outra legenda',
      body: 'Traduza uma legenda de fora do vídeo em vez de uma das faixas dele.',
      bodyNone:
        'Escolha um arquivo de legenda que você já tem ou busque no OpenSubtitles uma que combine com este vídeo.',
      file: 'Escolher arquivo de legenda',
      search: 'Buscar no OpenSubtitles',
      remove: 'Remover esta legenda',
      origin: { file: 'Seu arquivo', opensubtitles: 'OpenSubtitles' }
    },
    preset: {
      heading: 'Como traduzir',
      suggested: {
        anime: 'Sugerido: animação japonesa no TMDB',
        movie: 'Sugerido: filme no TMDB',
        series: 'Sugerido: série no TMDB',
        documentary: 'Sugerido: documentário no TMDB'
      },
      suggestedByName: 'Sugerido: o nome do preset bate com a mídia',
      useSuggestion: 'Voltar à sugestão ({{preset}})',
      noSuggestion: 'Sem obra do TMDB: vale o preset padrão das configurações.',
      preset: 'Preset',
      builtin: 'Embutido',
      temperature: 'Temperatura',
      mediaType: 'Tipo de mídia',
      glossary_one: '{{count}} entrada no glossário',
      glossary_other: '{{count}} entradas no glossário',
      noGlossary: 'Sem glossário',
      style: 'Estilo',
      model: 'Modelo',
      modelHint: 'Só seus modelos favoritos aparecem. Gerencie em Modelos.',
      language: 'Traduzir para',
      languageHint: 'Gerencie seus idiomas em Idiomas.',
      notes: 'Notas para o tradutor (opcional)',
      notesHint: 'Personagens, nomes, tom, tudo que a IA deve saber. {{count}}/{{max}}',
      notesPlaceholder: 'ex.: Frieren é uma maga elfa; mantenha honoríficos como -san.'
    },
    output: {
      heading: 'Onde salvar',
      folder: 'Pasta de saída',
      nextToSource: 'Ao lado do arquivo de origem',
      choose: 'Escolher pasta',
      reset: 'Usar a pasta de origem',
      setDefault: 'Tornar a legenda traduzida a faixa padrão',
      keepOriginal: 'Manter as legendas originais no vídeo',
      includeOriginal: 'Incluir também a legenda original no vídeo',
      mode: {
        label: 'O que salvar',
        mux: 'MKV com a legenda embutida',
        muxHint: 'Um novo MKV com a faixa traduzida; o vídeo original não é alterado.',
        subtitle: 'Só a legenda (para players)',
        subtitleHint:
          'Apenas o arquivo de legenda ao lado do vídeo, sem remux: mais rápido e sem usar espaço extra em disco.',
        both: 'Ambos',
        bothHint: 'O novo MKV e o arquivo de legenda separado.'
      },
      playersNote:
        'Players que leem legenda externa carregam o arquivo quando ele tem o mesmo nome do vídeo, por ex. Episodio.mkv e Episodio.pt-BR.ass na mesma pasta.',
      previewLabel: 'Arquivos que serão gravados',
      preview_one: 'Será salvo como',
      preview_other: 'Serão salvos como',
      previewExample: 'Exemplo, a partir do primeiro arquivo',
      overwrite: 'Substituir arquivos que já existem',
      overwriteOff: 'Se já existir um arquivo com esse nome, o novo é salvo como "nome (2)".',
      overwriteOn:
        'Um arquivo com esse nome é substituído. O arquivo de origem nunca é sobrescrito.',
      summary:
        'O BakaSub primeiro estima o custo e pede sua confirmação antes de enviar qualquer coisa para a IA.'
    }
  },
  choose: {
    toMany: 'Selecionar vários arquivos',
    toSingle: 'Traduzir um arquivo só',
    toSingleExtract: 'Extrair de um arquivo só',
    toSingleMux: 'Fazer mux de um vídeo só'
  },
  batch: {
    title: 'Traduzir vários arquivos',
    steps: {
      files: 'Arquivos',
      tracks: 'Faixas',
      works: 'Obras',
      preset: 'Preset e contexto',
      output: 'Saída',
      review: 'Revisar custo',
      run: 'Tradução'
    },
    review: {
      heading: 'Revise o custo antes de traduzir',
      body: 'Nada vai para a IA antes da sua confirmação.',
      lines_one: '{{n}} fala a traduzir',
      lines_other: '{{n}} falas a traduzir',
      noEstimate: 'Estimativa indisponível',
      nothing: 'Tudo em cache',
      dynamic: 'Preço dinâmico',
      approx: 'aproximado (tokenizer estimado)',
      quota: 'A cota diária do modelo grátis pode não bastar para todos os arquivos.',
      pending_one:
        '{{count}} arquivo ainda está sendo preparado; ele aparece aqui quando a estimativa ficar pronta.',
      pending_other:
        '{{count}} arquivos ainda estão sendo preparados; eles aparecem aqui quando as estimativas ficarem prontas.'
    },
    files: {
      heading: 'Escolha os arquivos',
      body: 'Uma pasta (com ou sem subpastas) ou vários arquivos: vídeos (.mkv) ou legendas (.srt, .ass, .ssa, .vtt).',
      drop: 'Arraste os arquivos para cá ou escolha uma pasta',
      folder: 'Escolher pasta',
      files: 'Escolher arquivos',
      recursive: 'Incluir subpastas',
      fromFolder: 'Pasta {{path}}',
      fromFolderRecursive: 'Pasta {{path}}, com subpastas',
      fromFiles_one: '{{count}} arquivo',
      fromFiles_other: '{{count}} arquivos',
      reading: 'Lendo os arquivos…',
      identifying: 'Reconhecendo as obras no TMDB…',
      none: 'Nenhum vídeo (.mkv) ou legenda foi encontrado.',
      truncated: 'A pasta tem mais de 2.000 arquivos: só os primeiros foram listados.',
      rejected_one: '{{count}} arquivo foi recusado:',
      rejected_other: '{{count}} arquivos foram recusados:',
      both: 'Esta pasta tem vídeos e legendas. Traduzir:',
      videos: 'Vídeos ({{count}})',
      subtitles: 'Legendas ({{count}})'
    },
    status: {
      ready: '{{count}} prontos',
      attention: '{{count}} precisam de você',
      out: '{{count}} ficam de fora',
      works_one: '{{count}} obra reconhecida',
      works_other: '{{count}} obras reconhecidas'
    },
    problem: {
      'no-subtitle-track': 'Sem legendas',
      'no-supported-track': 'Só legendas em imagem (PGS, VobSub)',
      'probe-failed': 'Não foi possível ler o arquivo',
      'no-track': 'Escolha uma faixa',
      'odd-track': 'Faixa diferente das outras da obra',
      'no-episode': 'Episódio não reconhecido'
    },
    issues: {
      title: 'Precisam de você ({{count}})',
      leaveOut: 'Deixar de fora',
      include: 'Incluir',
      out_one: '{{count}} fica de fora',
      out_other: '{{count}} ficam de fora',
      byYou: 'Deixado de fora por você',
      previousOutput: 'Parece uma tradução anterior (o original está no lote)'
    },
    tracks: {
      heading: 'Confira as faixas de legenda',
      body: 'Escolhidas pelo seu idioma de origem preferido. Troque a faixa de uma obra inteira ou de um arquivo só.',
      workTrack: 'Faixa da obra inteira',
      pick: 'Escolha a faixa',
      different_one: '{{count}} diferente',
      different_other: '{{count}} diferentes',
      noLanguage: 'Sem idioma',
      image: 'imagem'
    },
    works: {
      heading: 'Confira as obras',
      body: 'Os arquivos foram agrupados pela obra reconhecida no nome. Cada obra tem o próprio contexto (sinopse, preset, notas). Mova arquivos, junte ou separe obras quando o agrupamento errar.',
      regroup: 'Reagrupar automaticamente',
      name: 'Nome da obra',
      manual: 'Manual',
      tmdb: 'Obra no TMDB',
      series: 'Série',
      movie: 'Filme',
      notFound: 'Não encontrada no TMDB',
      mergeInto: 'Juntar com…',
      moveTo: 'Mover para…',
      newWork: 'Nova obra',
      season: 'Temporada',
      episode: 'Episódio',
      files_one: '{{count}} arquivo',
      files_other: '{{count}} arquivos'
    },
    preset: {
      heading: 'Como traduzir cada obra',
      body: 'Cada obra tem o próprio preset, sugerido pelo TMDB. Modelo, idioma e notas seguem o padrão do lote, a menos que você personalize a obra.',
      defaults: 'Modelo, idioma e notas do lote',
      defaultsHint: 'Usados pelas obras não personalizadas',
      modelLang: 'Modelo, idioma e notas',
      useDefaults: 'Usar o padrão do lote',
      customize: 'Personalizar esta obra',
      customized: 'Personalizada',
      usesDefaults: 'Padrão do lote',
      notes: 'Notas desta obra (opcional)',
      includeMedia: 'Enviar título, gêneros e sinopse (e a do episódio) para a IA',
      withMedia: 'com sinopse',
      withoutMedia: 'sem sinopse',
      noMedia: 'Sem obra do TMDB: só as notas vão como contexto. Escolha a obra em Obras.'
    },
    output: {
      heading: 'Onde salvar',
      mirror: 'As subpastas de {{root}} são recriadas dentro da pasta de saída.',
      confirm: 'Mostrar o custo total e esperar minha confirmação antes de traduzir',
      confirmOff:
        'Cada arquivo começa assim que o custo é estimado (as Configurações ainda podem pedir confirmação acima de um limite).',
      review: 'Revisão'
    },
    summary: {
      title: 'Resumo',
      source: 'Origem',
      files: 'Arquivos',
      tracks: 'Faixas',
      tracksValue: 'Por obra, pelo seu idioma preferido',
      works: 'Obras',
      preset: 'Preset',
      customized_one: '{{preset}} · {{count}} obra personalizada',
      customized_other: '{{preset}} · {{count}} obras personalizadas',
      model: 'Modelo',
      language: 'Para',
      output: 'Saída',
      write: 'Gravar',
      nextToEach: 'Ao lado de cada arquivo',
      included: '{{count}} incluídos',
      includedOut: '{{count}} incluídos · {{out}} de fora',
      none: '—',
      empty: 'Escolha uma pasta ou arquivos para começar.'
    },
    start_one: 'Traduzir {{count}} arquivo',
    start_other: 'Traduzir {{count}} arquivos',
    starting: 'Criando…',
    detail: {
      progress: '{{done}} de {{total}} concluídos',
      awaiting_one: '{{count}} arquivo aguarda sua confirmação.',
      awaiting_other: '{{count}} arquivos aguardam sua confirmação.',
      cost: 'Total estimado: US$ {{low}} – {{high}}',
      costUnknown: 'Alguns arquivos usam um modelo com preço dinâmico: o total não é conhecido.',
      confirmAll: 'Confirmar todos',
      cancelAll: 'Cancelar todos',
      openQueue: 'Abrir na fila',
      another: 'Traduzir outros arquivos',
      files: 'Arquivos',
      phase: {
        review: 'Preparando os arquivos e estimando o custo…',
        awaiting: 'Aguardando você confirmar o custo.',
        run: 'Traduzindo…',
        done: 'Lote concluído.'
      },
      count: {
        estimating_one: '{{count}} em preparação',
        estimating_other: '{{count}} em preparação',
        awaiting_one: '{{count}} aguardando você',
        awaiting_other: '{{count}} aguardando você',
        translating_one: '{{count}} traduzindo',
        translating_other: '{{count}} traduzindo',
        failed_one: '{{count}} com falha ou cancelado',
        failed_other: '{{count}} com falha ou cancelados'
      }
    }
  },
  jobs: {
    step: {
      probe: 'Lendo o vídeo',
      extract: 'Extraindo a legenda',
      preflight: 'Estimando o custo',
      translate: 'Traduzindo',
      validate: 'Conferindo o resultado',
      write: 'Salvando a legenda',
      mux: 'Colocando a legenda no vídeo'
    },
    status: {
      queued: 'Aguardando na fila',
      running: 'Em andamento',
      'awaiting-confirmation': 'Aguardando sua confirmação',
      completed: 'Concluído',
      'completed-with-warnings': 'Concluído, com avisos',
      failed: 'Falhou',
      cancelled: 'Cancelado',
      interrupted: 'Interrompido quando o BakaSub fechou'
    },
    lines: '{{done}} de {{total}} falas',
    cached: '{{count}} do cache de tradução',
    chunk: '{{completed}} de {{count}} partes',
    inFlight_one: '{{count}} em andamento',
    inFlight_other: '{{count}} em andamento',
    timeline: {
      label: 'Andamento da tradução ao longo do episódio',
      translated: 'Traduzido',
      cached: 'Do cache',
      active: 'Em andamento',
      failed: 'Não traduzido',
      lanes: '{{count}} trechos em paralelo'
    },
    cost: 'US$ {{value}} até agora',
    confirm: 'Traduzir',
    cancel: 'Cancelar job',
    cancelHint: 'Alguns provedores continuam gerando (e cobrando) a requisição já em andamento.',
    retry: 'Tentar de novo',
    memory: {
      option: 'Ignorar traduções salvas',
      off: 'As falas que já estão no cache de tradução são reaproveitadas sem custo.',
      on: 'Todas as falas vão para a IA de novo. As traduções novas substituem as salvas, então o cache fica corrigido para a próxima vez.',
      skipped: 'Traduções salvas ignoradas'
    },
    openQueue: 'Abrir na fila',
    logs: 'Ver logs',
    subtitleFile: 'Legenda salva em',
    outputs: 'Salvo em',
    finalCost: 'Custo: US$ {{value}}',
    costIncomplete: 'Não foi possível determinar todo o custo das requisições interrompidas.',
    costCompare: {
      heading: 'Custo estimado vs real',
      estimated: 'Estimado',
      actual: 'Real',
      actualPartial: 'pelo menos US$ {{value}}',
      notEstimated: 'Não estimado',
      nothingSent: 'nada foi enviado à IA',
      within: 'Dentro da faixa estimada',
      below: '{{pct}}% abaixo da estimativa',
      belowNoPct: 'Abaixo da estimativa',
      above: '{{pct}}% acima da estimativa',
      aboveNoPct: 'Acima da estimativa',
      free: 'Modelo gratuito: nada foi cobrado.',
      unavailable:
        'Este modelo tem preço variável, então o preflight estimou só os tokens: não há com o que comparar.',
      unfinished:
        'O job não terminou: o custo real cobre só o que rodou até parar, então não há comparação.',
      partial:
        'Valor parcial: o custo de algumas requisições interrompidas não foi informado, então não há comparação.',
      cached_one: 'A estimativa já contava {{count}} fala vinda das traduções salvas.',
      cached_other: 'A estimativa já contava {{count}} falas vindas das traduções salvas.',
      cacheSkipped: 'A estimativa ignorou as traduções salvas, como este job.',
      approximate:
        'Estimativa aproximada: este modelo conta tokens de um jeito diferente do BakaSub.'
    },
    untranslated_one: '{{count}} fala ficou no idioma original',
    untranslated_other: '{{count}} falas ficaram no idioma original',
    noExcerpt: 'Fala nº {{id}} (texto não registrado neste job)',
    warnings: 'Avisos',
    remove: 'Remover',
    kind: {
      extract: 'Extração',
      translate: 'Tradução',
      mux: 'Mux',
      pipeline: 'Tradução de vídeo'
    },
    batch: 'Lote',
    setup: {
      title: 'Configuração da tradução',
      model: 'Modelo',
      preset: 'Preset',
      languages: 'Idiomas',
      autoSource: 'Detectado automaticamente',
      notRecorded: 'Não registrado'
    },
    attempt: 'tentativa {{n}}',
    created: 'Adicionado {{when}}',
    finished: 'Concluído {{when}}',
    hint: {
      dailyLimit:
        'Os modelos gratuitos zeram às {{time}} (00:00 UTC). Tente de novo nesse horário ou escolha um modelo pago.',
      credits: 'Adicionar créditos no OpenRouter',
      privacy: 'Abrir as configurações de privacidade do OpenRouter',
      sidecar: 'Instalar ferramentas',
      interrupted:
        'O BakaSub fechou enquanto este job rodava. Tentar de novo continua da etapa em que parou.'
    }
  },
  queue: {
    new: 'Nova tradução',
    clearDone: 'Limpar concluídos',
    clearDoneHint:
      'Remove da fila os jobs concluídos e cancelados. Os que falharam ou foram interrompidos ficam.',
    cleared_one: '{{count}} job removido',
    cleared_other: '{{count}} jobs removidos',
    filters: { all: 'Todos', active: 'Ativos', attention: 'Precisam de você', done: 'Concluídos' },
    empty: {
      all: 'Nenhum job ainda. Traduza um arquivo para vê-lo aqui.',
      active: 'Nada em andamento.',
      attention: 'Nada precisa da sua atenção.',
      done: 'Nenhum job concluído.'
    },
    select: 'Selecione um job para ver os detalhes.',
    pill: {
      queued: 'Na fila',
      running: 'Rodando',
      'awaiting-confirmation': 'Confirmar',
      completed: 'Concluído',
      'completed-with-warnings': 'Avisos',
      failed: 'Falhou',
      cancelled: 'Cancelado',
      interrupted: 'Interrompido'
    },
    list: 'Jobs'
  },
  preflight: {
    heading: 'Revise antes de traduzir',
    lines: 'Falas',
    total: '{{count}} no total',
    cachedLines: '{{count}} já traduzidas (cache)',
    toTranslate: '{{count}} para traduzir',
    copies_one: '{{count}} repetição (traduzida uma vez só)',
    copies_other: '{{count}} repetições (traduzidas uma vez só)',
    parts_one: '{{count}} requisição',
    parts_other: '{{count}} requisições',
    tokens: 'Tokens',
    tokensIn: '{{low}}–{{high}} de entrada',
    tokensOut: '{{low}}–{{high}} de saída',
    cost: 'Custo estimado',
    costRange: 'US$ {{low}} – {{high}}',
    free: 'Grátis',
    dynamic: 'Este modelo tem preço variável; só os tokens são estimados.',
    approximate: 'Aproximado: este modelo conta tokens de um jeito diferente do BakaSub.',
    model: 'Modelo {{model}}, preços de {{date}}',
    nothing: 'Tudo já está no cache de tradução: nada será enviado para a IA.',
    cacheSkipped: 'cache ignorado',
    quotaShort:
      'Isto precisa de cerca de {{requests}} requisições de modelos gratuitos e restam {{remaining}} hoje. O job pode parar quando o limite diário for atingido; o que já tiver sido traduzido fica salvo.',
    quotaUnknown:
      'Os modelos gratuitos têm limite diário e o BakaSub não conseguiu verificar quanto resta.',
    quotaOk: '{{remaining}} de {{limit}} requisições gratuitas restantes hoje'
  },
  promptSummary: {
    heading: 'O que a IA recebe',
    blocks: 'Blocos no prompt de sistema',
    glossary_one: 'Glossário ({{count}} termo)',
    glossary_other: 'Glossário ({{count}} termos)',
    noGlossary: 'Sem glossário',
    style: 'Estilo do preset',
    tmdb: 'Detalhes do TMDB',
    episode: 'Sinopse do episódio',
    notes: 'Suas notas',
    noMedia: 'Sem contexto de mídia',
    unavailable:
      'Não foi possível buscar os detalhes do TMDB escolhidos; eles não serão enviados: {{reason}}',
    tokens: 'Prompt de sistema: {{value}} tokens, repetido em cada requisição',
    show: 'Ver o prompt de sistema exato',
    missing: 'Este job foi criado antes de o BakaSub registrar o prompt.'
  },
  termConsistency: {
    heading_one: '{{count}} termo traduzido de formas diferentes',
    heading_other: '{{count}} termos traduzidos de formas diferentes',
    hint: 'Verificado localmente, sem a IA. Escolha a forma certa de cada termo, ou “Não adicionar” para deixá-lo de fora, e salve num glossário: os próximos jobs vão usá-las. Termos que já estão no glossário não aparecem.',
    builtin:
      '“{{name}}” é um preset embutido e não pode ser editado. Salve os termos escolhidos num preset seu:',
    newName: 'Nome do novo preset',
    newDetails:
      'Tipo de mídia “{{type}}” e as instruções de estilo de “{{name}}”, mais o glossário dele.',
    create_one: 'Criar preset com {{count}} termo',
    create_other: 'Criar preset com {{count}} termos',
    created_one: 'Preset “{{name}}” criado com {{count}} termo. Escolha-o nos próximos jobs.',
    created_other: 'Preset “{{name}}” criado com {{count}} termos. Escolha-o nos próximos jobs.',
    target: 'Glossário de “{{name}}”',
    noneChosen: 'Nenhuma forma escolhida.',
    add_one: 'Adicionar {{count}} termo ao glossário',
    add_other: 'Adicionar {{count}} termos ao glossário',
    added_one: '{{count}} termo adicionado a “{{name}}”',
    added_other: '{{count}} termos adicionados a “{{name}}”',
    presetGone: 'O preset deste job (“{{name}}”) não existe mais.',
    lines_one: 'em {{count}} fala',
    lines_other: 'em {{count}} falas',
    unmatched_one: '{{count}} fala reformulou',
    unmatched_other: '{{count}} falas reformularam',
    glossaryIgnored: 'Glossário diz “{{target}}”',
    kept: 'mantido como no original',
    choose: 'Usar “{{target}}” para “{{term}}”',
    skip: 'Não adicionar ao glossário',
    skipTerm: 'Não adicionar “{{term}}” ao glossário',
    dismiss: 'Dispensar',
    ignoredHeading_one: 'A IA não seguiu o glossário em {{count}} termo',
    ignoredHeading_other: 'A IA não seguiu o glossário em {{count}} termos'
  },
  fileCheck: {
    heading: 'Verificação do arquivo',
    encHeading: 'Encoding',
    enc: 'Detectado como {{name}}',
    bom: 'com BOM',
    conf: {
      certain: 'Certeza',
      high: 'Confiança alta',
      medium: 'Confiança média',
      low: 'Confiança baixa'
    },
    lowWarn:
      'O BakaSub não tem certeza do encoding deste arquivo ({{name}}, pontuação {{score}}). Letras acentuadas podem sair corrompidas. Abra o arquivo num editor de texto, salve como UTF-8 e escolha-o de novo. Ainda não dá para escolher o encoding aqui.',
    medWarn:
      'O encoding ({{name}}) foi adivinhado com confiança média. Confira se as letras acentuadas estão certas nas falas abaixo.',
    eol: {
      lf: 'LF',
      crlf: 'CRLF',
      cr: 'CR',
      mixed: 'quebras de linha misturadas'
    },
    diagHeading: 'Avisos do parser',
    diagNone: 'Nenhum aviso do parser.',
    diagSummary: '{{errors}} erros, {{warns}} avisos, {{infos}} notas',
    show: 'Mostrar lista',
    hide: 'Ocultar lista',
    fixAuto: 'corrigido na exportação',
    fixReport: 'só informativo',
    col: 'col. {{column}}',
    truncated: 'Só os primeiros {{count}} estão listados.',
    unknown: 'Problema não reconhecido ({{code}})',
    loading: 'Verificando o arquivo…',
    failed: 'Não foi possível verificar este arquivo: {{message}}',
    sev: {
      info: 'Nota',
      warn: 'Aviso',
      error: 'Erro'
    },
    covHeading: 'Cobertura',
    covBar: '{{translatable}} de {{total}} falas vão para a IA',
    covTranslatable: 'Para traduzir: {{count}}',
    covSkipped: 'Nada a traduzir: {{count}}',
    covUntranslated: 'Ficam como estão: {{count}}',
    code: {
      SRT_INDEX_MISSING: 'Número da fala ausente',
      SRT_INDEX_DUPLICATE: 'Número de fala repetido',
      SRT_INDEX_GAP: 'Falha na numeração das falas',
      SRT_TIME_DOT_SEPARATOR: 'Tempo usa ponto em vez de vírgula',
      SRT_TIME_PADDING: 'Tempo sem zeros à esquerda',
      SRT_ARROW_SPACING: 'Espaçamento estranho em volta da seta “-->”',
      SRT_MISSING_BLANK_LINE: 'Falta a linha em branco entre falas',
      SRT_EOL_MIXED: 'Quebras de linha misturadas',
      SRT_TRAILING_WS: 'Espaços no fim da linha',
      SRT_END_BEFORE_START: 'Uma fala termina antes de começar',
      SRT_OVERLAP: 'Uma fala se sobrepõe à seguinte',
      SRT_EMPTY_CUE: 'Fala vazia',
      SRT_TIMING_EXTRAS: 'Parâmetros extras após o tempo',
      SRT_UNPARSEABLE_BLOCK: 'Um bloco não pôde ser lido',
      SRT_AMBIGUOUS_INDEX: 'O número da fala é ambíguo',
      ASS_EOL_MIXED: 'Quebras de linha misturadas',
      ASS_TIME_FORMAT: 'Tempo em formato incomum',
      ASS_END_BEFORE_START: 'Uma fala termina antes de começar',
      ASS_FORMAT_MISSING: 'Falta a linha “Format:”',
      ASS_UNPARSEABLE_EVENT: 'Uma linha de evento não pôde ser lida',
      VTT_HEADER_MISSING: 'Falta o cabeçalho WEBVTT',
      VTT_TIME_FORMAT: 'Tempo em formato incomum',
      VTT_MISSING_BLANK_LINE: 'Falta a linha em branco entre falas',
      VTT_ARROW_IN_TEXT: 'O texto contém “-->”',
      VTT_EOL_MIXED: 'Quebras de linha misturadas',
      VTT_END_BEFORE_START: 'Uma fala termina antes de começar',
      VTT_EMPTY_CUE: 'Fala vazia',
      VTT_UNPARSEABLE_BLOCK: 'Um bloco não pôde ser lido',
      SEG_MARKER_COLLISION: 'O texto parece um marcador interno'
    }
  },
  coverage: {
    heading: 'Não traduzido',
    summaryUntranslated_one: '{{count}} fala vai ficar no idioma original',
    summaryUntranslated_other: '{{count}} falas vão ficar no idioma original',
    summarySkipped_one: '{{count}} fala não tem nada para traduzir',
    summarySkipped_other: '{{count}} falas não têm nada para traduzir',
    showSkipped: 'Mostrar falas sem nada para traduzir',
    hideSkipped: 'Esconder',
    line: 'Linha {{line}}',
    truncated: 'Só as primeiras entradas estão listadas.',
    kept: 'Mantida como está',
    omitted: 'Deixada de fora',
    reason: {
      SKIPPED_NO_LETTERS: 'Sem palavras (notas musicais, pontos ou só formatação)',
      SKIPPED_EMPTY_BLOCK: 'Bloco vazio',
      SKIPPED_DRAWING: 'Desenho vetorial',
      SKIPPED_FRAGMENT: 'Fragmento de karaokê (uma letra/sílaba por evento)',
      UNTRANSLATED_ORPHAN_TEXT: 'Texto sem tempo depois de uma linha em branco',
      UNTRANSLATED_BAD_TIMING: 'Não foi possível ler a linha de tempo',
      UNTRANSLATED_MARKER_COLLISION: 'O texto parece um marcador interno',
      UNTRANSLATED_KARAOKE: 'Linha de karaokê (traduzir quebraria o tempo das sílabas)',
      UNTRANSLATED_BAD_EVENT: 'Linha de diálogo que não pôde ser lida',
      UNTRANSLATED_VALIDATION_FAILED: 'O modelo continuou quebrando a formatação',
      UNTRANSLATED_CONTENT_FILTERED: 'Bloqueada pela moderação do provedor'
    }
  },
  opensubtitles: {
    title: 'Buscar no OpenSubtitles',
    description: 'Legendas para este vídeo, para usar como origem da tradução.',
    needsKey:
      'O OpenSubtitles pede uma API key sua. É gratuita; adicione uma vez e ela fica salva em Configurações.',
    query: 'Título',
    season: 'Temporada',
    episode: 'Episódio',
    languages: 'Idiomas da legenda',
    anyLanguage: 'Qualquer idioma',
    search: 'Buscar',
    searching: 'Buscando…',
    noResults: 'Nenhuma legenda encontrada. Tente outro título ou outro idioma.',
    downloads_one: '{{count}} download',
    downloads_other: '{{count}} downloads',
    badge: {
      hashMatch: 'Feita para este arquivo',
      trusted: 'Uploader confiável',
      hearingImpaired: 'Para surdos',
      machine: 'Tradução automática'
    },
    use: 'Usar',
    useNamed: 'Usar {{name}}',
    downloading: 'Baixando…',
    syncNote:
      'As legendas marcadas com "Feita para este arquivo" batem com este vídeo exato. As outras podem estar fora de sincronia; o BakaSub não ajusta o tempo.',
    quotaNote:
      'Os downloads têm limite diário. Conecte sua conta do OpenSubtitles em Configurações para um limite maior.',
    remaining_one: 'Resta {{count}} download do OpenSubtitles hoje',
    remaining_other: 'Restam {{count}} downloads do OpenSubtitles hoje'
  },
  media: {
    heading: 'O que é isto?',
    title: 'Título',
    searchPlaceholder: 'Buscar pelo título',
    searching: 'Buscando “{{query}}” no TMDB…',
    noResults: 'Nenhum resultado para “{{query}}”. Tente o título original ou acrescente o ano.',
    identifying: 'Procurando “{{query}}” pelo nome do arquivo…',
    notFound:
      '“{{query}}” não foi reconhecido pelo nome do arquivo. Busque o título acima ou continue sem contexto de mídia.',
    offline: 'Não foi possível acessar o TMDB. Você pode continuar sem contexto de mídia.',
    noKey:
      'Adicione uma key gratuita do TMDB para o BakaSub reconhecer o filme ou a série e passar a sinopse para a IA.',
    addKey: 'Adicionar key do TMDB',
    recognized: 'Reconhecido pelo nome do arquivo',
    series: 'Série',
    movie: 'Filme',
    runtime: '{{minutes}} min',
    movieWithEpisode:
      'O nome do arquivo parece de um episódio ({{code}}), mas este título é um filme. Confira se o título certo foi escolhido.',
    season: 'Temporada',
    episode: 'Episódio',
    fromFile: 'Pelo nome do arquivo',
    changed: 'Alterado (o nome do arquivo diz {{code}})',
    notInFile: 'Não está no nome do arquivo',
    absolute:
      'O nome do arquivo usa numeração absoluta (episódio {{n}}). Informe a temporada e o episódio.',
    enterEpisode: 'Informe a temporada e o episódio para carregar a sinopse.',
    noEpisode: 'O TMDB não tem o episódio {{code}} deste título.',
    noEpisodeSynopsis: 'Sem sinopse para {{code}} no TMDB.',
    include: 'Incluir estes detalhes no contexto da tradução',
    attribution: 'Metadados e imagens do TMDB'
  },
  settings: {
    range: 'Use um valor de {{min}} a {{max}}.',
    tabs: {
      general: 'Geral',
      keys: 'Keys',
      cache: 'Cache de tradução',
      tools: 'Ferramentas de vídeo',
      advanced: 'Avançado',
      about: 'Sobre'
    },
    general: {
      language: 'Idioma',
      languageHint: 'Idioma do BakaSub e dos guias.',
      system: 'Igual ao do sistema',
      defaultPreset: 'Preset padrão',
      defaultPresetHint: 'Preset selecionado quando você começa uma tradução.',
      concurrency: 'Jobs ao mesmo tempo',
      concurrencyHint: 'Mais jobs terminam antes, mas dividem os mesmos limites do OpenRouter.',
      autoConfirm: 'Traduzir sem perguntar abaixo de',
      autoConfirmHint:
        'Pula a confirmação de custo quando a estimativa máxima fica abaixo deste valor. Vazio: sempre perguntar.',
      always: 'Sempre perguntar',
      sourceLangs: 'Idiomas de origem preferidos',
      sourceLangsHint:
        'Usados para escolher a faixa de legenda de um vídeo automaticamente, nesta ordem.',
      verbose: 'Logs detalhados',
      verboseHint: 'Escreve mais detalhes no log. Útil para relatar um problema.',
      notifications: 'Notificações do sistema',
      notificationsHint:
        'Avisa quando algo precisa de você com o BakaSub em segundo plano. Nunca com a janela em primeiro plano.',
      notifyJob: 'Jobs',
      notifyJobHint: 'Um arquivo terminou, falhou ou aguarda a confirmação do custo.',
      notifyBatch: 'Lotes',
      notifyBatchHint: 'Uma notificação por lote, quando ele termina ou aguarda confirmação.',
      notifyUpdate: 'Atualizações',
      notifyUpdateHint: 'Uma versão nova do BakaSub está disponível.',
      add: 'Adicionar idioma',
      remove: 'Remover {{name}}',
      up: 'Subir {{name}}'
    },
    keys: {
      openrouter: 'OpenRouter',
      openrouterHint: 'Usada em toda tradução. Conferida com o OpenRouter antes de salvar.',
      tmdb: 'TMDB (opcional)',
      tmdbHint: 'Reconhece filmes e séries para passar a sinopse para a IA.',
      opensubtitles: 'OpenSubtitles (opcional)',
      opensubtitlesHint:
        'Encontra legendas para vídeos que não têm nenhuma, para poder traduzi-las.',
      opensubtitlesLabel: 'API key do OpenSubtitles',
      opensubtitlesAccount: 'Conta (opcional)',
      opensubtitlesAccountHint:
        'Entrar na conta aumenta o limite diário de downloads. A senha fica guardada criptografada neste computador.',
      opensubtitlesAccountHintNoKeyring: 'Entrar na conta aumenta o limite diário de downloads.',
      opensubtitlesNoKeyring:
        'Seu sistema não tem chaveiro, então a senha fica salva sem proteção do sistema. Quem tiver acesso à sua conta de usuário pode lê-la.',
      username: 'Usuário',
      usernameHint: 'Seu nome de usuário do OpenSubtitles, não o e-mail.',
      password: 'Senha',
      connect: 'Entrar',
      connecting: 'Entrando…',
      connectedAs: 'Conectado como {{username}}',
      disconnect: 'Sair'
    },
    cache: {
      body: 'Toda fala traduzida fica guardada, então traduzir o mesmo arquivo de novo (ou um parecido) é instantâneo e gratuito.',
      size: 'Tamanho em disco',
      entries_one: '{{n}} fala',
      entries_other: '{{n}} falas',
      oldest: 'Menos usada em {{date}}',
      empty: 'O cache está vazio.',
      prune: 'Remover falas não usadas há',
      days_one: 'mais de {{count}} dia',
      days_other: 'mais de {{count}} dias',
      all: 'Limpar tudo',
      apply: 'Remover',
      busy: 'Espere os jobs em andamento terminarem antes de limpar o cache.',
      done_one: '{{count}} fala removida.',
      done_other: '{{count}} falas removidas.'
    },
    tools: {
      body: 'O mkvtoolnix é necessário para ler as faixas de legenda dos vídeos e colocar as traduções de volta.',
      check: 'Verificar de novo',
      required: 'precisa da {{version}}'
    },
    advanced: {
      body: 'Como o BakaSub espaça e repete as requisições ao OpenRouter. Os padrões seguem os limites do OpenRouter; mude só se tiver um motivo.',
      quota: 'Modelos gratuitos hoje',
      quotaValue: '{{used}} de {{limit}} requisições usadas',
      quotaUnknown: 'Não foi possível verificar a cota dos modelos gratuitos.',
      tierAuto: 'Detectado: {{tier}}',
      tier: {
        'under-10': 'menos de US$ 10 em créditos comprados (50 requisições/dia)',
        'over-10': 'US$ 10 ou mais comprados (1000 requisições/dia)',
        unknown: 'desconhecido (assumindo 50 requisições/dia)'
      },
      freeTier: 'Patamar dos modelos gratuitos',
      freeTierHint:
        'Detectado pela sua conta do OpenRouter. Defina à mão só se a detecção estiver errada.',
      tierOptions: {
        auto: 'Automático',
        'under-10': 'Menos de US$ 10 comprados',
        'over-10': 'US$ 10 ou mais comprados'
      },
      freeRpm: 'Requisições por minuto (modelos gratuitos)',
      freeRpmHint: 'Diminua se usar a mesma conta do OpenRouter em outros apps ao mesmo tempo.',
      paidRpm: 'Requisições por minuto (modelos pagos)',
      paidRpmHint:
        'Defina um limite se um provedor vive pedindo para diminuir o ritmo. Vazio: sem limite.',
      noLimit: 'Sem limite',
      maxThrottled: 'Tentativas quando pedem para diminuir o ritmo',
      maxTransient: 'Tentativas após erros',
      baseDelay: 'Primeira espera entre tentativas',
      maxDelay: 'Maior espera entre tentativas',
      chunkRounds: 'Rodadas extras para uma parte que falhou',
      chunkRoundsHint:
        'Espera 1 minuto, depois 2, antes de tentar essa parte de novo. 0 falha na hora.',
      seconds: 's',
      parallel: 'Tradução em paralelo',
      parallelBody:
        'O arquivo é dividido em trechos traduzidos ao mesmo tempo. Cada trecho segue em ordem, com o contexto das falas anteriores; os cortes ficam nas pausas entre falas. Não muda o custo.',
      recommended: 'Recomendado: {{min}} a {{max}}',
      perJob: 'Trechos em paralelo por arquivo',
      perJobHint: '1 traduz em sequência, como antes.',
      perJobLow: 'Tradução sequencial: bem mais lenta.',
      perJobHigh: 'Mais cortes no arquivo: pode reduzir a coerência entre trechos.',
      paidMaxInFlight: 'Requisições simultâneas (modelos pagos)',
      paidMaxInFlightHint: 'Somando todos os jobs. Diminui sozinho se o provedor pedir calma.',
      paidMaxInFlightLow: 'Limita o paralelismo quando há mais de um job rodando.',
      paidMaxInFlightHigh:
        'Mais chance de o provedor pedir para diminuir o ritmo e de faltar saldo reservado no OpenRouter.',
      freeMaxInFlight: 'Requisições simultâneas (modelos gratuitos)',
      freeMaxInFlightHint:
        'Somando todos os jobs. O OpenRouter aceita 20 por minuto nos gratuitos.',
      freeMaxInFlightLow: 'Tradução sequencial: bem mais lenta.',
      freeMaxInFlightHigh: 'Vai esbarrar no limite de 20 por minuto e gerar esperas.',
      starved:
        '{{jobs}} jobs × {{perJob}} trechos = {{wanted}} requisições, mas o limite é {{limit}}: alguns trechos vão esperar vaga.',
      restore: 'Restaurar padrões',
      restored: 'Padrões restaurados.'
    },
    about: {
      version: 'Versão {{version}}',
      platform: '{{platform}} {{arch}}, Electron {{electron}}',
      upToDate: 'Você está na versão mais recente.',
      update: 'A versão {{latest}} está disponível.',
      offline: 'Não foi possível verificar atualizações.',
      updates: 'Atualizações',
      checking: 'Verificando atualizações…',
      download: 'Baixar',
      downloading: 'Baixando a versão {{latest}}…',
      downloadProgress: 'Progresso do download',
      cancelDownload: 'Cancelar download',
      ready: 'A versão {{latest}} foi baixada e está pronta para instalar.',
      install: 'Reiniciar e atualizar',
      retry: 'Tentar novamente',
      failed: 'A atualização falhou. Nada foi alterado.',
      unsupported:
        'Esta versão não se atualiza sozinha (build de desenvolvimento ou fora de um AppImage/instalador). Baixe a nova versão na página da release.',
      devBuild:
        'Build de desenvolvimento: a atualização pelo app fica desligada e nunca mexe no BakaSub instalado.',
      openRelease: 'Abrir a página da release',
      showNotes: 'Ver notas da versão',
      hideNotes: 'Ocultar notas da versão',
      autoCheck: 'Verificar atualizações automaticamente',
      autoCheckHint:
        'Procura uma versão nova alguns segundos depois de abrir o BakaSub. Nunca baixa nem instala nada sozinho.',
      dialogTitle: 'Uma nova versão do BakaSub está disponível',
      dialogVersions: 'Versão {{current}} → {{latest}}',
      dialogPublished: 'Publicada em {{date}}',
      dialogSize: 'Tamanho do download: {{size}}',
      dialogNotes: 'Notas da versão',
      dialogHint:
        'Nada é baixado até você decidir. O download e a reinicialização acontecem em Configurações → Sobre.',
      dialogLater: 'Depois',
      dialogUpdate: 'Ver detalhes e atualizar',
      jobsTitle: 'Há jobs em andamento',
      jobsBody:
        'Instalar a atualização reinicia o BakaSub e cancela os jobs em andamento. Eles são mantidos e retomados na próxima abertura do app.',
      jobsBack: 'Voltar',
      jobsConfirm: 'Cancelar jobs e atualizar',
      author: 'Feito por {{name}}',
      checkUpdate: 'Verificar atualização',
      repository: 'Repositório',
      coffee: 'Me pague um café',
      data: 'Pasta de dados',
      tmdb: 'Este produto usa a API do TMDB, mas não é endossado nem certificado pelo TMDB. Os metadados e imagens de filmes e séries vêm do The Movie Database (TMDB).',
      openrouter: 'As traduções são feitas pelos modelos que você escolhe, por meio do OpenRouter.'
    }
  },
  modelsPage: {
    keepOne: 'Mantenha pelo menos um modelo favorito: é com ele que o BakaSub traduz.',
    body: 'Marque com estrela os modelos que devem aparecer ao traduzir e escolha o padrão. Os preços são por 1M tokens, do OpenRouter.',
    refresh: 'Atualizar lista',
    updated: 'Lista de {{date}}',
    missingTitle: 'Não estão mais no OpenRouter',
    missingBody:
      'Estes favoritos saíram do catálogo do OpenRouter. Traduzir com eles vai falhar; remova-os ou escolha outro padrão.',
    remove: 'Remover',
    quota: '{{remaining}} de {{limit}} requisições de modelos gratuitos restantes hoje'
  },
  languagesPage: {
    keepOne: 'Mantenha pelo menos um idioma de destino.',
    body: 'Os idiomas para os quais você traduz. Eles aparecem ao traduzir; o padrão vem selecionado.'
  },
  common: { resize: 'Redimensionar a lista' },
  help: {
    label: 'O que é isto?',
    preset:
      'Um preset reúne como traduzir: instruções de estilo, glossário e temperatura. Os embutidos servem para tipos comuns de mídia; crie os seus em Presets.',
    model:
      'O modelo de IA que traduz. Preço, velocidade e qualidade variam muito; veja “Como escolher um modelo” em Modelos.',
    language: 'O idioma para o qual a legenda é traduzida.',
    notes:
      'Enviadas para a IA junto com cada parte do arquivo. Fatos curtos funcionam melhor: quem são os personagens, como falam, palavras a manter.',
    temperature:
      'Quão livre é a tradução. Baixa (0,1–0,3) fica literal e consistente; perto de 0,5 soa mais natural; acima de 1 começa a improvisar.',
    mediaType:
      'Para que tipo de mídia é este preset. É um rótulo para você e também vai para a IA como contexto.',
    prompt:
      'Instruções de tom e registro (formal, casual, manter piadas, adaptar expressões…). As regras fixas que protegem tempo e formatação são adicionadas pelo BakaSub.',
    glossary:
      'Termos que a IA deve traduzir sempre do mesmo jeito (nomes, lugares, técnicas). Use uma regra para coisas como “manter honoríficos”. Presets prontos não têm glossário: para criar um, duplique um preset existente ou crie um novo.',
    include:
      'Passa para a IA o título, os gêneros e a sinopse (e a do episódio) para nomes e tom saírem certos. Só o título é enviado ao TMDB.',
    concurrency:
      'Quantos jobs rodam ao mesmo tempo. Cada job faz as próprias requisições, então mais jobs atingem os limites do OpenRouter antes.',
    autoConfirm:
      'Quando o custo máximo estimado fica abaixo deste valor, o job começa sem perguntar. A estimativa continua visível na fila.',
    sourceLangs:
      'Quando um vídeo tem várias faixas de legenda, o BakaSub escolhe como origem a primeira desta lista.',
    track:
      'A faixa de legenda que será traduzida. Faixas em imagem (PGS, VobSub) não podem ser traduzidas.'
  },
  presetNames: { movie: 'Filme', series: 'Série', anime: 'Anime', documentary: 'Documentário' },
  mediaTypes: { movie: 'Filme', series: 'Série', anime: 'Anime', documentary: 'Documentário' },
  presets: {
    copyOf: '{{name}} (cópia)',
    copyOfN: '{{name}} (cópia {{n}})',
    body: 'Como a IA deve traduzir: tom, um glossário de nomes e termos, e quanta liberdade ela tem.',
    new: 'Novo preset',
    import: 'Importar',
    exportPlain: 'Exportar',
    export_one: 'Exportar {{count}} preset',
    export_other: 'Exportar {{count}} presets',
    exported_one: '{{count}} preset exportado.',
    exported_other: '{{count}} presets exportados.',
    builtin: 'Embutidos',
    custom: 'Seus presets',
    noCustom: 'Você ainda não tem presets próprios. Crie um ou duplique um embutido.',
    select: 'Selecionar para exportar',
    readOnly: 'Presets embutidos não podem ser alterados. Duplique para criar sua própria versão.',
    copiedFrom: 'Copiado de {{name}}',
    duplicate: 'Duplicar',
    delete: 'Excluir',
    deleteTitle: 'Excluir “{{name}}”?',
    deleteBody: 'Jobs que já o usaram mantêm as configurações.',
    save: 'Salvar preset',
    saved: 'Preset salvo.',
    discard: 'Descartar alterações',
    unsaved: 'Alterações não salvas',
    field: {
      name: 'Nome',
      mediaType: 'Tipo de mídia',
      mediaTypeHint: 'Texto livre, ex.: anime, sitcom, documentário.',
      temperature: 'Temperatura',
      temperatureHint: 'Mais baixa fica mais perto do original; mais alta é mais livre. De 0 a 2.',
      prompt: 'Instruções de estilo',
      promptHint: 'Tom, registro, o que manter ou adaptar. {{count}}/{{max}}',
      promptLanguageTip:
        'Dica: escreva estas instruções em inglês, mesmo traduzindo para outro idioma. Os modelos seguem instruções em inglês com mais precisão; o idioma de destino é definido no job.'
    },
    error: {
      required: 'Obrigatório.',
      tooLong: 'Longo demais.',
      range: 'Use um valor de 0 a 2.'
    },
    glossary: {
      title: 'Glossário',
      hint: 'Nomes e termos que a IA deve seguir. Cada entrada precisa de uma tradução ou de uma regra.',
      source: 'Termo',
      target: 'Tradução',
      rule: 'Regra',
      targetPlaceholder: 'Repita para manter',
      rulePlaceholder: 'ex.: manter o honorífico',
      add: 'Adicionar entrada',
      remove: 'Remover entrada {{n}}',
      empty: 'Nenhuma entrada.',
      tooMany: 'No máximo 500 entradas.',
      rowError: {
        source: 'Informe o termo.',
        targetOrRule: 'Adicione uma tradução ou uma regra.',
        tooLong: 'Longo demais.'
      }
    },
    importDialog: {
      title: 'Importar presets',
      body: 'Escolha um arquivo de presets do BakaSub (.json). Os presets importados sempre entram como novos.',
      conflict: 'Quando o nome já existir',
      rename: 'Importar com outro nome',
      skip: 'Pular',
      choose: 'Escolher arquivo',
      imported_one: '{{count}} preset importado.',
      imported_other: '{{count}} presets importados.',
      skipped: 'Não importados',
      entry: 'Entrada {{n}}',
      nameTaken: 'o nome já existe',
      done: 'Concluir'
    }
  },
  files: {
    videoOnly: 'Escolha um arquivo de vídeo (.mkv).'
  },
  extractPage: {
    body: 'Salve faixas de legenda de um vídeo como arquivos separados, sem traduzir.',
    video: 'Vídeo',
    tracks: 'Faixas para salvar',
    output: 'Pasta de saída',
    start: 'Extrair',
    another: 'Extrair de outro vídeo',
    naming: {
      label: 'Nome dos arquivos',
      player: 'Compatível com player',
      playerHint: 'Episodio.en.ass: os players carregam ao lado do vídeo sozinhos.',
      track: 'Com o id da faixa',
      trackHint: 'Episodio.3.en.ass: sempre único, mas os players não carregam sozinhos.',
      ambiguous_one:
        '{{count}} arquivo tem várias faixas no mesmo idioma: elas recebem um número (Episodio.en.1.ass, Episodio.en.2.ass), que os players podem não carregar sozinhos.',
      ambiguous_other:
        '{{count}} arquivos têm várias faixas no mesmo idioma: elas recebem um número (Episodio.en.1.ass, Episodio.en.2.ass), que os players podem não carregar sozinhos.'
    },
    needTrack: 'Escolha pelo menos uma faixa.'
  },
  extractBatch: {
    title: 'Extrair de vários arquivos',
    steps: { files: 'Arquivos', tracks: 'Faixas', output: 'Saída' },
    files: {
      heading: 'Escolha os vídeos',
      body: 'Uma pasta (com ou sem subpastas) ou vários vídeos (.mkv).',
      none: 'Nenhum vídeo (.mkv) foi encontrado.'
    },
    tracks: {
      heading: 'Escolha as faixas a salvar',
      body: 'Por padrão todas as faixas de texto são salvas (legendas em imagem não podem ser extraídas). Altere um tipo de faixa para todos os arquivos ou ajuste cada arquivo.',
      kinds: 'Tipos de faixa',
      perFile: 'Arquivos',
      count: '{{selected}} de {{total}} faixas'
    },
    output: {
      heading: 'Onde salvar',
      next: 'Sem pasta de saída, cada arquivo é salvo ao lado do seu vídeo.',
      review: '{{tracks}} faixas de {{files}} arquivos serão salvas.'
    },
    error: {
      item: 'O lote não foi criado: corrija {{file}}.',
      fix: 'Ir para as faixas'
    },
    start: 'Extrair {{count}} arquivos',
    starting: 'Iniciando…',
    run: {
      running: 'Extraindo…',
      done: 'Lote concluído.',
      progress: '{{done}} de {{total}} concluídos',
      another: 'Extrair outros arquivos',
      count: {
        queued_one: '{{count}} aguardando',
        queued_other: '{{count}} aguardando',
        running_one: '{{count}} extraindo',
        running_other: '{{count}} extraindo',
        done_one: '{{count}} concluído',
        done_other: '{{count}} concluídos',
        failed_one: '{{count}} com falha ou cancelado',
        failed_other: '{{count}} com falha ou cancelados'
      }
    }
  },
  muxBatch: {
    title: 'Fazer mux de vários vídeos',
    steps: { files: 'Arquivos', pairing: 'Pareamento', output: 'Saída' },
    files: {
      heading: 'Escolha os vídeos e as legendas',
      body: 'Uma pasta (com ou sem subpastas) ou vários arquivos, com os vídeos (.mkv) e suas legendas (.srt, .ass, .ssa, .vtt) juntos. O BakaSub pareia pelo nome ou por temporada e episódio.',
      none: 'Nenhum vídeo (.mkv) encontrado.'
    },
    status: {
      paired_one: '{{count}} vídeo com legendas',
      paired_other: '{{count}} vídeos com legendas',
      unmatchedSubs_one: '{{count}} legenda precisa de um vídeo',
      unmatchedSubs_other: '{{count}} legendas precisam de um vídeo',
      unmatchedVideos_one: '{{count}} vídeo sem legendas',
      unmatchedVideos_other: '{{count}} vídeos sem legendas'
    },
    pairing: {
      heading: 'Confira o pareamento',
      body: 'Cada vídeo recebe as suas legendas. Mude o idioma, o nome e a faixa padrão, mova uma legenda para outro vídeo ou remova-a.',
      unresolved_one: '{{count}} legenda sem vídeo',
      unresolved_other: '{{count}} legendas sem vídeo',
      unresolvedBody:
        'O BakaSub não soube dizer a que vídeo estas pertencem (sem correspondência, ou mais de uma). Escolha o vídeo ou ignore a legenda.',
      assignTo: 'Vídeo de {{name}}',
      moveTo: 'Mover {{name}} para outro vídeo',
      addTo: 'Adicionar uma legenda a {{name}}',
      addSubtitle: 'Adicionar uma legenda…',
      pickVideo: 'Escolha o vídeo',
      ignore: 'Ignorar',
      restore: 'Restaurar',
      videos_one: '{{count}} vídeo será criado',
      videos_other: '{{count}} vídeos serão criados',
      noneIncluded: 'Nenhum vídeo com legendas para criar.',
      count_one: '{{count}} legenda',
      count_other: '{{count}} legendas',
      max: 'Mais de 32 legendas',
      byEpisode: 'Por episódio',
      manual: 'Escolhida por você',
      undetermined: 'Indeterminado',
      out_one: '{{count}} item de fora',
      out_other: '{{count}} itens de fora',
      noSubtitle: 'sem legendas',
      blocked_one: 'Resolva {{count}} legenda para continuar.',
      blocked_other: 'Resolva {{count}} legendas para continuar.'
    },
    output: {
      heading: 'Onde salvar',
      next: 'Sem pasta de saída, cada vídeo novo é salvo ao lado do original.',
      collision_one:
        '{{count}} nome se repete na pasta de saída; um número é acrescentado aos seguintes:',
      collision_other:
        '{{count}} nomes se repetem na pasta de saída; um número é acrescentado aos seguintes:',
      space:
        'Cada vídeo é copiado por inteiro: deixe espaço livre de cerca do tamanho de todos os vídeos.',
      review: '{{videos}} vídeos com {{subs}} legendas serão criados.'
    },
    error: {
      item: 'O lote não foi criado: corrija {{file}}.',
      fix: 'Ir para o pareamento'
    },
    start: 'Criar {{count}} vídeos',
    starting: 'Iniciando…',
    run: {
      running: 'Fazendo mux…',
      done: 'Lote concluído.',
      progress: '{{done}} de {{total}} concluídos',
      another: 'Fazer mux de outros vídeos',
      count: {
        queued_one: '{{count}} na espera',
        queued_other: '{{count}} na espera',
        running_one: '{{count}} em mux',
        running_other: '{{count}} em mux',
        done_one: '{{count}} concluído',
        done_other: '{{count}} concluídos',
        failed_one: '{{count}} falhou ou foi cancelado',
        failed_other: '{{count}} falharam ou foram cancelados'
      }
    }
  },
  muxPage: {
    body: 'Coloque arquivos de legenda em um vídeo. O vídeo é copiado; o original nunca é alterado.',
    video: 'Vídeo',
    subtitles: 'Legendas para adicionar',
    add: 'Adicionar arquivos de legenda',
    empty: 'Nenhuma legenda ainda.',
    language: 'Idioma',
    name: 'Nome da faixa',
    namePlaceholder: 'Por padrão, o nome do idioma',
    isDefault: 'Faixa padrão',
    remove: 'Remover {{name}}',
    guessed: 'Pelo nome do arquivo',
    needLanguage: 'Escolha o idioma.',
    options: 'Opções',
    setDefault: 'Tornar a escolhida a única legenda padrão',
    keepOriginal: 'Manter as legendas que já estão no vídeo',
    start: 'Criar vídeo',
    another: 'Fazer mux de outro vídeo',
    max: 'No máximo 32 legendas.'
  },
  logsPage: {
    body: 'O que o BakaSub fez, com o mais recente embaixo. Keys e outros segredos nunca são escritos aqui.',
    live: 'Ao vivo',
    level: 'Nível',
    levels: {
      debug: 'Tudo',
      info: 'Info e acima',
      warn: 'Avisos e erros',
      error: 'Só erros'
    },
    scope: 'Área',
    allScopes: 'Todas as áreas',
    job: 'Job {{id}}',
    filteredJob: 'Só o job: {{name}}',
    clearJob: 'Mostrar todos os jobs',
    openFolder: 'Abrir pasta de logs',
    export: 'Exportar logs',
    exported: 'Logs exportados.',
    empty: 'Nenhuma linha de log corresponde aos filtros.',
    jobEmpty: 'Este job não tem mais linhas de log. Elas podem ter sido rotacionadas ou apagadas.',
    loading: 'Lendo o log…',
    loadFailed: 'Não foi possível ler o log.',
    details: 'Detalhes',
    count: '{{shown}} de {{total}} linhas'
  }
}
