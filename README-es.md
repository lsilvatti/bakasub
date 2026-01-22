# 🍜 BakaSub

> *"¡N-No es como si hubiera hecho esta herramienta de subtítulos para ti ni nada... B-Baka!"*

**BakaSub** es una herramienta de traducción de subtítulos ultrarrápida impulsada por IA, construida para usuarios avanzados que exigen **cero desincronización** y **estética nativa de terminal**. Nacido de la frustración con interfaces web torpes y desastres de sincronización de subtítulos, BakaSub trae automatización de traducción de nivel profesional a tu terminal.

Piensa en ello como `btop` se encuentra con `lazygit`, pero para subtítulos. Sin necesidad de mouse, sin GUI hinchada, solo eficiencia pura orientada al teclado.

## ✨ Características

- **🤖 Traducción Impulsada por IA**: Soporte para OpenRouter, Google Gemini, OpenAI o LLM local
- **⚡ Protocolo Cero Desincronización**: Contexto de ventana deslizante + puertas de calidad aseguran sincronización perfecta
- **💾 Caché Inteligente**: Coincidencia difusa basada en SQLite te ahorra dinero en traducciones repetidas
- **🎨 TUI Neón Nativo**: Interfaz inspirada en btop que se ve *chef's kiss* en tu terminal
- **📦 Binario Primero**: Ejecutable único, sin dependencias (excepto FFmpeg/MKVToolNix)
- **🔄 Modo Observador**: Suelta archivos en una carpeta, vete, deja que BakaSub se encargue
- **🛠️ Caja de Herramientas MKV**: Extraer, muxear, editar encabezados, gestionar fuentes - todo en un lugar
- **🌍 Trilingüe**: Interfaz disponible en English, Português (BR) y Español

### ¿Por Qué BakaSub?

| 💀 Forma Antigua | ✨ Forma BakaSub |
|------------------|------------------|
| Exportar subtítulos manualmente | Auto-extrae del MKV |
| Copiar y pegar en traductor web | Llamadas de API por lotes con contexto |
| Corregir desincronización durante 2 horas | Protocolo anti-desincronización integrado |
| Remuxear manualmente en video | Muxeo en un paso con respaldos |
| Esperar no haber arruinado nada | Puerta de calidad detecta errores |

## 🚀 Instalación

### Instalación Rápida (Linux/macOS)

```bash
curl -fsSL https://raw.githubusercontent.com/lsilvatti/bakasub/main/install.sh | bash
```

### Instalación Manual

