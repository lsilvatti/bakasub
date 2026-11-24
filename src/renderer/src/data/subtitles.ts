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
    queryFn: () => unwrap(window.api.subtitles.search(request!)),
    enabled: request !== null,
    staleTime: 5 * 60 * 1000,
    retry: false
  })

export function useSubtitleDownload(): UseMutationResult<SubtitleDownload, Error, number> {
  return useMutation({ mutationFn: (fileId) => unwrap(window.api.subtitles.download(fileId)) })
}
