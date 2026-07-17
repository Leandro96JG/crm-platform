'use server';
import Topbar from '@/app/components/common/topbar';
import CustomersManager from '@/app/components/customers/customers-manager';
import { getCurrentUser } from '@/app/libs/session';
import { fetchCustomers } from '@/app/services/customers';
import { getInitials } from '@/app/utils/user-display';
import { redirect } from 'next/navigation';

type CustomersPageParams = {
  searchParams: Promise<{
    query?: string;
    page?: number;
  }>;
};

export default async function CustomersPage({
  searchParams,
}: CustomersPageParams) {
  const user = await getCurrentUser();
  if (!user) {
    redirect('/login');
  }

  const { query, page } = await searchParams;
  const resp = await fetchCustomers(query || '', page || 1);

  if (resp.unauthorized) {
    redirect('/login');
  }

  const customers = resp.data?.result ?? [];
  const total = resp.data?.paging.total ?? 0;

  return (
    <>
      <Topbar
        title="Clientes"
        subtitle="Gestión de clientes del taller"
        initials={getInitials(user.name ?? user.username)}
      />
      <div className="cut-divider" />
      <main className="animate-fadeIn px-8 pb-10">
        <CustomersManager
          customers={customers}
          total={total}
          initialSearch={query || ''}
        />
      </main>
    </>
  );
}
