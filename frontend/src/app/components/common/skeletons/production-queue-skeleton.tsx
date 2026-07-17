export default function ProductionQueueSkeleton({
  items = 5,
}: {
  items?: number;
}) {
  return (
    <div className="flex flex-col gap-2">
      {Array.from({ length: items }).map((_, i) => (
        <div
          key={i}
          className="flex items-center gap-3 rounded-md border border-line p-2.5"
        >
          <div className="h-[26px] w-[26px] flex-none animate-pulse rounded-sm bg-paper-dim" />
          <div className="h-[30px] w-[30px] flex-none animate-pulse rounded-[9px] bg-paper-dim" />
          <div className="min-w-0 flex-1 space-y-1.5">
            <div className="h-3 w-2/3 animate-pulse rounded bg-paper-dim" />
            <div className="h-2.5 w-16 animate-pulse rounded bg-paper-dim" />
          </div>
          <div className="h-4 w-16 flex-none animate-pulse rounded-full bg-paper-dim" />
        </div>
      ))}
    </div>
  );
}
