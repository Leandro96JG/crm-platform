'use server';
import { unauthorizedRedirect } from '@/app/libs/auth-redirect';
import { getCurrentUser } from '@/app/libs/session';
import { Plancha } from '@/app/types/plancha';
import { SearchResponse } from '@/app/types/search_response';
import { redirect } from 'next/navigation';
import StickersTable from '@/app/components/stickers/table';
import Topbar from '@/app/components/common/topbar';
import { fetchPlanchas } from '@/app/services/planchas';
import { getInitials } from '@/app/utils/user-display';

type StickersPageParams = {
  searchParams: Promise<{
    category?: string;
    page?: number;
  }>;
};

async function getData(
  category: string,
  page: number
): Promise<SearchResponse<Plancha>> {
  let query = '';
  if (category) {
    query = `category=${category}`;
  }

  const { success, unauthorized, data } = await fetchPlanchas(query, page);
  if (!success || !data) {
    if (unauthorized) {
      unauthorizedRedirect();
    }
    return { result: [], paging: { total: 0, limit: 10, offset: 0 } };
  }

  return data;
}

export default async function StickersPage({
  searchParams,
}: StickersPageParams) {
  const user = await getCurrentUser();
  if (!user) {
    redirect('/login');
  }

  const { category, page } = await searchParams;
  const data = await getData(category || '', page || 1);

  return (
    <>
      <Topbar
        title="Catálogo de planchas"
        subtitle="Diseños disponibles en el taller"
        initials={getInitials(user.name ?? user.username)}
      />
      <div className="cut-divider" />
      <main className="animate-fadeIn px-8 pb-10">
        <StickersTable planchas={data.result} paging={data.paging} />
      </main>
    </>
  );
}