1. **Descarga** el último release para tu plataforma:
   - [Linux (AMD64)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-linux-amd64)
   - [Windows (AMD64)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-windows-amd64.exe)
   - [macOS (Intel)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-amd64)
   - [macOS (Apple Silicon)](https://github.com/lsilvatti/bakasub/releases/latest/download/bakasub-darwin-arm64)

2. **Hazlo ejecutable** (Linux/macOS):
   ```bash
   chmod +x bakasub-*
   sudo mv bakasub-* /usr/local/bin/bakasub
   ```

3. **Verifica la instalación**:
   ```bash
   bakasub --version
   ```

### Dependencias

BakaSub necesita estas herramientas externas (el asistente ofrecerá descargarlas):

- **FFmpeg**: Procesamiento de medios
- **MKVToolNix**: Manipulación de contenedores

## 🎬 Inicio Rápido

### Primera Ejecución (Asistente de Configuración)

En el primer lanzamiento, BakaSub te guía a través de:

1. **Configuración del Proveedor de IA**: Elige tu servicio (OpenRouter recomendado) e ingresa la clave API
2. **Verificación de Dependencias**: Descarga automáticamente FFmpeg/MKVToolNix si faltan
3. **Valores Predeterminados**: Establece tu idioma objetivo y modelo preferido

```bash
bakasub
```

### Flujo Básico: Modo Proceso Completo

El caso de uso más común - traducir todo de una vez:

1. Inicia BakaSub
2. Ingresa la ruta a tu archivo o carpeta MKV
3. Selecciona el modo **"Proceso Completo"**
4. Presiona **Enter** para iniciar
5. Toma un café mientras BakaSub hace su magia ☕

### Modo Observador (Configúralo y Olvídalo)

Perfecto para automatización o procesamiento por lotes:

1. Crea una carpeta (ej: `~/anime-entrante`)
2. En BakaSub, selecciona **"Modo Observador"**
3. Apúntalo a tu carpeta
4. Suelta archivos en la carpeta
5. BakaSub procesa automáticamente nuevos archivos a medida que aparecen

*Como la carpeta de descargas de un adulto responsable, pero que realmente se limpia sola.*

## ⌨️ Atajos de Teclado

### Panel Principal

| Tecla | Acción |
|-------|--------|
| `1-4` | Lanzar módulos (Extraer, Traducir, Muxear, Revisar) |
| `5-8` | Abrir caja de herramientas (Editor de Encabezado, Glosario, etc.) |
| `m` | Cambiar modelo de IA |
| `c` | Abrir configuración |
| `q` | Salir |

### Configuración de Trabajo

| Tecla | Acción |
|-------|--------|
| `Enter` | Iniciar trabajo |
| `d` | Ejecución de prueba (estimación de costo) |
| `r` | Resolver conflictos de pista |
| `Esc` | Volver al panel |

### Editor de Revisión Manual

| Tecla | Acción |
|-------|--------|
| `↑/↓` | Navegar líneas |
| `Enter` | Confirmar edición y siguiente |
| `Ctrl+S` | Guardar archivo |
| `g` | Ir al número de línea |
| `Esc` | Salir del editor |

### Editor de Encabezado

| Tecla | Acción |
|-------|--------|
| `↑/↓` | Navegar pistas |
| `Space` | Alternar banderas (Predeterminado/Forzado) |
| `Enter` | Aplicar cambios |
| `Esc` | Cancelar |

## 🎭 Configuración

La configuración está en `~/.config/bakasub/config.json`. Ajustes clave:

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

### Perfiles de Prompt

BakaSub viene con prompts especializados para diferentes tipos de contenido:

- **Anime**: Preserva honoríficos (-san, -kun), mantiene nombres de ataques
- **Película**: Tono formal, modismos localizados
- **Serie**: Estilo equilibrado para contenido episódico
- **Documental**: Precisión técnica sobre creatividad
- **YouTube**: Tono casual, consciente de jerga de internet

Puedes clonar perfiles de fábrica y personalizarlos.

## 🛠️ Módulos de la Caja de Herramientas

### Operaciones Independientes

1. **Extraer Pistas**: Extraer subtítulos/audio del MKV
2. **Traducir Subtítulo**: Traducción de IA con tus ajustes
3. **Muxear Contenedor**: Combinar pistas en MKV
4. **Revisión Manual**: Editor de vista dividida para correcciones

### Herramientas MKVToolNix

5. **Editar Banderas/Metadatos**: Establecer pistas predeterminadas, subtítulos forzados
6. **Gestionar Adjuntos**: Agregar/eliminar fuentes del MKV
7. **Agregar/Eliminar Pistas**: Remuxeador rápido con selección de pistas
8. **Glosario del Proyecto**: Definir términos para traducción consistente

## 🌍 Localización

La interfaz de BakaSub soporta:

- 🇬🇧 **English** (predeterminado)
- 🇧🇷 **Português (Brasil)**
- 🇪🇸 **Español**

Cambia en `Configuración > General > Idioma de Interfaz`.

## 🐛 Solución de Problemas

### "Error de API 401"

Tu clave API es inválida o expiró. Ejecuta `bakasub` → `c` (config) → Proveedores de IA → reingresa la clave.

### "Conflicto de Pista Detectado"

Múltiples pistas de subtítulos coinciden con tu idioma objetivo. BakaSub necesita que elijas:
- Presiona `r` en Configuración de Trabajo
- Selecciona la pista de **diálogo completo** (generalmente el tamaño de archivo más grande)
- Las pistas de Señales/Canciones son típicamente más pequeñas

### "FFmpeg No Encontrado"

Instala FFmpeg:
- **Ubuntu/Debian**: `sudo apt install ffmpeg`
- **macOS**: `brew install ffmpeg`
- **Windows**: Descarga desde [ffmpeg.org](https://ffmpeg.org)

O deja que el Asistente de Configuración lo descargue por ti.

### Subtítulos Desincronizados

Esto NUNCA debería suceder gracias a nuestro protocolo anti-desincronización. Si sucede:
1. Verifica que seleccionaste la pista de subtítulos correcta (Señales/Canciones ≠ Diálogo Completo)
2. Verifica que el MKV de origen no esté ya corrupto (`mkvmerge -i file.mkv`)
3. Abre una issue en GitHub con la información del archivo

## 🤝 Contribuyendo

¿Encontraste un bug? ¿Quieres una característica? ¡Las contribuciones son bienvenidas!

1. Haz fork del repositorio
2. Crea una rama de característica (`git checkout -b caracteristica-genial`)
3. Haz commit de tus cambios (`git commit -am 'Agrega característica genial'`)
4. Haz push a la rama (`git push origin caracteristica-genial`)
5. Abre un Pull Request

### Configuración de Desarrollo

```bash
git clone https://github.com/lsilvatti/bakasub.git
cd bakasub
go mod download
make build-linux
./bin/bakasub-linux-amd64
```

## 📜 Licencia

Licencia MIT - ver [LICENSE](LICENSE) para detalles.

## 💖 Apoyo

¿Te gusta BakaSub? Considera apoyar el desarrollo:

- ⭐ Dale una estrella al repositorio
- ☕ [Cómprame un café](https://ko-fi.com/lsilvatti) *(también aceptamos caricias)*
- 📢 Comparte con amigos que sufren del infierno de subtítulos

---

**Hecho con 💜 por alguien que vio demasiado anime con subtítulos terribles**

*"Omae wa mou... traducido." - BakaSub, probablemente*
