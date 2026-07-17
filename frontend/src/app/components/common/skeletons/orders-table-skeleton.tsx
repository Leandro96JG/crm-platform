const columnWidths = [
  'w-16',
  'w-32',
  'w-24',
  'w-20',
  'w-14',
];

export default function OrdersTableSkeleton({ rows = 6 }: { rows?: number }) {
  return (
    <table className="w-full border-collapse">
      <thead>
        <tr>
          {['N° de pedido', 'Cliente', 'Origen', 'Estado'].map((h) => (
            <th
              key={h}
              className="border-y border-line bg-paper-dim px-5 py-2.5 text-left text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint"
            >
              {h}
            </th>
          ))}
          <th className="border-y border-line bg-paper-dim px-5 py-2.5 text-right text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint">
            Total
          </th>
        </tr>
      </thead>
      <tbody>
        {Array.from({ length: rows }).map((_, i) => (
          <tr key={i} className="border-b border-line last:border-none">
            {columnWidths.map((w, c) => (
              <td key={c} className="px-5 py-3.5">
                <div
                  className={`h-3.5 animate-pulse rounded bg-paper-dim ${w} ${
                    c === columnWidths.length - 1 ? 'ml-auto' : ''
                  }`}
                />
              </td>
            ))}
          </tr>
        ))}
      </tbody>
    </table>
  );
}
