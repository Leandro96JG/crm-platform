'use server';
import Topbar from '@/app/components/common/topbar';
import PanelCard from '@/app/components/common/panel-card';
import KpiCardSkeleton from '@/app/components/common/kpi-card/skeleton';
import OrdersTableSkeleton from '@/app/components/common/skeletons/orders-table-skeleton';
import ProductionQueueSkeleton from '@/app/components/common/skeletons/production-queue-skeleton';
import KpiSection from '@/app/components/home/kpi-section';
import RecentOrders from '@/app/components/home/recent-orders';
import ProductionQueue from '@/app/components/home/production-queue';
import { getCurrentUser } from '@/app/libs/session';
import { getInitials } from '@/app/utils/user-display';
import { redirect } from 'next/navigation';
import { Suspense } from 'react';

export default async function Page() {
  const user = await getCurrentUser();
  if (!user) {
    redirect('/login');
  }

  const now = new Date().toLocaleDateString('es-AR', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
  });

  return (
    <>
      <Topbar
        title="Panel general"
        subtitle={`Taller · ${now}`}
        initials={getInitials(user.name ?? user.username)}
      />
      <div className="cut-divider" />

      <main className="flex flex-col gap-6 px-8 pb-10">
        <section className="grid grid-cols-1 gap-4 sm:grid-cols-2 xl:grid-cols-4">
          <Suspense
            fallback={
              <>
                {Array.from({ length: 4 }).map((_, i) => (
                  <KpiCardSkeleton key={i} />
                ))}
              </>
            }
          >
            <KpiSection />
          </Suspense>
        </section>

        <section className="grid grid-cols-1 items-start gap-5 xl:grid-cols-[1fr_340px]">
          <div className="animate-fadeIn">
            <PanelCard
              title="Pedidos recientes"
              action={{ label: 'Ver todos', href: '/orders' }}
            >
              <Suspense fallback={<OrdersTableSkeleton />}>
                <RecentOrders />
              </Suspense>
            </PanelCard>
          </div>

          <div className="animate-fadeIn" style={{ animationDelay: '120ms' }}>
            <PanelCard
              title="Cola de producción"
              subtitle="Trabajos activos"
              action={{ label: 'Ver todo', href: '/printing' }}
              bodyPadding
            >
              <Suspense fallback={<ProductionQueueSkeleton />}>
                <ProductionQueue />
              </Suspense>
            </PanelCard>
          </div>
        </section>
      </main>
    </>
  );
}
