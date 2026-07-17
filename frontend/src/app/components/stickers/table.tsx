import StatusBadge from '@/app/components/common/status-badge';
import { Plancha } from '@/app/types/plancha';
import Image from 'next/image';

interface StickersTableProps {
  planchas: Plancha[];
  paging: { total: number; limit: number; offset: number };
}

export default function StickersTable({ planchas, paging }: StickersTableProps) {
  if (planchas.length === 0) {
    return (
      <div className="rounded-lg bg-card p-10 text-center text-sm text-ink-soft shadow-card">
        No hay planchas disponibles
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-lg bg-card shadow-card">
      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr>
              {['Plancha', 'Categoría', 'Subcategoría', 'Estado'].map((h) => (
                <th
                  key={h}
                  className="border-y border-line bg-paper-dim px-5 py-2.5 text-left text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint"
                >
                  {h}
                </th>
              ))}
            </tr>
          </thead>
          <tbody>
            {planchas.map((plancha, i) => (
              <tr
                key={plancha.plancha_id}
                style={{ animationDelay: `${i * 40}ms` }}
                className="animate-fadeIn border-b border-line transition-colors duration-150 last:border-none hover:bg-paper-dim/40"
              >
                <td className="px-5 py-3">
                  <div className="flex items-center gap-3">
                    <div className="relative h-11 w-11 flex-none overflow-hidden rounded-md border border-line bg-paper-dim">
                      {plancha.preview_image_url ? (
                        <Image
                          src={plancha.preview_image_url}
                          alt={plancha.name}
                          fill
                          sizes="44px"
                          className="object-cover"
                          unoptimized
                        />
                      ) : (
                        <div className="flex h-full w-full items-center justify-center text-[10px] text-ink-faint">
                          s/img
                        </div>
                      )}
                    </div>
                    <div className="min-w-0">
                      <div className="truncate text-[13px] font-semibold text-ink">
                        {plancha.name}
                      </div>
                      <div className="truncate text-[11.5px] text-ink-soft">
                        {plancha.description}
                      </div>
                    </div>
                  </div>
                </td>
                <td className="px-5 py-3 text-[13px] text-ink">
                  {plancha.category || '—'}
                </td>
                <td className="px-5 py-3 text-[13px] text-ink-soft">
                  {plancha.subcategory || '—'}
                </td>
                <td className="px-5 py-3">
                  {plancha.is_active ? (
                    <StatusBadge
                      label="Activa"
                      className="bg-st-ok-bg text-st-ok-fg"
                    />
                  ) : (
                    <StatusBadge
                      label="Inactiva"
                      className="bg-paper-dim text-ink-soft"
                    />
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
      <div className="border-t border-line px-5 py-3 text-[12px] text-ink-soft">
        Total: {paging.total} planchas
      </div>
    </div>
  );
}
