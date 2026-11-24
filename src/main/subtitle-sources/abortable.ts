import { AppError } from '../errors'

export const cancelled = (what: string): AppError => new AppError('CANCELLED', `${what} cancelled`)

/**
 * Settles with `shared`, or rejects with `CANCELLED` as soon as `signal` aborts, whichever first.
 * `shared` itself keeps running: use it when the work is shared with other callers, so one
 * caller's abort only detaches that caller.
 */
export function detach<T>(shared: Promise<T>, signal: AbortSignal, what: string): Promise<T> {
  if (signal.aborted) return Promise.reject(cancelled(what))
  return new Promise<T>((resolve, reject) => {
    const onAbort = (): void => reject(cancelled(what))
    signal.addEventListener('abort', onAbort, { once: true })
    shared.then(resolve, reject).finally(() => signal.removeEventListener('abort', onAbort))
  })
}
