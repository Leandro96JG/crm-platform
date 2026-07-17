'use server';
import Topbar from '@/app/components/common/topbar';
import StatusBadge from '@/app/components/common/status-badge';
import OrderStatusActions from '@/app/components/orders/status-actions';
import { PencilIcon, WhatsAppIcon } from '@/app/components/common/icons';
import { getCurrentUser } from '@/app/libs/session';
import { fetchOrder } from '@/app/services/orders';
import { fetchPlanchas, fetchMaterials } from '@/app/services/planchas';
import {
  formatCurrency,
  formatDateTime,
  getOrderStatusMeta,
} from '@/app/utils/status';
import { getInitials } from '@/app/utils/user-display';
import { notFound, redirect } from 'next/navigation';
import Link from 'next/link';

type Params = { params: Promise<{ order_id: string }> };

export default async function OrderDetailPage({ params }: Params) {
  const user = await getCurrentUser();
  if (!user) {
    redirect('/login');
  }

  const { order_id } = await params;

  const [orderResp, planchasResp, materialsResp] = await Promise.all([
    fetchOrder(order_id),
    fetchPlanchas('', 1, 100),
    fetchMaterials(),
  ]);

  if (orderResp.unauthorized) {
    redirect('/login');
  }
  if (!orderResp.success || !orderResp.data) {
    notFound();
  }

  const order = orderResp.data;
  const planchaNames = new Map(
    (planchasResp.data?.result ?? []).map((p) => [p.plancha_id, p.name])
  );
  const materialNames = new Map(
    (materialsResp.data?.result ?? []).map((m) => [m.material_id, m.name])
  );

  const statusMeta = getOrderStatusMeta(order.status);
  const isWhatsapp = order.source?.toLowerCase().includes('whatsapp');
  const activeStatus = order.status === 'in_production';

  return (
    <>
      <Topbar
        title={`Pedido ${order.order_number}`}
        subtitle={order.customer_id}
        initials={getInitials(user.name ?? user.username)}
      />
      <div className="cut-divider" />
      <main className="animate-fadeIn px-8 pb-10">
        <div className="mb-4 flex items-center justify-between">
          <Link
            href="/orders"
            className="text-[13px] font-semibold text-cut-dark hover:underline"
          >
            ← Volver a pedidos
          </Link>
          <OrderStatusActions orderID={order.order_id} status={order.status} />
        </div>

        <div className="grid grid-cols-1 gap-6 lg:grid-cols-3">
          <section className="lg:col-span-2 rounded-lg bg-card p-5 shadow-card">
            <div className="mb-4 flex items-center justify-between">
              <h2 className="text-[15px] font-bold text-ink">Ítems</h2>
              <span className="text-[12px] text-ink-soft">
                {order.items?.length ?? 0} ítem
                {(order.items?.length ?? 0) === 1 ? '' : 's'}
              </span>
            </div>
            <div className="overflow-x-auto">
              <table className="w-full border-collapse">
                <thead>
                  <tr className="border-b border-line text-left text-[11.5px] uppercase tracking-wide text-ink-faint">
                    <th className="py-2 pr-3 font-semibold">Plancha</th>
                    <th className="py-2 pr-3 font-semibold">Material</th>
                    <th className="py-2 pr-3 text-right font-semibold">Cant.</th>
                    <th className="py-2 pr-3 text-right font-semibold">
                      Precio unit.
                    </th>
                    <th className="py-2 text-right font-semibold">Subtotal</th>
                  </tr>
                </thead>
                <tbody>
                  {(order.items ?? []).map((item) => (
                    <tr
                      key={item.item_id}
                      className="border-b border-line/60 text-[13px] text-ink"
                    >
                      <td className="py-2.5 pr-3">
                        {planchaNames.get(item.plancha_id) ?? item.plancha_id}
                        {item.custom_design_notes && (
                          <span className="block text-[11px] text-ink-faint">
                            {item.custom_design_notes}
                          </span>
                        )}
                      </td>
                      <td className="py-2.5 pr-3">
                        {materialNames.get(item.material_id) ??
                          item.material_id}
                      </td>
                      <td className="py-2.5 pr-3 text-right font-mono">
                        {item.sheet_quantity}
                      </td>
                      <td className="py-2.5 pr-3 text-right font-mono">
                        {formatCurrency(item.unit_price)}
                      </td>
                      <td className="py-2.5 text-right font-mono font-semibold">
                        {formatCurrency(item.subtotal)}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
            <div className="mt-4 flex items-center justify-end gap-3 border-t border-line pt-4">
              <span className="text-[13px] text-ink-soft">Total:</span>
              <span className="font-mono text-xl font-bold text-ink">
                {formatCurrency(order.total)}
              </span>
            </div>
          </section>

          <aside className="flex flex-col gap-6">
            <section className="rounded-lg bg-card p-5 shadow-card">
              <h2 className="mb-4 text-[15px] font-bold text-ink">
                Información
              </h2>
              <dl className="flex flex-col gap-3 text-[13px]">
                <div className="flex items-center justify-between">
                  <dt className="text-ink-soft">Estado</dt>
                  <dd>
                    <StatusBadge
                      label={statusMeta.label}
                      className={statusMeta.badge}
                      pulse={activeStatus}
                    />
                  </dd>
                </div>
                <div className="flex items-center justify-between">
                  <dt className="text-ink-soft">Cliente</dt>
                  <dd className="text-ink">{order.customer_id}</dd>
                </div>
                <div className="flex items-center justify-between">
                  <dt className="text-ink-soft">Origen</dt>
                  <dd className="inline-flex items-center gap-1.5 text-ink">
                    {isWhatsapp ? (
                      <>
                        <WhatsAppIcon className="h-3.5 w-3.5" />
                        {order.ai_handled ? 'WhatsApp (IA)' : 'WhatsApp'}
                      </>
                    ) : (
                      <>
                        <PencilIcon className="h-3.5 w-3.5" />
                        Manual
                      </>
                    )}
                  </dd>
                </div>
                <div className="flex items-center justify-between">
                  <dt className="text-ink-soft">Urgencia</dt>
                  <dd className="text-ink">
                    {order.urgency === 'urgent' ? 'Urgente' : 'Normal'}
                  </dd>
                </div>
                {order.assigned_to && (
                  <div className="flex items-center justify-between">
                    <dt className="text-ink-soft">Asignado a</dt>
                    <dd className="text-ink">{order.assigned_to}</dd>
                  </div>
                )}
              </dl>
            </section>

            <section className="rounded-lg bg-card p-5 shadow-card">
              <h2 className="mb-4 text-[15px] font-bold text-ink">Fechas</h2>
              <dl className="flex flex-col gap-3 text-[13px]">
                <div className="flex items-center justify-between">
                  <dt className="text-ink-soft">Creado</dt>
                  <dd className="text-ink">{formatDateTime(order.created_at)}</dd>
                </div>
                <div className="flex items-center justify-between">
                  <dt className="text-ink-soft">Actualizado</dt>
                  <dd className="text-ink">{formatDateTime(order.updated_at)}</dd>
                </div>
                <div className="flex items-center justify-between">
                  <dt className="text-ink-soft">Completado</dt>
                  <dd className="text-ink">
                    {formatDateTime(order.completed_at)}
                  </dd>
                </div>
                {order.created_by && (
                  <div className="flex items-center justify-between">
                    <dt className="text-ink-soft">Creado por</dt>
                    <dd className="text-ink">{order.created_by}</dd>
                  </div>
                )}
              </dl>
            </section>

            {order.notes && (
              <section className="rounded-lg bg-card p-5 shadow-card">
                <h2 className="mb-3 text-[15px] font-bold text-ink">Notas</h2>
                <p className="text-[13px] text-ink-soft">{order.notes}</p>
              </section>
            )}
          </aside>
        </div>
      </main>
    </>
  );
}
