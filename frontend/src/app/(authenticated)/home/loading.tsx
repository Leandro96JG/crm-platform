import KpiCardSkeleton from '@/app/components/common/kpi-card/skeleton';
import OrdersTableSkeleton from '@/app/components/common/skeletons/orders-table-skeleton';
import ProductionQueueSkeleton from '@/app/components/common/skeletons/production-queue-skeleton';
import PanelCard from '@/app/components/common/panel-card';

export default function Loading() {
  return (
    <>
      <div className="flex items-center justify-between px-8 py-5">
        <div className="space-y-2">
          <div className="h-5 w-40 animate-pulse rounded bg-paper-dim" />
          <div className="h-3 w-56 animate-pulse rounded bg-paper-dim" />
        </div>
        <div className="h-9 w-9 animate-pulse rounded-full bg-paper-dim" />
      </div>
      <div className="cut-divider" />

      <main className="flex flex-col gap-6 px-8 pb-10">
        <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          {Array.from({ length: 4 }).map((_, i) => (
            <KpiCardSkeleton key={i} />
          ))}
        </section>

        <section className="grid grid-cols-1 items-start gap-5 xl:grid-cols-[1fr_340px]">
          <PanelCard title="Pedidos recientes">
            <OrdersTableSkeleton />
          </PanelCard>
          <PanelCard title="Cola de producción" bodyPadding>
            <ProductionQueueSkeleton />
          </PanelCard>
        </section>
      </main>
    </>
  );
}
