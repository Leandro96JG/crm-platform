'use server';
import { unauthorizedRedirect } from '@/app/libs/auth-redirect';
import { getCurrentUser } from '@/app/libs/session';
import { PrintJob } from '@/app/types/print_job';
import { SearchResponse } from '@/app/types/search_response';
import { redirect } from 'next/navigation';
import PrintQueue from '@/app/components/printing/queue';
import Topbar from '@/app/components/common/topbar';
import { fetchPrintJobs } from '@/app/services/printing';
import { getInitials } from '@/app/utils/user-display';

async function getPrintData(): Promise<{
  printQueue: SearchResponse<PrintJob>;
  cutQueue: SearchResponse<PrintJob>;
}> {
  const printResult = await fetchPrintJobs('print', ['queued', 'printing', 'printed']);
  const cutResult = await fetchPrintJobs('cut', ['queued', 'cutting', 'cut']);

  const printData = printResult.success && printResult.data
    ? printResult.data
    : { result: [], paging: { total: 0, limit: 100, offset: 0 } };

  const cutData = cutResult.success && cutResult.data
    ? cutResult.data
    : { result: [], paging: { total: 0, limit: 100, offset: 0 } };

  if (printResult.unauthorized || cutResult.unauthorized) {
    unauthorizedRedirect();
  }

  return { printQueue: printData, cutQueue: cutData };
}

export default async function PrintingPage() {
  const user = await getCurrentUser();
  if (!user) {
    redirect('/login');
  }

  const { printQueue, cutQueue } = await getPrintData();

  return (
    <>
      <Topbar
        title="Producción"
        subtitle="Cola de impresión y corte"
        initials={getInitials(user.name ?? user.username)}
      />
      <div className="cut-divider" />
      <main className="animate-fadeIn px-8 pb-10">
        <div className="grid gap-5 md:grid-cols-2">
          <PrintQueue
            title="Cola de impresión"
            type="print"
            jobs={printQueue.result}
            paging={printQueue.paging}
          />
          <PrintQueue
            title="Cola de corte"
            type="cut"
            jobs={cutQueue.result}
            paging={cutQueue.paging}
          />
        </div>
      </main>
    </>
  );
}
