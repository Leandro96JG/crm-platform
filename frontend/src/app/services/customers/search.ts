'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { Customer } from '@/app/types/customer';
import { SearchResponse } from '@/app/types/search_response';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export async function fetchCustomers(
  search: string,
  page: number,
  limit: number = 10
): Promise<ServiceResponse<SearchResponse<Customer>>> {
  try {
    page = page - 1;
    const jwt = (await cookies()).get('jwt')?.value;
    const searchParam = search ? `&search=${encodeURIComponent(search)}` : '';
    const url = `${crmCoreEndpoint}/crm/core/api/v1/customers?offset=${page * limit}&limit=${limit}${searchParam}`;

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
        message: await getApiErrorMessage(
          resp,
          'fallo en la búsqueda de los clientes'
        ),
        unauthorized: resp.status === 401,
      };
    }

    const respData = (await resp.json()) as SearchResponse<Customer>;
    return { success: true, data: respData, message: '' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
