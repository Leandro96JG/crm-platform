'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { PrintJob } from '@/app/types/print_job';
import { SearchResponse } from '@/app/types/search_response';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export async function fetchPrintJobs(
  jobType: string,
  status: string[],
  page: number = 0,
  limit: number = 100
): Promise<ServiceResponse<SearchResponse<PrintJob>>> {
  try {
    const jwt = (await cookies()).get('jwt')?.value;
    const statusQuery = status.map(s => `status=${s}`).join('&');
    const url = `${crmCoreEndpoint}/crm/core/api/v1/print-jobs?job_type=${jobType}&${statusQuery}&offset=${page}&limit=${limit}`;

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
        message: await getApiErrorMessage(resp, 'fallo en la búsqueda de los trabajos'),
        unauthorized: resp.status === 401,
      };
    }

    const respData = (await resp.json()) as SearchResponse<PrintJob>;
    return { success: true, data: respData, message: '' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
