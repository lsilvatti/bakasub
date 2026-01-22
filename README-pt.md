# 🍜 BakaSub

> *"N-Não é como se eu tivesse feito essa ferramenta de legendas pra você ou algo assim... B-Baka!"*

**BakaSub** é uma ferramenta de tradução de legendas alimentada por IA, ultrarrápida e construída para usuários avançados que exigem **zero dessincronia** e **estética nativa de terminal**. Nascido da frustração com interfaces web desajeitadas e desastres de timing de legendas, BakaSub traz automação de tradução de nível profissional para o seu terminal.

Pense nisso como `btop` encontra `lazygit`, mas para legendas. Sem necessidade de mouse, sem GUI inchada, apenas eficiência pura orientada ao teclado.

## ✨ Recursos

- **🤖 Tradução Alimentada por IA**: Suporte para OpenRouter, Google Gemini, OpenAI ou LLM local
- **⚡ Protocolo Zero Dessincronia**: Contexto de janela deslizante + portões de qualidade garantem sincronização perfeita
- **💾 Cache Inteligente**: Correspondência difusa baseada em SQLite economiza seu dinheiro em traduções repetidas
- **🎨 TUI Neon Nativo**: Interface inspirada em btop que fica *chef's kiss* no seu terminal
- **📦 Binário Primeiro**: Executável único, sem dependências (exceto FFmpeg/MKVToolNix)
- **🔄 Modo Observador**: Solte arquivos em uma pasta, vá embora, deixe o BakaSub cuidar
- **🛠️ Caixa de Ferramentas MKV**: Extrair, muxar, editar cabeçalhos, gerenciar fontes - tudo em um lugar
- **🌍 Trilíngue**: Interface disponível em English, Português (BR) e Español

### Por Que BakaSub?

| 💀 Jeito Antigo | ✨ Jeito BakaSub |
|-----------------|------------------|
| Exportar legendas manualmente | Auto-extrai do MKV |
| Copiar e colar em tradutor web | Chamadas de API em lote com contexto |
| Corrigir dessincronia por 2 horas | Protocolo anti-dessincronia integrado |
| Remuxar manualmente no vídeo | Muxagem em uma etapa com backups |
| Torcer para não ter bagunçado | Portão de qualidade detecta erros |

## 🚀 Instalação

