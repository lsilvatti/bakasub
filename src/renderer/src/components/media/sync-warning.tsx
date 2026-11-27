import { AlertTriangle, Info } from 'lucide-react'
import { useSyncMessages, type SyncCheck } from '@renderer/lib/sync-check'
import { cn } from '@renderer/lib/utils'

/** Non-blocking notice: a frame-rate mismatch is a warning, an unverified hash only a note. */
export function SyncWarning({
  check,
  className
}: {
  check: SyncCheck
  className?: string
}): React.JSX.Element {
  const messages = useSyncMessages()(check, 'long')
  const severe = check.risks.includes('fps-mismatch')
  const Icon = severe ? AlertTriangle : Info
  return (
    <div
      role="status"
      className={cn(
        'flex gap-2 rounded-md border p-3 text-sm',
        severe
          ? 'border-warning/40 bg-warning/10 text-warning'
          : 'bg-muted/40 text-muted-foreground',
        className
      )}
    >
      <Icon className="mt-0.5 size-4 shrink-0" />
      <div className="grid gap-1">
        {messages.map((m) => (
          <p key={m}>{m}</p>
        ))}
      </div>
    </div>
  )
}
