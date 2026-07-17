import { PrinterIcon, ScissorsIcon } from '@/app/components/common/icons';
import { fetchPrintJobs } from '@/app/services/printing';
import { PrintJob } from '@/app/types/print_job';
import { getPrintJobStatusMeta } from '@/app/utils/status';

const activePulseStatuses = new Set(['printing', 'cutting']);

export default async function ProductionQueue() {
  const [printResp, cutResp] = await Promise.all([
    fetchPrintJobs('print', ['queued', 'printing'], 0, 5),
    fetchPrintJobs('cut', ['queued', 'cutting'], 0, 5),
  ]);

  const queue: PrintJob[] = [
    ...(printResp.data?.result ?? []),
    ...(cutResp.data?.result ?? []),
  ]
    .sort((a, b) => a.queue_position - b.queue_position)
    .slice(0, 6);

  if (queue.length === 0) {
    return <p className="text-sm text-ink-soft">Sin trabajos activos.</p>;
  }

  return (
    <div className="flex flex-col gap-2">
      {queue.map((job) => {
        const meta = getPrintJobStatusMeta(job.status);
        const isCut = job.job_type === 'cut';
        return (
          <div
            key={job.job_id}
            className="flex items-center gap-3 rounded-md border border-line p-2.5 transition-colors duration-150 hover:bg-paper-dim/40"
          >
            <div className="flex h-[26px] w-[26px] flex-none items-center justify-center rounded-sm bg-paper-dim font-mono text-[11.5px] font-bold text-ink-soft">
              {job.queue_position || '—'}
            </div>
            <div
              className={`flex h-[30px] w-[30px] flex-none items-center justify-center rounded-[9px] [&>svg]:h-3.5 [&>svg]:w-3.5 ${
                isCut
                  ? 'bg-st-prod-bg text-st-prod-fg'
                  : 'bg-st-danger-bg text-cut-dark'
              }`}
            >
              {isCut ? <ScissorsIcon /> : <PrinterIcon />}
            </div>
            <div className="min-w-0 flex-1">
              <div className="truncate text-[12.5px] font-semibold text-ink">
                {job.notes || `Trabajo ${job.job_id.slice(0, 8)}`}
              </div>
              <div className="text-[11px] text-ink-soft">
                {isCut ? 'Cortar' : 'Imprimir'}
              </div>
            </div>
            <span
              className={`flex-none whitespace-nowrap rounded-full px-2 py-0.5 text-[10.5px] font-bold ${meta.badge} ${
                activePulseStatuses.has(job.status) ? 'animate-pulse' : ''
              }`}
            >
              {meta.label}
            </span>
          </div>
        );
      })}
    </div>
  );
}