### Instalação Rápida (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/lsilvatti/bakasub/main/install.sh | bash
```

### Instalação Manual

1. **Baixe** o último release para sua plataforma:
   - [Linux (AMD64)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-linux-amd64)
   - [Windows (AMD64)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-windows-amd64.exe)
   - [macOS (Intel)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-amd64)
   - [macOS (Apple Silicon)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-arm64)

2. **Torne executável** (Linux/macOS):
   ```bash
   chmod +x bakasub-*
   sudo mv bakasub-* /usr/local/bin/bakasub
   ```

3. **Verifique a instalação**:
   ```bash
   bakasub --version
   ```

### Dependências

BakaSub precisa dessas ferramentas externas (o assistente oferecerá para baixá-las):

- **FFmpeg**: Processamento de mídia
- **MKVToolNix**: Manipulação de contêiner

## 🎬 Início Rápido

### Primeira Execução (Assistente de Configuração)

No primeiro lançamento, BakaSub te guia por:

1. **Configuração do Provedor de IA**: Escolha seu serviço (OpenRouter recomendado) e insira a chave da API
2. **Verificação de Dependências**: Baixa automaticamente FFmpeg/MKVToolNix se estiverem faltando
3. **Padrões**: Defina seu idioma alvo e modelo preferido

```bash
bakasub
```

### Fluxo Básico: Modo Processo Completo

O caso de uso mais comum - traduzir tudo de uma vez:

1. Inicie o BakaSub
2. Digite o caminho para seu arquivo ou pasta MKV
3. Selecione o modo **"Processo Completo"**
4. Pressione **Enter** para iniciar
5. Pegue um café enquanto o BakaSub faz sua mágica ☕

### Modo Observador (Configure e Esqueça)

Perfeito para automação ou processamento em lote:

1. Crie uma pasta (ex: `~/anime-chegando`)
2. No BakaSub, selecione **"Modo Observador"**
3. Aponte para sua pasta
4. Solte arquivos na pasta
5. BakaSub processa automaticamente novos arquivos conforme aparecem

*Como a pasta de downloads de um adulto responsável, mas que realmente se limpa sozinha.*

## ⌨️ Atalhos de Teclado

### Painel Principal

| Tecla | Ação |
|-------|------|
| `1-4` | Lançar módulos (Extrair, Traduzir, Muxar, Revisar) |
| `5-8` | Abrir caixa de ferramentas (Editor de Cabeçalho, Glossário, etc.) |
| `m` | Mudar modelo de IA |
| `c` | Abrir configuração |
| `q` | Sair |

### Configuração de Trabalho

| Tecla | Ação |
|-------|------|
| `Enter` | Iniciar trabalho |
| `d` | Execução teste (estimativa de custo) |
| `r` | Resolver conflitos de faixa |
| `Esc` | Voltar ao painel |

### Editor de Revisão Manual

| Tecla | Ação |
|-------|------|
| `↑/↓` | Navegar linhas |
| `Enter` | Confirmar edição e próxima |
| `Ctrl+S` | Salvar arquivo |
| `g` | Ir para número de linha |
| `Esc` | Sair do editor |

### Editor de Cabeçalho

| Tecla | Ação |
|-------|------|
| `↑/↓` | Navegar faixas |
| `Space` | Alternar flags (Padrão/Forçado) |
| `Enter` | Aplicar mudanças |
| `Esc` | Cancelar |

## 🎭 Configuração

A config fica em `~/.config/bakasub/config.json`. Configurações principais:

```json
{
  "api_provider": "openrouter",
  "api_key": "sk-or-...",
  "target_lang": "pt-br",
  "remove_hi_tags": true,
  "global_temp": 0.3,
  "touchless_mode": false,
  "prompt_profile": "anime"
}
```

### Perfis de Prompt

BakaSub vem com prompts especializados para diferentes tipos de conteúdo:

- **Anime**: Preserva honoríficos (-san, -kun), mantém nomes de ataques
- **Filme**: Tom formal, expressões idiomáticas localizadas
- **Série**: Estilo equilibrado para conteúdo episódico
- **Documentário**: Precisão técnica sobre criatividade
- **YouTube**: Tom casual, consciente de gírias da internet

Você pode clonar perfis de fábrica e personalizá-los.

## 🛠️ Módulos da Caixa de Ferramentas

### Operações Independentes

1. **Extrair Faixas**: Extrair legendas/áudio do MKV
2. **Traduzir Legenda**: Tradução de IA com suas configurações
3. **Muxar Contêiner**: Combinar faixas em MKV
4. **Revisão Manual**: Editor de visão dividida para correções

### Ferramentas MKVToolNix

5. **Editar Flags/Metadados**: Definir faixas padrão, legendas forçadas
6. **Gerenciar Anexos**: Adicionar/remover fontes do MKV
7. **Adicionar/Remover Faixas**: Remuxador rápido com seleção de faixas
8. **Glossário do Projeto**: Definir termos para tradução consistente

## 🌍 Localização

A interface do BakaSub suporta:

- 🇬🇧 **English** (padrão)
- 🇧🇷 **Português (Brasil)**
- 🇪🇸 **Español**

Mude em `Configuração > Geral > Idioma da Interface`.

## 🐛 Resolução de Problemas

### "Erro de API 401"

Sua chave de API é inválida ou expirou. Execute `bakasub` → `c` (config) → Provedores de IA → reinsira a chave.

### "Conflito de Faixa Detectado"

Múltiplas faixas de legenda correspondem ao seu idioma alvo. BakaSub precisa que você escolha:
- Pressione `r` na Configuração de Trabalho
- Selecione a faixa de **diálogo completo** (geralmente o tamanho de arquivo maior)
- Faixas de Sinais/Músicas são tipicamente menores

### "FFmpeg Não Encontrado"

Instale o FFmpeg:
- **Ubuntu/Debian**: `sudo apt install ffmpeg`
- **macOS**: `brew install ffmpeg`
- **Windows**: Baixe de [ffmpeg.org](https://ffmpeg.org)

Ou deixe o Assistente de Configuração baixá-lo para você.

### Legendas Dessincronizadas

Isso NUNCA deveria acontecer graças ao nosso protocolo anti-dessincronia. Se acontecer:
1. Verifique se você selecionou a faixa de legenda correta (Sinais/Músicas ≠ Diálogo Completo)
2. Verifique se o MKV de origem já não está corrompido (`mkvmerge -i file.mkv`)
3. Abra uma issue no GitHub com as informações do arquivo

## 🤝 Contribuindo

Encontrou um bug? Quer um recurso? Contribuições são bem-vindas!

1. Faça um fork do repositório
2. Crie uma branch de recurso (`git checkout -b recurso-legal`)
3. Faça commit de suas mudanças (`git commit -am 'Adiciona recurso legal'`)
4. Faça push para a branch (`git push origin recurso-legal`)
5. Abra um Pull Request

### Configuração de Desenvolvimento

```bash
git clone https://github.com/lsilvatti/bakasub.git
cd bakasub
go mod download
make build-linux
./bin/bakasub-linux-amd64
```

## 📜 Licença

Licença MIT - veja [LICENSE](LICENSE) para detalhes.

## 💖 Apoio

Gostou do BakaSub? Considere apoiar o desenvolvimento:

- ⭐ Dê uma estrela no repositório
- ☕ [Me pague um café](https://ko-fi.com/lsilvatti) *(aceitamos cafunés também)*
- 📢 Compartilhe com amigos que sofrem com legendas ruins

---

**Feito com 💜 por alguém que assistiu muito anime com legendas terríveis**

*"Você já está... traduzido." - BakaSub, provavelmente*
