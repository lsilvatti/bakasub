/**
 * Renderer-cancellable IPC requests. The renderer tags a request with an id it generated and may
 * later cancel it on a separate channel. Keys are scoped by sender (`webContents.id`), so a frame
 * can only cancel its own requests, and an id reused while still in flight aborts the older one.
 */
export class AbortRegistry {
  private readonly live = new Map<string, AbortController>()

  /** Runs `work` with a signal the matching `cancel` aborts; the entry is dropped when it settles. */
  async run<T>(
    senderId: number,
    id: string,
    work: (signal: AbortSignal) => Promise<T>
  ): Promise<T> {
    const key = `${senderId}:${id}`
    this.live.get(key)?.abort()
    const ctrl = new AbortController()
    this.live.set(key, ctrl)
    try {
      return await work(ctrl.signal)
    } finally {
      if (this.live.get(key) === ctrl) this.live.delete(key)
    }
  }

  /** `false` when nothing is in flight under that id (already settled, or never started). */
  cancel(senderId: number, id: string): boolean {
    const key = `${senderId}:${id}`
    const ctrl = this.live.get(key)
    if (!ctrl) return false
    this.live.delete(key)
    ctrl.abort()
    return true
  }

  get size(): number {
    return this.live.size
  }
}
