import StatusBadge from '@/app/components/common/status-badge';
import OrderStatusActions from '@/app/components/orders/status-actions';
import { PencilIcon, WhatsAppIcon } from '@/app/components/common/icons';
import { Order } from '@/app/types/order';
import { formatCurrency, getOrderStatusMeta } from '@/app/utils/status';

interface OrdersTableProps {
  orders: Order[];
  paging: { total: number; limit: number; offset: number };
}

function OriginTag({ order }: { order: Order }) {
  const isWhatsapp = order.source?.toLowerCase().includes('whatsapp');
  return (
    <span className="inline-flex items-center gap-1.5 text-[11.5px] text-ink-soft">
      {isWhatsapp ? (
        <>
          <WhatsAppIcon className="h-3 w-3" />
          {order.ai_handled ? 'WhatsApp (IA)' : 'WhatsApp'}
        </>
      ) : (
        <>
          <PencilIcon className="h-3 w-3" />
          Manual
        </>
      )}
    </span>
  );
}

export default function OrdersTable({ orders, paging }: OrdersTableProps) {
  if (orders.length === 0) {
    return (
      <div className="rounded-lg bg-card p-10 text-center text-sm text-ink-soft shadow-card">
        No hay pedidos disponibles
      </div>
    );
  }

  return (
    <div className="overflow-hidden rounded-lg bg-card shadow-card">
      <div className="overflow-x-auto">
        <table className="w-full border-collapse">
          <thead>
            <tr>
              {['N° de pedido', 'Cliente', 'Origen', 'Urgencia', 'Estado'].map(
                (h) => (
                  <th
                    key={h}
                    className="border-y border-line bg-paper-dim px-5 py-2.5 text-left text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint"
                  >
                    {h}
                  </th>
                )
              )}
              <th className="border-y border-line bg-paper-dim px-5 py-2.5 text-right text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint">
                Total
              </th>
              <th className="border-y border-line bg-paper-dim px-5 py-2.5 text-right text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint">
                Acciones
              </th>
            </tr>
          </thead>
          <tbody>
            {orders.map((order, i) => {
              const meta = getOrderStatusMeta(order.status);
              const urgent = order.urgency === 'urgent';
              return (
                <tr
                  key={order.order_id}
                  style={{ animationDelay: `${i * 40}ms` }}
                  className="animate-fadeIn border-b border-line transition-colors duration-150 last:border-none hover:bg-paper-dim/40"
                >
                  <td className="px-5 py-3.5 font-mono text-[12.5px] text-ink-soft">
                    {order.order_number}
                  </td>
                  <td className="px-5 py-3.5 text-[13px] font-semibold text-ink">
                    {order.customer_id || '—'}
                  </td>
                  <td className="px-5 py-3.5">
                    <OriginTag order={order} />
                  </td>
                  <td className="px-5 py-3.5">
                    <span
                      className={`text-[12.5px] font-medium ${urgent ? 'text-st-danger-fg' : 'text-ink-soft'}`}
                    >
                      {urgent ? 'Urgente' : 'Normal'}
                    </span>
                  </td>
                  <td className="px-5 py-3.5">
                    <StatusBadge label={meta.label} className={meta.badge} />
                  </td>
                  <td className="px-5 py-3.5 text-right font-mono text-[13px] font-semibold text-ink">
                    {formatCurrency(order.total)}
                  </td>
                  <td className="px-5 py-3.5 text-right">
                    <OrderStatusActions
                      orderID={order.order_id}
                      status={order.status}
                    />
                  </td>
                </tr>
              );
            })}
          </tbody>
        </table>
      </div>
      <div className="border-t border-line px-5 py-3 text-[12px] text-ink-soft">
        Total: {paging.total} pedidos
      </div>
    </div>
  );
}
