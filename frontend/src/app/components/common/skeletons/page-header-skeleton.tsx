export default function PageHeaderSkeleton() {
  return (
    <>
      <div className="flex items-center justify-between px-8 py-5">
        <div className="space-y-2">
          <div className="h-5 w-44 animate-pulse rounded bg-paper-dim" />
          <div className="h-3 w-60 animate-pulse rounded bg-paper-dim" />
        </div>
        <div className="h-9 w-9 animate-pulse rounded-full bg-paper-dim" />
      </div>
      <div className="cut-divider" />
    </>
  );
}
