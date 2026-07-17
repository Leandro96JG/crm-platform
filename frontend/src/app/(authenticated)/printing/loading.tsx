import PageHeaderSkeleton from '@/app/components/common/skeletons/page-header-skeleton';
import ProductionQueueSkeleton from '@/app/components/common/skeletons/production-queue-skeleton';
import PanelCard from '@/app/components/common/panel-card';

export default function Loading() {
  return (
    <>
      <PageHeaderSkeleton />
      <main className="px-8 pb-10">
        <div className="grid gap-5 md:grid-cols-2">
          <PanelCard title="Cola de impresión" bodyPadding>
            <ProductionQueueSkeleton items={6} />
          </PanelCard>
          <PanelCard title="Cola de corte" bodyPadding>
            <ProductionQueueSkeleton items={6} />
          </PanelCard>
        </div>
      </main>
    </>
  );
}
