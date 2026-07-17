'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { Order } from '@/app/types/order';
import { SearchResponse } from '@/app/types/search_response';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export async function fetchOrders(
  query: string,
  page: number,
  limit: number = 10
): Promise<ServiceResponse<SearchResponse<Order>>> {
  try {
    page = page - 1;
    const jwt = (await cookies()).get('jwt')?.value;
    const url = `${crmCoreEndpoint}/crm/core/api/v1/orders?offset=${page * limit}&limit=${limit}&${query}`;

    const resp = await fetch(url, {
      method: 'GET',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': crmCoreApiKey || '',
        Authorization: `Bearer ${jwt}`,
      },
    });

    if (!resp.ok) {
      return {
        success: false,
        message: await getApiErrorMessage(resp, 'fallo en la búsqueda de los pedidos'),
        unauthorized: resp.status === 401,
      };
    }

    const respData = (await resp.json()) as SearchResponse<Order>;
    return { success: true, data: respData, message: '' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
