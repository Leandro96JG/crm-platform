import StatusBadge from '@/app/components/common/status-badge';
import { PencilIcon, WhatsAppIcon } from '@/app/components/common/icons';
import { fetchOrders } from '@/app/services/orders';
import { Order } from '@/app/types/order';
import { formatCurrency, getOrderStatusMeta } from '@/app/utils/status';

const originLabel = (order: Order) => {
  const isWhatsapp = order.source?.toLowerCase().includes('whatsapp');
  if (isWhatsapp) {
    return (
      <span className="inline-flex items-center gap-1.5 text-[11.5px] text-ink-soft">
        <WhatsAppIcon className="h-3 w-3" />
        {order.ai_handled ? 'WhatsApp (IA)' : 'WhatsApp'}
      </span>
    );
  }
  return (
    <span className="inline-flex items-center gap-1.5 text-[11.5px] text-ink-soft">
      <PencilIcon className="h-3 w-3" />
      Manual
    </span>
  );
};

export default async function RecentOrders() {
  const resp = await fetchOrders('', 1, 6);
  const orders = resp.data?.result ?? [];

  if (orders.length === 0) {
    return (
      <p className="px-5 pb-5 text-sm text-ink-soft">No hay pedidos todavía.</p>
    );
  }

  return (
    <table className="w-full border-collapse">
      <thead>
        <tr>
          {['N° de pedido', 'Cliente', 'Origen', 'Estado'].map((h) => (
            <th
              key={h}
              className="border-y border-line bg-paper-dim px-5 py-2.5 text-left text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint"
            >
              {h}
            </th>
          ))}
          <th className="border-y border-line bg-paper-dim px-5 py-2.5 text-right text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint">
            Total
          </th>
        </tr>
      </thead>
      <tbody>
        {orders.map((order, i) => {
          const meta = getOrderStatusMeta(order.status);
          return (
            <tr
              key={order.order_id}
              style={{ animationDelay: `${i * 40}ms` }}
              className="animate-fadeIn border-b border-line transition-colors duration-150 last:border-none hover:bg-paper-dim/40"
            >
              <td className="px-5 py-3 font-mono text-[12.5px] text-ink-soft">
                {order.order_number}
              </td>
              <td className="px-5 py-3 text-[13px] font-semibold text-ink">
                {order.customer_id || '—'}
              </td>
              <td className="px-5 py-3">{originLabel(order)}</td>
              <td className="px-5 py-3">
                <StatusBadge
                  label={meta.label}
                  className={meta.badge}
                  pulse={order.status === 'in_production'}
                />
              </td>
              <td className="px-5 py-3 text-right font-mono text-[13px] font-semibold text-ink">
                {formatCurrency(order.total)}
              </td>
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}
