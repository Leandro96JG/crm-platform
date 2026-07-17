'use server';
import Topbar from '@/app/components/common/topbar';
import OrderForm from '@/app/components/orders/order-form';
import { getCurrentUser } from '@/app/libs/session';
import { fetchPlanchas, fetchMaterials } from '@/app/services/planchas';
import { fetchCustomers } from '@/app/services/customers';
import { getInitials } from '@/app/utils/user-display';
import { redirect } from 'next/navigation';
import Link from 'next/link';

export default async function NewOrderPage() {
  const user = await getCurrentUser();
  if (!user) {
    redirect('/login');
  }

  const [planchasResp, materialsResp, customersResp] = await Promise.all([
    fetchPlanchas('is_active=true', 1, 100),
    fetchMaterials(),
    fetchCustomers('', 1, 100),
  ]);

  const planchas = planchasResp.data?.result ?? [];
  const materials = materialsResp.data?.result ?? [];
  const customers = customersResp.data?.result ?? [];

  return (
    <>
      <Topbar
        title="Nuevo pedido"
        subtitle="Cargar un pedido manualmente"
        initials={getInitials(user.name ?? user.username)}
      />
      <div className="cut-divider" />
      <main className="animate-fadeIn px-8 pb-10">
        <div className="mb-4">
          <Link
            href="/orders"
            className="text-[13px] font-semibold text-cut-dark hover:underline"
          >
            ← Volver a pedidos
          </Link>
        </div>
        {planchas.length === 0 || materials.length === 0 ? (
          <div className="rounded-lg bg-card p-6 text-sm text-ink-soft shadow-card">
            {planchas.length === 0
              ? 'No hay planchas activas para crear un pedido.'
              : 'No hay materiales cargados para crear un pedido.'}
          </div>
        ) : (
          <OrderForm
            planchas={planchas}
            materials={materials}
            customers={customers}
          />
        )}
      </main>
    </>
  );
}
