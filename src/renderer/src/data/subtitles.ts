import {
  useMutation,
  useQuery,
  type UseMutationResult,
  type UseQueryResult
} from '@tanstack/react-query'
import type {
  SubtitleDownload,
  SubtitleSearchRequest,
  SubtitleSearchResponse
} from '@shared/ipc-contract'
import { unwrap } from '@renderer/lib/ipc'
import { qk } from './query-keys'

export const useSubtitleSearch = (
  request: SubtitleSearchRequest | null
): UseQueryResult<SubtitleSearchResponse> =>
  useQuery({
    queryKey: qk.subtitles.search(request),
    // Reading `signal` lets TanStack cancel the query when its last observer leaves (dialog
    // closed, search replaced): the Main then drops the OpenSubtitles requests in flight, and
    // the query reverts to its previous state instead of landing in error.
    queryFn: async ({ signal }) => {
      const searchId = crypto.randomUUID()
      const cancel = (): void => void window.api.subtitles.cancelSearch(searchId)
      signal.addEventListener('abort', cancel, { once: true })
      try {
        return await unwrap(window.api.subtitles.search(request!, searchId))
      } finally {
        signal.removeEventListener('abort', cancel)
      }
    },
    enabled: request !== null,
    staleTime: 5 * 60 * 1000,
    retry: false
  })

export function useSubtitleDownload(): UseMutationResult<SubtitleDownload, Error, number> {
  return useMutation({ mutationFn: (fileId) => unwrap(window.api.subtitles.download(fileId)) })
}
