'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { Customer } from '@/app/types/customer';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export async function fetchCustomer(
  customerID: string
): Promise<ServiceResponse<Customer>> {
  try {
    const jwt = (await cookies()).get('jwt')?.value;
    const url = `${crmCoreEndpoint}/crm/core/api/v1/customers/${customerID}`;

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
        message: await getApiErrorMessage(resp, 'fallo al obtener el cliente'),
        unauthorized: resp.status === 401,
      };
    }

    const respData = (await resp.json()) as Customer;
    return { success: true, data: respData, message: '' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
