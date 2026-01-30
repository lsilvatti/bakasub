# 🍜 BakaSub

> *"N-Não é como se eu tivesse feito essa ferramenta de legendas pra você ou algo assim... B-Baka!"*

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/lsilvatti)

**BakaSub** é uma ferramenta de tradução de legendas com IA para usuários avançados que exigem **zero dessincronia** e **estética nativa de terminal**. Nasceu da frustração com interfaces web desajeitadas e desastres de timing.

Pense em `btop` + `lazygit`, mas pra legendas. Sem mouse, sem inchaço—só eficiência via teclado.

---

## 📋 Índice

- [Recursos](#-recursos)
- [Instalação](#-instalação)
- [Dependências](#-dependências)
- [Início Rápido](#-início-rápido)
- [Guia de Uso](#-guia-de-uso)
- [Configuração](#-configuração)
- [Resolução de Problemas](#-resolução-de-problemas)
- [Para Desenvolvedores](#-para-desenvolvedores)
- [Apoio](#-apoio)

---

## ✨ Recursos

| Recurso | O que faz |
|---------|-----------|
| 🤖 **Tradução com IA** | Suporta OpenRouter, Google Gemini, OpenAI e LLMs locais (Ollama/LMStudio) |
| ⚡ **Zero Dessinc** | Janela deslizante + quality gates mantêm timing perfeito |
| 💾 **Cache Inteligente** | Fuzzy matching com SQLite—por que pagar duas vezes pela mesma linha? |
| 🎨 **TUI Neon** | Interface de terminal tão bonita que você esquece que GUIs existem |
| 📦 **Binário Único** | Um arquivo, sem Python, sem Node, sem drama |
| 🔄 **Watch Mode** | Joga arquivos numa pasta, BakaSub cuida do resto. Mágica! ✨ |
| 🛠️ **Toolbox MKV** | Extrair, muxar, editar headers, gerenciar fontes—tudo num lugar só |
| 🌍 **Interface Trilíngue** | English, Português (BR), Español |

---

## 🚀 Instalação

### Instalação em Uma Linha (Linux/macOS)

*"T-Tá bom, eu vou facilitar pra você... mas só dessa vez!"*

```bash
curl -fsSL https://raw.githubusercontent.com/lsilvatti/bakasub/main/install.sh | bash
```

### Download Manual

Escolha sua plataforma, baixe e pronto:

| Plataforma | Link de Download |
|------------|------------------|
| 🐧 Linux (AMD64) | [bakasub-linux-amd64](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-linux-amd64) |
| 🪟 Windows (AMD64) | [bakasub-windows-amd64.exe](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-windows-amd64.exe) |
| 🍎 macOS (Intel) | [bakasub-darwin-amd64](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-amd64) |
| 🍎 macOS (Apple Silicon) | [bakasub-darwin-arm64](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-arm64) |

**Setup Linux/macOS:**
```bash
chmod +x bakasub-*
sudo mv bakasub-* /usr/local/bin/bakasub
bakasub --version  # Verifica se funcionou!
```

**Windows:** Coloque o `.exe` no PATH ou execute direto.

---

## 🔧 Dependências

BakaSub precisa de duas ferramentas externas. *"N-Não me olhe assim! Você tem que instalar elas você mesmo... não é como se eu pudesse fazer tudo por você!"*

**Você PRECISA instalar antes de rodar o BakaSub:**

| Ferramenta | O que faz | Download |
|------------|-----------|----------|
| **FFmpeg** | Processamento de mídia, extração de streams | [ffmpeg.org](https://ffmpeg.org/download.html) |
| **MKVToolNix** | Manipulação de containers MKV | [mkvtoolnix.download](https://mkvtoolnix.download/downloads.html) |

### Comandos Rápidos de Instalação

**Ubuntu/Debian:**
```bash
sudo apt install ffmpeg mkvtoolnix
```

**Fedora:**
```bash
sudo dnf install ffmpeg mkvtoolnix
```

**Arch Linux:**
```bash
sudo pacman -S ffmpeg mkvtoolnix-cli
```

**macOS (Homebrew):**
```bash
brew install ffmpeg mkvtoolnix
```

**Windows:** Baixe os instaladores nos links acima, ou use [Chocolatey](https://chocolatey.org/):
```powershell
choco install ffmpeg mkvtoolnix
```

---

## 🎬 Início Rápido

### Primeira Execução

```bash
bakasub
```

Na primeira vez, um wizard te guia por:

1. **Provedor de IA** — Escolha seu serviço e insira a API key
2. **Verificação de Dependências** — Verifica se FFmpeg e MKVToolNix estão instalados
3. **Padrões** — Define idioma alvo e modelo preferido

*"E-Eu só tô ajudando porque você claramente não consegue sozinho!"*

### Fluxo Básico

**Modo Processo Completo** — O caso de uso mais comum:

1. Execute `bakasub`
2. Digite o caminho pro arquivo/pasta MKV
3. Selecione **Processo Completo**
4. Aperte **Enter**
5. ☕ Pegue um café. Você mereceu.

**Watch Mode** — Configure e esqueça:

1. Crie uma pasta (ex: `~/anime-chegando`)
2. Selecione **Watch Mode** no BakaSub
3. Aponte pra sua pasta
4. Jogue arquivos MKV lá quando quiser
5. BakaSub processa automaticamente novos arquivos

*Como uma pasta de downloads de adulto responsável que realmente se limpa sozinha.*

---

## 📖 Guia de Uso

### Teclas do Dashboard

| Tecla | Ação |
|-------|------|
| `1` | Extrair faixas do MKV |
| `2` | Traduzir arquivo de legenda |
| `3` | Muxar faixas no MKV |
| `4` | Editor de revisão manual |
| `5` | Editar flags/metadados de faixa |
| `6` | Gerenciar anexos (fontes) |
| `7` | Remuxer rápido |
| `8` | Glossário do projeto |
| `m` | Mudar modelo de IA |
| `c` | Abrir configuração |
| `q` | Sair |

### Teclas de Configuração de Job

| Tecla | Ação |
|-------|------|
| `Enter` | Iniciar o job |
| `d` | Dry run (estimativa de custo sem chamar API) |
| `r` | Resolver conflitos de faixa |
| `Esc` | Voltar ao dashboard |

### Teclas do Editor de Revisão

| Tecla | Ação |
|-------|------|
| `↑/↓` | Navegar linhas |
| `Enter` | Confirmar edição, ir pra próxima |
| `Ctrl+S` | Salvar arquivo |
| `g` | Ir para número de linha |
| `Esc` | Sair do editor |

### Módulos da Toolbox

| # | Módulo | Descrição |
|---|--------|-----------|
| 1 | **Extrair Faixas** | Extrai legendas ou áudio do MKV |
| 2 | **Traduzir Legenda** | Tradução com IA usando suas configurações |
| 3 | **Muxar Container** | Combina faixas num novo MKV |
| 4 | **Revisão Manual** | Editor split-view pra correções |
| 5 | **Editor de Header** | Define flags de faixa padrão/forçada |
| 6 | **Anexos** | Adiciona/remove fontes do MKV |
| 7 | **Remuxer** | Adiciona/remove faixas rápido |
| 8 | **Glossário** | Define termos pra tradução consistente entre episódios |

---

## 🎭 Configuração

A config fica em `~/.config/bakasub/config.json`

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

Conteúdos diferentes precisam de estilos de tradução diferentes:

| Perfil | Melhor pra |
|--------|------------|
| **anime** | Preserva honoríficos (-san, -kun), mantém nomes de ataques |
| **movie** | Tom formal, expressões idiomáticas localizadas |
| **series** | Estilo equilibrado pra conteúdo episódico |
| **documentary** | Precisão técnica sobre criatividade |
| **youtube** | Tom casual, consciente de gírias da internet |

Clone perfis de fábrica pra customizar. *"Eu fiz os padrões, mas você pode mudar... se acha que sabe mais!"*

### Idioma da Interface

BakaSub suporta: 🇬🇧 English (padrão) · 🇧🇷 Português · 🇪🇸 Español

Mude em `Configuração > Geral > Idioma da Interface`

---

## 🐛 Resolução de Problemas

### "Erro de API 401"

Sua API key é inválida ou expirou.

→ Aperte `c` → Provedores de IA → Reinsira sua key

### "Conflito de Faixa Detectado"

Múltiplas faixas de legenda correspondem ao seu idioma. BakaSub precisa que você escolha:

→ Aperte `r` na Configuração de Job  
→ Selecione a faixa de **diálogo completo** (geralmente arquivo maior)  
→ Faixas de Signs/Songs são tipicamente menores

### "FFmpeg Não Encontrado"

Instale o FFmpeg usando os comandos na seção [Dependências](#-dependências) acima.

*"Eu literalmente dei os comandos pra você... só copiar e colar! Baka!"*

### Legendas Dessincronizadas

*"Isso NUNCA deveria acontecer. Meu código é perfeito!"* ...mas se acontecer:

1. Verifique se selecionou a faixa certa (Signs/Songs ≠ Diálogo Completo)
2. Verifique se o MKV de origem não tá corrompido: `mkvmerge -i arquivo.mkv`
3. [Abra uma issue](https://github.com/lsilvatti/bakasub/issues) com info do arquivo

---

## 👨‍💻 Para Desenvolvedores

*"Ah, você quer contribuir? Q-Que ousadia..."*

### Compilando do Código-Fonte

**Requisitos:** Go 1.22+

```bash
git clone https://github.com/lsilvatti/bakasub.git
cd bakasub
go mod download
```

### Comandos de Build

```bash
make build-linux     # Linux AMD64
make build-windows   # Windows AMD64
make build-macos     # macOS Intel + ARM
make build-all       # Todas as plataformas
make install         # Build + instala em /usr/local/bin
```

### Desenvolvimento

```bash
make dev    # Roda sem compilar
make test   # Roda testes
make fmt    # Formata código
make lint   # Roda linter
```

### Contribuindo

1. Faça fork do repo
2. Crie uma branch: `git checkout -b recurso-legal`
3. Commit suas mudanças: `git commit -am 'Adiciona recurso legal'`
4. Push: `git push origin recurso-legal`
5. Abra um Pull Request

---

## 📜 Licença

Licença MIT — Faz o que quiser, só não me culpe.

---

## 💖 Apoio

*"N-Não é como se eu precisasse do seu apoio ou algo assim... mas se você insistir..."*

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/lsilvatti)

- ⭐ Dá uma estrela nesse repo
- 📢 Compartilha com amigos sofrendo com legendas ruins
- 🐛 Reporta bugs (mas seja gentil!)

---

**Feito com 💜 por alguém que assistiu muito anime com legendas terríveis**

*"Omae wa mou... traduzido." — BakaSub, provavelmente*
