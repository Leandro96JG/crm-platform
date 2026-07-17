'use client';

import { updateOrderStatus } from '@/app/services/orders';
import {
  getNextOrderStatuses,
  getOrderStatusMeta,
} from '@/app/utils/status';
import { useRouter } from 'next/navigation';
import { useState } from 'react';

interface OrderStatusActionsProps {
  orderID: string;
  status: string;
}

export default function OrderStatusActions({
  orderID,
  status,
}: OrderStatusActionsProps) {
  const router = useRouter();
  const [open, setOpen] = useState(false);
  const [pending, setPending] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const nextStatuses = getNextOrderStatuses(status);

  if (nextStatuses.length === 0) {
    return <span className="text-[12px] text-ink-faint">—</span>;
  }

  async function handleSelect(newStatus: string) {
    setPending(true);
    setError(null);
    const resp = await updateOrderStatus(orderID, newStatus);
    setPending(false);
    setOpen(false);
    if (!resp.success) {
      setError(resp.message);
      return;
    }
    router.refresh();
  }

  return (
    <div className="relative inline-block text-left">
      <button
        type="button"
        disabled={pending}
        onClick={() => setOpen((v) => !v)}
        className="inline-flex items-center gap-1 rounded-md border border-line bg-card px-2.5 py-1.5 text-[12px] font-medium text-ink-soft transition-colors hover:text-ink disabled:opacity-50"
      >
        {pending ? 'Guardando…' : 'Cambiar estado'}
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth={2}
          className="h-3.5 w-3.5"
        >
          <path d="M6 9l6 6 6-6" strokeLinecap="round" strokeLinejoin="round" />
        </svg>
      </button>

      {open && (
        <>
          <div
            className="fixed inset-0 z-10"
            onClick={() => setOpen(false)}
            aria-hidden
          />
          <div className="absolute right-0 z-20 mt-1 min-w-[160px] overflow-hidden rounded-md border border-line bg-card shadow-card">
            {nextStatuses.map((s) => {
              const meta = getOrderStatusMeta(s);
              return (
                <button
                  key={s}
                  type="button"
                  onClick={() => handleSelect(s)}
                  className="flex w-full items-center gap-2 px-3 py-2 text-left text-[12.5px] text-ink transition-colors hover:bg-paper-dim"
                >
                  <span
                    className={`inline-block rounded-full px-2 py-0.5 text-[10.5px] font-bold ${meta.badge}`}
                  >
                    {meta.label}
                  </span>
                </button>
              );
            })}
          </div>
        </>
      )}

      {error && (
        <p className="absolute right-0 mt-1 whitespace-nowrap text-[11px] text-st-danger-fg">
          {error}
        </p>
      )}
    </div>
  );
}
