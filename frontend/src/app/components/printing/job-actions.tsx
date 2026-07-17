'use client';

import { updatePrintJobStatus } from '@/app/services/printing';
import {
  getNextPrintJobStatuses,
  getPrintJobStatusMeta,
} from '@/app/utils/status';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

interface PrintJobActionsProps {
  jobID: string;
  status: string;
}

export default function PrintJobActions({ jobID, status }: PrintJobActionsProps) {
  const router = useRouter();
  const [pending, setPending] = useState<string | null>(null);
  const [error, setError] = useState<string | null>(null);

  const nextStatuses = getNextPrintJobStatuses(status);

  if (nextStatuses.length === 0) {
    return null;
  }

  async function handleClick(newStatus: string) {
    setPending(newStatus);
    setError(null);
    const resp = await updatePrintJobStatus(jobID, newStatus);
    setPending(null);
    if (!resp.success) {
      setError(resp.message);
      return;
    }
    router.refresh();
  }

  return (
    <div className="flex flex-col items-end gap-1">
      <div className="flex flex-wrap justify-end gap-1.5">
        {nextStatuses.map((s) => {
          const meta = getPrintJobStatusMeta(s);
          const isFailed = s === 'failed';
          return (
            <button
              key={s}
              type="button"
              disabled={pending !== null}
              onClick={() => handleClick(s)}
              className={`rounded-full px-2.5 py-1 text-[10.5px] font-bold transition-colors disabled:opacity-50 ${
                isFailed
                  ? 'bg-st-danger-bg text-st-danger-fg hover:brightness-95'
                  : meta.badge + ' hover:brightness-95'
              }`}
            >
              {pending === s ? '…' : meta.label}
            </button>
          );
        })}
      </div>
      {error && (
        <p className="text-[10.5px] text-st-danger-fg">{error}</p>
      )}
    </div>
  );
}
