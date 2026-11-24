// @vitest-environment jsdom
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { cleanup, render, screen } from '@testing-library/react'
import { I18nextProvider } from 'react-i18next'
import { afterEach, describe, expect, it } from 'vitest'
import type { ApiKeyStatus, OpenSubtitlesStatus } from '@shared/ipc-contract'
import { TooltipProvider } from '@renderer/components/ui/popover'
import { i18n } from '@renderer/i18n'
import { installFakeApi, ok } from '@renderer/test/fake-api'
import { KeysSection } from './keys-section'

const NO_KEYRING = /the password is stored without system protection/

function setup(protection: ApiKeyStatus['protection'], os: OpenSubtitlesStatus): void {
  installFakeApi({
    settings: {
      apiKeyStatus: () => ok({ present: true, protection }),
      tmdbKeyStatus: () => ok({ present: false }),
      openSubtitlesStatus: () => ok(os)
    }
  })
  const qc = new QueryClient({ defaultOptions: { queries: { retry: false } } })
  render(
    <I18nextProvider i18n={i18n}>
      <QueryClientProvider client={qc}>
        <TooltipProvider>
          <KeysSection />
        </TooltipProvider>
      </QueryClientProvider>
    </I18nextProvider>
  )
}

afterEach(async () => {
  cleanup()
  await i18n.changeLanguage('en')
})

describe('OpenSubtitles account → keyring warning', () => {
  it('warns before sign-in when the vault has no keyring, and drops the "encrypted" claim', async () => {
    setup('none', { keyPresent: true, username: null })
    expect(await screen.findByText(NO_KEYRING)).toBeTruthy()
    expect(screen.queryByText(/The password is stored encrypted/)).toBeNull()
  })

  it('keeps warning while signed in without a keyring', async () => {
    setup('none', { keyPresent: true, username: 'baka' })
    expect(await screen.findByText('Signed in as baka')).toBeTruthy()
    expect(screen.getByText(NO_KEYRING)).toBeTruthy()
  })

  it('does not warn when the OS keyring protects the vault', async () => {
    setup('os', { keyPresent: true, username: 'baka' })
    expect(await screen.findByText('Signed in as baka')).toBeTruthy()
    expect(screen.getByText(/The password is stored encrypted/)).toBeTruthy()
    expect(screen.queryByText(NO_KEYRING)).toBeNull()
  })

  it('renders in pt-BR', async () => {
    await i18n.changeLanguage('pt-BR')
    setup('none', { keyPresent: true, username: null })
    expect(await screen.findByText(/a senha fica salva sem proteção do sistema/)).toBeTruthy()
  })
})
