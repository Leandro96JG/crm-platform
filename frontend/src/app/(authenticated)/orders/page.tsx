'use server';
import { unauthorizedRedirect } from '@/app/libs/auth-redirect';
import { getCurrentUser } from '@/app/libs/session';
import { Order } from '@/app/types/order';
import { SearchResponse } from '@/app/types/search_response';
import { redirect } from 'next/navigation';
import OrdersTable from '@/app/components/orders/table';
import Topbar from '@/app/components/common/topbar';
import { fetchOrders } from '@/app/services/orders';
import { getInitials } from '@/app/utils/user-display';
import Link from 'next/link';

type OrdersPageParams = {
  searchParams: Promise<{
    status?: string;
    page?: number;
  }>;
};

async function getData(
  status: string,
  page: number
): Promise<SearchResponse<Order>> {
  let query = '';
  if (status) {
    query = `status=${status}`;
  }

  const { success, unauthorized, data } = await fetchOrders(query, page);
  if (!success || !data) {
    if (unauthorized) {
      unauthorizedRedirect();
    }
    return { result: [], paging: { total: 0, limit: 10, offset: 0 } };
  }

  return data;
}

export default async function OrdersPage({ searchParams }: OrdersPageParams) {
  const user = await getCurrentUser();
  if (!user) {
    redirect('/login');
  }

  const { status, page } = await searchParams;
  const data = await getData(status || '', page || 1);

  return (
    <>
      <Topbar
        title="Pedidos"
        subtitle="Gestión de pedidos del taller"
        initials={getInitials(user.name ?? user.username)}
      />
      <div className="cut-divider" />
      <main className="animate-fadeIn px-8 pb-10">
        <div className="mb-4 flex justify-end">
          <Link
            href="/orders/new"
            className="rounded-lg bg-cut-dark px-4 py-2 text-sm font-semibold text-white transition-colors hover:opacity-90"
          >
            + Nuevo pedido
          </Link>
        </div>
        <OrdersTable orders={data.result} paging={data.paging} />
      </main>
    </>
  );
}
