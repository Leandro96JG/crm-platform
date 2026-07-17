export default function KpiCardSkeleton() {
  return (
    <div className="register-mark relative overflow-hidden rounded-lg bg-card px-5 pb-4 pt-5 shadow-card">
      <div className="mb-3.5 flex items-center gap-2.5">
        <div className="h-[34px] w-[34px] flex-none animate-pulse rounded-[10px] bg-paper-dim" />
        <div className="h-3 w-24 animate-pulse rounded bg-paper-dim" />
      </div>
      <div className="mb-2 h-8 w-20 animate-pulse rounded bg-paper-dim" />
      <div className="h-3 w-28 animate-pulse rounded bg-paper-dim" />
    </div>
  );
}
