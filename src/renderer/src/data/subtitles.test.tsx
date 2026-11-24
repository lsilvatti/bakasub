// @vitest-environment jsdom
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { cleanup, renderHook, waitFor } from '@testing-library/react'
import type { ReactNode } from 'react'
import { afterEach, describe, expect, it, vi } from 'vitest'
import type { SubtitleSearchRequest } from '@shared/ipc-contract'
import { installFakeApi, ok } from '@renderer/test/fake-api'
import { qk } from './query-keys'
import { useSubtitleSearch } from './subtitles'

afterEach(cleanup)

const req = (query: string): SubtitleSearchRequest => ({
  videoPath: '/v.mkv',
  query,
  languages: []
})

function setup(): {
  qc: QueryClient
  search: ReturnType<typeof vi.fn>
  cancelSearch: ReturnType<typeof vi.fn>
  wrapper: (p: { children: ReactNode }) => React.JSX.Element
} {
  // The Main answers a cancelled search with CANCELLED; here it simply never answers.
  const search = vi.fn(() => new Promise<never>(() => undefined))
  const cancelSearch = vi.fn(() => ok({ cancelled: true }))
  installFakeApi({ subtitles: { search, cancelSearch } })
  const qc = new QueryClient()
  const wrapper = ({ children }: { children: ReactNode }): React.JSX.Element => (
    <QueryClientProvider client={qc}>{children}</QueryClientProvider>
  )
  return { qc, search, cancelSearch, wrapper }
}

describe('useSubtitleSearch cancellation', () => {
  it('unmounting (dialog closed) cancels the search in the Main and leaves no error behind', async () => {
    const { qc, search, cancelSearch, wrapper } = setup()
    const { unmount } = renderHook(() => useSubtitleSearch(req('a')), { wrapper })
    await waitFor(() => expect(search).toHaveBeenCalledTimes(1))
    const searchId = search.mock.calls[0]![1]
    expect(searchId).toMatch(/^[0-9a-f-]{36}$/)

    unmount()
    await waitFor(() => expect(cancelSearch).toHaveBeenCalledWith(searchId))
    const state = qc.getQueryState(qk.subtitles.search(req('a')))!
    expect(state.status).not.toBe('error')
    expect(state.fetchStatus).toBe('idle')
  })

  it('a new search cancels the previous one only', async () => {
    const { search, cancelSearch, wrapper } = setup()
    const { rerender } = renderHook(({ r }) => useSubtitleSearch(r), {
      wrapper,
      initialProps: { r: req('a') }
    })
    await waitFor(() => expect(search).toHaveBeenCalledTimes(1))
    rerender({ r: req('b') })
    await waitFor(() => expect(search).toHaveBeenCalledTimes(2))
    await waitFor(() => expect(cancelSearch).toHaveBeenCalledTimes(1))
    expect(cancelSearch).toHaveBeenCalledWith(search.mock.calls[0]![1])
    expect(search.mock.calls[1]![1]).not.toBe(search.mock.calls[0]![1])
  })
})
