# 🍜 BakaSub

> *"¡N-No es como si hubiera hecho esta herramienta de subtítulos para ti ni nada... B-Baka!"*

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/lsilvatti)

**BakaSub** es una herramienta de traducción de subtítulos con IA para usuarios avanzados que exigen **cero desincronización** y **estética nativa de terminal**. Nació de la frustración con interfaces web torpes y desastres de timing.

Piensa en `btop` + `lazygit`, pero para subtítulos. Sin mouse, sin hinchazón—solo eficiencia con teclado.

---

## 📋 Índice

- [Características](#-características)
- [Instalación](#-instalación)
- [Dependencias](#-dependencias)
- [Inicio Rápido](#-inicio-rápido)
- [Guía de Uso](#-guía-de-uso)
- [Configuración](#-configuración)
- [Solución de Problemas](#-solución-de-problemas)
- [Para Desarrolladores](#-para-desarrolladores)
- [Apoyo](#-apoyo)

---

## ✨ Características

| Característica | Qué hace |
|----------------|----------|
| 🤖 **Traducción con IA** | Soporta OpenRouter, Google Gemini, OpenAI y LLMs locales (Ollama/LMStudio) |
| ⚡ **Cero Desinc** | Ventana deslizante + quality gates mantienen timing perfecto |
| 💾 **Caché Inteligente** | Fuzzy matching con SQLite—¿por qué pagar dos veces por la misma línea? |
| 🎨 **TUI Neón** | Una interfaz de terminal tan bonita que olvidarás que las GUIs existen |
| 📦 **Binario Único** | Un archivo, sin Python, sin Node, sin drama |
| 🔄 **Watch Mode** | Suelta archivos en una carpeta, BakaSub se encarga del resto. ¡Magia! ✨ |
| 🛠️ **Toolbox MKV** | Extraer, muxear, editar headers, gestionar fuentes—todo en un lugar |
| 🌍 **Interfaz Trilingüe** | English, Português (BR), Español |

---

## 🚀 Instalación

### Instalación en Una Línea (Linux/macOS)

*"B-Bueno, te lo voy a hacer fácil... ¡pero solo esta vez!"*

```bash
curl -fsSL https://raw.githubusercontent.com/lsilvatti/bakasub/main/install.sh | bash
```

### Descarga Manual

Elige tu plataforma, descarga y listo:

| Plataforma | Link de Descarga |
|------------|------------------|
| 🐧 Linux (AMD64) | [bakasub-linux-amd64](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-linux-amd64) |
| 🪟 Windows (AMD64) | [bakasub-windows-amd64.exe](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-windows-amd64.exe) |
| 🍎 macOS (Intel) | [bakasub-darwin-amd64](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-amd64) |
| 🍎 macOS (Apple Silicon) | [bakasub-darwin-arm64](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-arm64) |

**Setup Linux/macOS:**
```bash
chmod +x bakasub-*
sudo mv bakasub-* /usr/local/bin/bakasub
bakasub --version  # ¡Verifica que funciona!
```

**Windows:** Pon el `.exe` en el PATH o ejecútalo directamente.

---

## 🔧 Dependencias

BakaSub necesita dos herramientas externas. *"¡N-No me mires así! Tienes que instalarlas tú mismo... ¡no es como si pudiera hacer todo por ti!"*

**DEBES instalarlas antes de ejecutar BakaSub:**

| Herramienta | Qué hace | Descarga |
|-------------|----------|----------|
| **FFmpeg** | Procesamiento de medios, extracción de streams | [ffmpeg.org](https://ffmpeg.org/download.html) |
| **MKVToolNix** | Manipulación de containers MKV | [mkvtoolnix.download](https://mkvtoolnix.download/downloads.html) |

### Comandos Rápidos de Instalación

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

**Windows:** Descarga los instaladores de los links arriba, o usa [Chocolatey](https://chocolatey.org/):
```powershell
choco install ffmpeg mkvtoolnix
```

---

## 🎬 Inicio Rápido

### Primera Ejecución

```bash
bakasub
```

En la primera ejecución, un wizard te guía por:

1. **Proveedor de IA** — Elige tu servicio e ingresa la API key
2. **Verificación de Dependencias** — Verifica que FFmpeg y MKVToolNix estén instalados
3. **Predeterminados** — Define idioma objetivo y modelo preferido

*"¡S-Solo te estoy ayudando porque claramente no puedes solo!"*

### Flujo Básico

**Modo Proceso Completo** — El caso de uso más común:

1. Ejecuta `bakasub`
2. Ingresa la ruta al archivo/carpeta MKV
3. Selecciona **Proceso Completo**
4. Presiona **Enter**
5. ☕ Tómate un café. Te lo ganaste.

**Watch Mode** — Configúralo y olvídalo:

1. Crea una carpeta (ej: `~/anime-entrante`)
2. Selecciona **Watch Mode** en BakaSub
3. Apunta a tu carpeta
4. Suelta archivos MKV ahí cuando quieras
5. BakaSub procesa automáticamente nuevos archivos

*Como una carpeta de descargas de adulto responsable que realmente se limpia sola.*

---

## 📖 Guía de Uso

### Teclas del Dashboard

| Tecla | Acción |
|-------|--------|
| `1` | Extraer pistas del MKV |
| `2` | Traducir archivo de subtítulos |
| `3` | Muxear pistas en MKV |
| `4` | Editor de revisión manual |
| `5` | Editar flags/metadatos de pista |
| `6` | Gestionar adjuntos (fuentes) |
| `7` | Remuxeador rápido |
| `8` | Glosario del proyecto |
| `m` | Cambiar modelo de IA |
| `c` | Abrir configuración |
| `q` | Salir |

### Teclas de Configuración de Job

| Tecla | Acción |
|-------|--------|
| `Enter` | Iniciar el job |
| `d` | Dry run (estimación de costo sin llamar API) |
| `r` | Resolver conflictos de pista |
| `Esc` | Volver al dashboard |

### Teclas del Editor de Revisión

| Tecla | Acción |
|-------|--------|
| `↑/↓` | Navegar líneas |
| `Enter` | Confirmar edición, ir a la siguiente |
| `Ctrl+S` | Guardar archivo |
| `g` | Ir a número de línea |
| `Esc` | Salir del editor |

### Módulos del Toolbox

| # | Módulo | Descripción |
|---|--------|-------------|
| 1 | **Extraer Pistas** | Extrae subtítulos o audio del MKV |
| 2 | **Traducir Subtítulo** | Traducción con IA usando tu configuración |
| 3 | **Muxear Container** | Combina pistas en un nuevo MKV |
| 4 | **Revisión Manual** | Editor split-view para correcciones |
| 5 | **Editor de Header** | Define flags de pista predeterminada/forzada |
| 6 | **Adjuntos** | Agrega/elimina fuentes del MKV |
| 7 | **Remuxeador** | Agrega/elimina pistas rápido |
| 8 | **Glosario** | Define términos para traducción consistente entre episodios |

---

## 🎭 Configuración

La config está en `~/.config/bakasub/config.json`

```json
{
  "api_provider": "openrouter",
  "api_key": "sk-or-...",
  "target_lang": "es",
  "remove_hi_tags": true,
  "global_temp": 0.3,
  "touchless_mode": false,
  "prompt_profile": "anime"
}
```

### Perfiles de Prompt

Diferentes contenidos necesitan diferentes estilos de traducción:

| Perfil | Mejor para |
|--------|------------|
| **anime** | Preserva honoríficos (-san, -kun), mantiene nombres de ataques |
| **movie** | Tono formal, expresiones idiomáticas localizadas |
| **series** | Estilo equilibrado para contenido episódico |
| **documentary** | Precisión técnica sobre creatividad |
| **youtube** | Tono casual, consciente de jerga de internet |

Clona perfiles de fábrica para personalizarlos. *"Yo hice los predeterminados, pero puedes cambiarlos... ¡si crees que sabes más!"*

### Idioma de la Interfaz

BakaSub soporta: 🇬🇧 English (predeterminado) · 🇧🇷 Português · 🇪🇸 Español

Cambia en `Configuración > General > Idioma de Interfaz`

---

## 🐛 Solución de Problemas

### "Error de API 401"

Tu API key es inválida o expiró.

→ Presiona `c` → Proveedores de IA → Reingresa tu key

### "Conflicto de Pista Detectado"

Múltiples pistas de subtítulos coinciden con tu idioma. BakaSub necesita que elijas:

→ Presiona `r` en Configuración de Job  
→ Selecciona la pista de **diálogo completo** (generalmente archivo más grande)  
→ Pistas de Signs/Songs son típicamente más pequeñas

### "FFmpeg No Encontrado"

Instala FFmpeg usando los comandos en la sección [Dependencias](#-dependencias) arriba.

*"¡Literalmente te di los comandos... solo cópialos y pégalos! ¡Baka!"*

### Subtítulos Desincronizados

*"¡Esto NUNCA debería pasar. Mi código es perfecto!"* ...pero si pasa:

1. Verifica que seleccionaste la pista correcta (Signs/Songs ≠ Diálogo Completo)
2. Verifica que el MKV fuente no esté corrupto: `mkvmerge -i archivo.mkv`
3. [Abre un issue](https://github.com/lsilvatti/bakasub/issues) con info del archivo

---

## 👨‍💻 Para Desarrolladores

*"Oh, ¿quieres contribuir? Q-Qué osadía..."*

### Compilando desde el Código Fuente

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
make build-all       # Todas las plataformas
make install         # Build + instala en /usr/local/bin
```

### Desarrollo

```bash
make dev    # Ejecuta sin compilar
make test   # Ejecuta tests
make fmt    # Formatea código
make lint   # Ejecuta linter
```

### Contribuyendo

1. Haz fork del repo
2. Crea una rama: `git checkout -b caracteristica-genial`
3. Commit tus cambios: `git commit -am 'Agrega característica genial'`
4. Push: `git push origin caracteristica-genial`
5. Abre un Pull Request

---

## 📜 Licencia

Licencia MIT — Haz lo que quieras, solo no me culpes.

---

## 💖 Apoyo

*"N-No es como si necesitara tu apoyo ni nada... ¡pero si insistes!"*

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/lsilvatti)

- ⭐ Dale una estrella a este repo
- 📢 Comparte con amigos sufriendo con subtítulos malos
- 🐛 Reporta bugs (¡pero sé amable!)

---

**Hecho con 💜 por alguien que vio demasiado anime con subtítulos terribles**

*"Omae wa mou... traducido." — BakaSub, probablemente*
