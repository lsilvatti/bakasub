import type { LocaleShape } from '../../locale-shape'
import type { errors as en } from '../en/errors'

export const errors: LocaleShape<typeof en> = {
  NETWORK: 'Não foi possível acessar o servidor. Verifique sua conexão e tente de novo.',
  INVALID_MANIFEST: 'A lista de download das ferramentas é inválida. Tente mais tarde.',
  CHECKSUM_MISMATCH: 'O arquivo baixado não passou na verificação de integridade e foi descartado.',
  EXTRACT_FAILED: 'Não foi possível descompactar a ferramenta baixada.',
  PROBE_FAILED: 'A ferramenta instalada não iniciou. Reinstale em Configurações.',
  PROCESS_TIMEOUT: 'Uma ferramenta externa demorou demais e foi interrompida.',
  UNSUPPORTED_PLATFORM: 'Esta ferramenta não está disponível para o seu sistema.',
  SIDECAR_BUSY: 'As ferramentas já estão sendo instaladas.',
  SIDECAR_MISSING: 'O mkvtoolnix não está instalado. Instale para trabalhar com vídeos.',
  UNKNOWN_TOOL: 'Ferramenta desconhecida.',
  CANCELLED: 'Cancelado.',
  MEMORY_BUSY: 'O cache de tradução está ocupado. Tente quando os jobs em andamento terminarem.',
  LLM_AUTH: 'O OpenRouter recusou a API key. Confira em Configurações.',
  LLM_CREDITS: 'Seus créditos do OpenRouter acabaram. Adicione créditos e tente o job de novo.',
  LLM_BAD_REQUEST: 'O modelo recusou a requisição. Tente outro modelo.',
  LLM_UNAVAILABLE:
    'O modelo está indisponível no momento. Tente mais tarde ou escolha outro modelo.',
  LLM_DAILY_LIMIT:
    'Limite diário dos modelos gratuitos atingido. Ele zera às 00:00 UTC; o que já foi traduzido fica salvo, então tentar depois continua de onde parou.',
  LLM_UNKNOWN_MODEL: 'Este modelo não está mais disponível no OpenRouter. Escolha outro.',
  LLM_TOO_MANY_FAILURES: 'Muitas falas não puderam ser traduzidas. Tente outro modelo ou preset.',
  LLM_NO_ENDPOINT:
    'Nenhum provedor do OpenRouter aceitou este modelo para a sua conta. Confira as restrições de provedor nas configurações de privacidade do OpenRouter (provedores permitidos/ignorados, ZDR) ou escolha outro modelo.',
  LLM_REASONING_BUDGET:
    'O modelo gastou todo o orçamento de saída em raciocínio e não devolveu a tradução. Escolha outro modelo.',
  INVALID_PATH: 'Este arquivo ou pasta não pode ser usado. Escolha de novo.',
  DISK_FULL: 'Não há espaço livre em disco suficiente para a saída.',
  PRESET_NOT_FOUND: 'Este preset não existe mais.',
  PRESET_IN_USE: 'Este preset está sendo usado por um job na fila ou em andamento.',
  PRESET_NAME_TAKEN: 'Já existe um preset com este nome.',
  PRESET_READ_ONLY: 'Presets embutidos não podem ser alterados. Duplique para editar uma cópia.',
  MODEL_NOT_SET: 'Escolha um modelo primeiro: marque um em Modelos e defina como padrão.',
  INVALID_PRESET_FILE: 'Este não é um arquivo de presets do BakaSub válido.',
  INVALID_INPUT: 'Parte das informações enviadas era inválida.',
  JOB_NOT_FOUND: 'Este job não existe mais.',
  JOB_INVALID_STATE: 'Esta ação não está disponível para o job agora.',
  BATCH_NOT_FOUND: 'Este lote não existe mais.',
  MEDIA_TOOL_FAILED: 'O mkvtoolnix não conseguiu processar este arquivo.',
  INVALID_SUBTITLE: 'Não foi possível ler este arquivo de legenda.',
  SETTINGS_CONFLICT:
    'O modelo e o idioma padrão precisam continuar nas suas listas. Escolha outro padrão antes de remover.',
  TARGET_LANG_NOT_SET:
    'Escolha um idioma de destino primeiro: adicione um em Idiomas e defina como padrão.',
  TMDB_AUTH: 'O TMDB recusou a key. Confira em Configurações.',
  TMDB_KEY_MISSING: 'Adicione uma key do TMDB em Configurações para buscar títulos. É opcional.',
  OPENSUBTITLES_AUTH:
    'O OpenSubtitles recusou a API key ou a conta. Confira em Configurações (a conta usa o nome de usuário, não o e-mail).',
  OPENSUBTITLES_KEY_MISSING:
    'Adicione uma API key do OpenSubtitles em Configurações para buscar legendas.',
  OPENSUBTITLES_QUOTA:
    'O limite de downloads do OpenSubtitles de hoje foi atingido. Tente amanhã ou entre na sua conta em Configurações para um limite maior.',
  URL_NOT_ALLOWED: 'Este link não pode ser aberto pelo BakaSub.',
  UPDATE_UNSUPPORTED: 'Atualizações só são instaladas no AppImage ou no instalador do Windows.',
  UPDATE_NOT_READY: 'Nenhuma atualização está pronta para instalar.',
  UPDATE_JOBS_ACTIVE: 'Há jobs em execução. Confirme para cancelá-los e instalar a atualização.',
  FORBIDDEN: 'A requisição foi bloqueada.',
  INTERNAL: 'Algo deu errado dentro do BakaSub. Os detalhes estão nos logs.',
  unknown: 'Algo deu errado ({{code}}). Os detalhes estão nos logs.'
}
