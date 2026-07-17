interface TableSkeletonProps {
  headers: string[];
  rows?: number;
}

export default function TableSkeleton({ headers, rows = 8 }: TableSkeletonProps) {
  return (
    <div className="overflow-hidden rounded-lg bg-card shadow-card">
      <table className="w-full border-collapse">
        <thead>
          <tr>
            {headers.map((h, i) => (
              <th
                key={h}
                className={`border-y border-line bg-paper-dim px-5 py-2.5 text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint ${
                  i === headers.length - 1 ? 'text-right' : 'text-left'
                }`}
              >
                {h}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {Array.from({ length: rows }).map((_, r) => (
            <tr key={r} className="border-b border-line last:border-none">
              {headers.map((_, c) => {
                const isLast = c === headers.length - 1;
                const width = ['w-20', 'w-32', 'w-24', 'w-16', 'w-28'][
                  c % 5
                ];
                return (
                  <td key={c} className="px-5 py-3.5">
                    <div
                      className={`h-3.5 animate-pulse rounded bg-paper-dim ${width} ${
                        isLast ? 'ml-auto' : ''
                      }`}
                    />
                  </td>
                );
              })}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
