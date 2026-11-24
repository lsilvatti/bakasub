<h1 align="center">
  <img src="./docs/bakasub-logo-pt.svg" alt="BakaSub" width="560" />
</h1>

🇬🇧 [English](./README.md) · 🇧🇷 **Português**

**Extraia, traduza com IA e remuxe as legendas dos seus vídeos `.mkv`**, num app desktop, na sua própria máquina.

Windows x64 e Linux x64 · [GPL-3.0](./LICENSE) · [Site](https://bakasub.silvatti.com.br/pt/)

<p align="center">
  <a href="https://buymeacoffee.com/leosilvatto"><img src="https://img.shields.io/badge/Me%20pague%20um%20caf%C3%A9-ffdd00?style=for-the-badge&logo=buy-me-a-coffee&logoColor=black" alt="Me pague um café" /></a>
</p>

## O que faz

- **Traduz** `.srt`, `.ass`, `.ssa` e `.vtt` com qualquer modelo do [OpenRouter](https://openrouter.ai). Só o texto vai ao modelo: timestamps, estilos e tags voltam exatamente onde estavam.
- **Cuida de todo o fluxo do MKV**: extrai uma faixa de legenda, traduz e remuxa no vídeo (ou só salva o arquivo de legenda). Use cada etapa separada ou rode tudo de uma vez.
- **Mostra o custo antes de você pagar**: contagem de tokens e estimativa de preço com os preços atuais do OpenRouter, mais a lista do que _não_ será traduzido.
- **Lembra o que já traduziu**: um cache local deixa o reprocessamento de um arquivo quase de graça.
- **Funciona quando o vídeo não tem legenda**: use um arquivo de legenda seu ou busque no [OpenSubtitles](https://www.opensubtitles.com/) pelo app, e traduza como qualquer outra.
- **Sabe o que você está assistindo**: lê o nome do release e pode buscar sinopse e dados do episódio no TMDB para melhorar a tradução.
- **Mantém seus termos consistentes** com presets (Filme, Série, Anime, Documentário ou os seus) e um glossário para nomes, honoríficos e jargões.
- **Trabalha em lote**: traduza, extraia ou remuxe uma pasta ou temporada inteira de uma vez.
- Interface em inglês e português do Brasil. Sem conta, sem telemetria; sua chave e seus arquivos ficam no seu computador.

## Instalação

Baixe a última versão na [página de Releases](https://github.com/lsilvatti/bakasub/releases):

| Sistema | Arquivo                 | Observações                                                                                                                             |
| ------- | ----------------------- | --------------------------------------------------------------------------------------------------------------------------------------- |
| Windows | `bakasub-<v>-setup.exe` | O Windows pode mostrar um aviso do SmartScreen porque o instalador não é assinado. Escolha "Mais informações" → "Executar assim mesmo". |
| Linux   | `bakasub-<v>.AppImage`  | `chmod +x bakasub-<v>.AppImage` e execute. A atualização automática só funciona com o AppImage.                                         |

Você pode conferir o download com o `SHA256SUMS.txt` da mesma release.

## Primeiros passos

Você precisa de uma **API key do OpenRouter**: crie uma em [openrouter.ai/keys](https://openrouter.ai/keys). Alguns modelos são gratuitos (com limite diário); os outros são cobrados por uso.

1. Abra o BakaSub. O guia de primeiro uso pede sua chave do OpenRouter, um modelo favorito e padrão, e os idiomas de destino.
2. Deixe-o baixar o **ffmpeg e o mkvtoolnix** quando pedir (necessários para `.mkv`; só baixa com o seu consentimento). Arquivos de legenda sozinhos não precisam deles.
3. Opcionalmente, adicione uma chave gratuita do [TMDB](https://www.themoviedb.org/settings/api) para enriquecer as traduções com sinopses.
4. Vá em **Traduzir**, solte um vídeo ou uma legenda, escolha um preset, revise o custo e confirme.

Progresso, custo e resultados ficam na **Fila**. O passo a passo completo está no [guia do usuário](./docs/user-guide.md) (em inglês).

## Bom saber

- **Custo e privacidade.** Você paga o OpenRouter diretamente pelo que traduz. O texto das legendas, suas notas e o glossário são enviados ao OpenRouter; o nome do arquivo vai ao TMDB se você adicionar uma chave, e ao OpenSubtitles (com uma impressão digital do arquivo de vídeo) se você adicionar a chave dele e buscar por lá. Nada é enviado sem uma ação sua, e nada vai para o autor. Detalhes no [guia do usuário](./docs/user-guide.md#privacy-what-leaves-your-computer).
- **Seus dados** ficam na pasta de dados do usuário. Atualizar ou desinstalar o app não os apaga. Builds de desenvolvimento (`npm run dev`) usam uma pasta separada e nunca tocam nos dados do app instalado:

  |         | Instalado           | Desenvolvimento         |
  | ------- | ------------------- | ----------------------- |
  | Linux   | `~/.config/bakasub` | `~/.config/bakasub-dev` |
  | Windows | `%APPDATA%\bakasub` | `%APPDATA%\bakasub-dev` |

  Banco, logs, ferramentas baixadas e cache ficam todos dentro dela. `BAKASUB_USER_DATA=<dir>` sobrescreve a pasta (testes, scripts). Detalhes em [Development](./docs/development.md#dev-data-folder) (em inglês).

- **Não suportado**: legendas em imagem (PGS, VobSub, DVB) e macOS.
- **Limitações conhecidas**: não há tela geral do lote depois de sair da página que o criou (os jobs aparecem um a um na Fila); não há override manual de encoding, então um arquivo mal detectado deve ser salvo de novo como UTF-8.

## Documentação

A documentação detalhada está em inglês.

- **Usando o app**: [Guia do usuário](./docs/user-guide.md) · [Reportar um problema](./docs/reporting-issues.md)
- **Contribuindo**: [CONTRIBUTING.md](./CONTRIBUTING.md) · [Desenvolvimento](./docs/development.md) · [Arquitetura](./docs/architecture.md)
- **Internos**: [IPC](./docs/ipc.md) · [Sidecars](./docs/sidecars.md) · [Banco de dados](./docs/database.md) · [Parsing de legendas](./docs/parsing.md) · [Job Memory](./docs/job-memory.md) · [Tradução LLM](./docs/llm.md) · [Jobs](./docs/jobs.md) · [Preflight](./docs/preflight.md) · [Limites e retries](./docs/rate-limits.md) · [Presets](./docs/presets.md) · [Contexto de mídia](./docs/media-context.md) · [Legendas externas](./docs/subtitle-sources.md) · [Logging](./docs/logging.md) · [Renderer](./docs/renderer.md) · [Tokens de marca](./docs/brand-tokens.md) · [Site](./docs/site.md) · [Build e release](./docs/release.md)

## Suporte e contribuição

Achou um bug ou tem uma ideia? [Abra uma issue](https://github.com/lsilvatti/bakasub/issues/new/choose); o [guia](./docs/reporting-issues.md) diz o que incluir. Pull requests são bem-vindos depois de uma conversa numa issue: veja o [CONTRIBUTING.md](./CONTRIBUTING.md).

## Licença

O BakaSub é software livre sob a [GNU GPL v3](./LICENSE). ffmpeg e mkvtoolnix são programas separados, baixados por você e executados como processos externos; mantêm suas próprias licenças. Os dados de filmes e séries vêm do TMDB (atribuição exibida no app); este produto usa a API do TMDB, mas não é endossado nem certificado por ele.
