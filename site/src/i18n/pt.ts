import type { Dictionary } from './en'

/**
 * Brazilian Portuguese strings, served under `/pt/`. Wording follows `README-pt.md`.
 * `satisfies Dictionary` makes a missing key, or a key `en.ts` does not have, a type error.
 */
export const pt = {
  'meta.title': 'BakaSub: tradutor de legendas por IA para MKV',
  'meta.description':
    'Extraia, traduza com IA e remuxe as legendas dos seus vídeos .mkv, num app desktop, na sua própria máquina.',
  'nav.home': 'Início do BakaSub',
  'nav.docs': 'Docs',
  'nav.docsLabel': 'Docs (em inglês)',
  'nav.github': 'GitHub',
  'nav.githubLabel': 'BakaSub no GitHub',
  'nav.primary': 'Principal',
  'nav.language': 'Idioma',
  'hero.kicker': 'Tradutor de legendas por IA',
  'hero.titleStrong': 'Extraia, traduza com IA e remuxe',
  'hero.titleRest': 'as legendas dos seus vídeos .mkv.',
  'hero.lead':
    'Um app desktop que roda na sua própria máquina. Use sua chave do OpenRouter; seus arquivos nunca saem do seu computador.',
  'download.windows': 'Baixar para Windows (.exe)',
  'download.linux': 'Baixar para Linux (.AppImage)',
  'download.others': 'Outras versões',
  'download.mobile': 'O BakaSub é um app desktop para Windows e Linux: abra esta página no computador para baixá-lo.',
  'download.version': 'Versão {version}',
  'hero.source': 'Ver código-fonte',
  'hero.platforms': 'Windows x64 · Linux x64 · Software livre e de código aberto',
  'screenshots.title': 'Por dentro do app',
  'screenshots.lead': 'Do primeiro arquivo ao remux final, cada etapa fica à sua frente.',
  'carousel.label': 'Capturas de tela do BakaSub',
  'carousel.controls': 'Escolha uma captura de tela',
  'carousel.prev': 'Captura anterior',
  'carousel.next': 'Próxima captura',
  'carousel.slide': '{index} de {total}',
  'carousel.goTo': 'Mostrar captura {index}',
  'carousel.placeholder': 'Captura de tela em breve.',
  'screenshots.wizard.caption': 'Assistente de tradução',
  'screenshots.wizard.alt': 'O assistente de tradução escolhendo uma faixa de legenda de um arquivo .mkv.',
  'screenshots.cost.caption': 'Estimativa de custo',
  'screenshots.cost.alt': 'Tela de preflight com a contagem de tokens, a taxa de acerto do cache e o custo estimado.',
  'screenshots.queue.caption': 'Fila de jobs',
  'screenshots.queue.alt': 'A fila de jobs com traduções em andamento, aguardando e concluídas.',
  'screenshots.tmdb.caption': 'Contexto do TMDB',
  'screenshots.tmdb.alt': 'Contexto de mídia buscado no TMDB, com o backdrop e a sinopse do episódio.',
  'screenshots.presets.caption': 'Presets e glossário',
  'screenshots.presets.alt': 'Um preset de tradução com as regras do glossário.',
  'footer.license': 'Software livre e de código aberto, sob a licença',
  'footer.docs': 'Documentação (em inglês)',
  'footer.releases': 'Todas as versões',
  'footer.coffee': 'Me pague um café',
  'footer.source': 'Código-fonte',
  'features.kicker': 'O que faz',
  'features.title': 'Feito em torno do arquivo de legenda, não só do texto',
  'features.lead':
    'Traduz .srt, .ass, .ssa e .vtt com qualquer modelo do OpenRouter e cuida do MKV em volta deles.',
  'features.cost.title': 'Mostra o custo antes de você pagar',
  'features.cost.body':
    'Contagem de tokens e estimativa de preço com os preços atuais do OpenRouter, mais a lista do que não será traduzido.',
  'features.memory.title': 'Lembra o que já traduziu',
  'features.memory.body':
    'As linhas traduzidas ficam num cache local, então reprocessar um arquivo, ou outro com as mesmas linhas, sai quase de graça.',
  'features.tmdb.title': 'Sabe o que você está assistindo',
  'features.tmdb.body':
    'Lê o nome do release e pode buscar a sinopse e os dados do episódio no TMDB para melhorar a tradução.',
  'features.presets.title': 'Mantém seus termos consistentes',
  'features.presets.body':
    'Presets para Filme, Série, Anime, Documentário ou os seus, e um glossário para nomes, honoríficos e jargões.',
  'features.batch.title': 'Trabalha em lote',
  'features.batch.body': 'Traduza, extraia ou remuxe uma pasta ou temporada inteira de uma vez.',
  'features.tags.title': 'Não mexe na formatação',
  'features.tags.body':
    'Só o texto vai ao modelo: timestamps, estilos e tags voltam exatamente onde estavam.',
  'features.privacy.title': 'Fica no seu computador',
  'features.privacy.body': 'Sem conta, sem telemetria. Sua chave e seus arquivos ficam na sua máquina.',
  'pipeline.kicker': 'Como funciona',
  'pipeline.title': 'Três etapas, um job',
  'pipeline.lead':
    'Use cada etapa separada ou rode tudo de uma vez. Um arquivo de legenda sozinho vai direto para a tradução.',
  'pipeline.extract.title': 'Extrair',
  'pipeline.extract.body': 'Escolha uma faixa de legenda em texto do .mkv e extraia com o mkvtoolnix.',
  'pipeline.translate.title': 'Traduzir',
  'pipeline.translate.body':
    'As linhas de texto vão para o modelo do OpenRouter que você escolheu, com as linhas vizinhas como contexto. Linhas que já estão no cache não são enviadas de novo.',
  'pipeline.remux.title': 'Remuxar',
  'pipeline.remux.body':
    'Remuxe a faixa traduzida no vídeo, mantendo as outras faixas e as fontes anexadas, ou só salve o arquivo de legenda.',
  'pipeline.stepLabel': 'Etapa',
  'privacy.kicker': 'Privacidade',
  'privacy.title': 'Sua máquina, sua chave, seus arquivos',
  'privacy.lead':
    'O BakaSub é um app desktop, não um serviço. Não há cadastro, e nada é enviado ao autor.',
  'privacy.account.title': 'Sem conta',
  'privacy.account.body':
    'Instale e use sua própria chave do OpenRouter. Você paga o OpenRouter diretamente pelo que traduz.',
  'privacy.telemetry.title': 'Sem telemetria',
  'privacy.telemetry.body':
    'Os logs ficam na sua máquina. O arquivo de log exportado esconde sua API key e troca sua pasta pessoal por ~.',
  'privacy.local.title': 'Chave e arquivos ficam locais',
  'privacy.local.body':
    'Sua chave é guardada criptografada no chaveiro do sistema, quando houver um. Configurações, cache e logs ficam na sua pasta de dados do usuário.',
  'privacy.sentTitle': 'O que sai do seu computador',
  'privacy.sent.openrouter':
    'O OpenRouter recebe o texto das legendas a traduzir, suas notas, as instruções do preset e o glossário. Timestamps e tags não são enviados.',
  'privacy.sent.tmdb': 'O TMDB recebe o título extraído do nome do arquivo, só se você adicionar uma chave do TMDB.',
  'privacy.sent.none': 'Nada é enviado até você iniciar um job.',
  'privacy.guideLink': 'Leia as notas completas de privacidade no guia do usuário (em inglês)'
} as const satisfies Dictionary
