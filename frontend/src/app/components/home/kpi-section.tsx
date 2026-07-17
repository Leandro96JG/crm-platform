import KpiCard from '@/app/components/common/kpi-card';
import {
  ClockIcon,
  GearIcon,
  GridIcon,
  OrderBoxIcon,
} from '@/app/components/common/icons';
import { fetchOrders } from '@/app/services/orders';
import { fetchPlanchas } from '@/app/services/planchas';
import { fetchPrintJobs } from '@/app/services/printing';

export default async function KpiSection() {
  const [totalResp, pendingResp, planchasResp, printResp, cutResp] =
    await Promise.all([
      fetchOrders('', 1, 1),
      fetchOrders('status=pending', 1, 1),
      fetchPlanchas('is_active=true', 1, 1),
      fetchPrintJobs('print', ['queued', 'printing'], 0, 1),
      fetchPrintJobs('cut', ['queued', 'cutting'], 0, 1),
    ]);

  const cards = [
    {
      label: 'Pedidos totales',
      value: totalResp.data?.paging.total ?? 0,
      tone: 'cut' as const,
      icon: <OrderBoxIcon />,
      delta: 'Historial completo',
      deltaKind: 'flat' as const,
    },
    {
      label: 'Pedidos pendientes',
      value: pendingResp.data?.paging.total ?? 0,
      tone: 'warn' as const,
      icon: <ClockIcon />,
      delta: 'Requieren revisión',
      deltaKind: 'attn' as const,
    },
    {
      label: 'Planchas activas',
      value: planchasResp.data?.paging.total ?? 0,
      tone: 'info' as const,
      icon: <GridIcon />,
      delta: 'En taller ahora',
      deltaKind: 'flat' as const,
    },
    {
      label: 'Trabajos en producción',
      value:
        (printResp.data?.paging.total ?? 0) + (cutResp.data?.paging.total ?? 0),
      tone: 'teal' as const,
      icon: <GearIcon />,
      delta: 'Imprimiendo / cortando',
      deltaKind: 'flat' as const,
    },
  ];

  return (
    <>
      {cards.map((card, i) => (
        <div
          key={card.label}
          className="animate-fadeIn"
          style={{ animationDelay: `${i * 80}ms` }}
        >
          <KpiCard {...card} />
        </div>
      ))}
    </>
  );
}
