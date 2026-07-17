import { PrinterIcon, ScissorsIcon } from '@/app/components/common/icons';
import PrintJobActions from '@/app/components/printing/job-actions';
import { PrintJob } from '@/app/types/print_job';
import { getPrintJobStatusMeta } from '@/app/utils/status';

interface PrintQueueProps {
  title: string;
  type: 'print' | 'cut';
  jobs: PrintJob[];
  paging: { total: number; limit: number; offset: number };
}

export default function PrintQueue({
  title,
  type,
  jobs,
  paging,
}: PrintQueueProps) {
  const isCut = type === 'cut';

  return (
    <div className="overflow-hidden rounded-lg bg-card shadow-card">
      <div className="flex items-center justify-between px-5 pb-3.5 pt-4">
        <h2 className="text-[15px] font-bold text-ink">{title}</h2>
        <span className="text-[11.5px] text-ink-soft">
          {paging.total} trabajos
        </span>
      </div>

      {jobs.length === 0 ? (
        <p className="px-5 pb-5 text-sm text-ink-soft">Sin trabajos pendientes</p>
      ) : (
        <div className="flex flex-col gap-2 px-4 pb-4">
          {jobs.map((job, i) => {
            const meta = getPrintJobStatusMeta(job.status);
            return (
              <div
                key={job.job_id}
                style={{ animationDelay: `${i * 40}ms` }}
                className="flex animate-fadeIn items-center gap-3 rounded-md border border-line p-2.5 transition-colors duration-150 hover:bg-paper-dim/40"
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
                  <div className="flex items-center gap-2 text-[11px] text-ink-soft">
                    {isCut ? 'Cortar' : 'Imprimir'} · {job.copies}x
                    <span
                      className={`whitespace-nowrap rounded-full px-2 py-0.5 text-[10px] font-bold ${meta.badge}`}
                    >
                      {meta.label}
                    </span>
                  </div>
                </div>
                <PrintJobActions jobID={job.job_id} status={job.status} />
              </div>
            );
          })}
        </div>
      )}
    </div>
  );
}
