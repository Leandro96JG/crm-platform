'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export async function deleteCustomer(
  customerID: string
): Promise<ServiceResponse<null>> {
  try {
    if (!customerID) {
      return { success: false, message: 'ID del cliente no informado' };
    }

    const jwt = (await cookies()).get('jwt')?.value;
    const url = `${crmCoreEndpoint}/crm/core/api/v1/customers/${customerID}`;

    const resp = await fetch(url, {
      method: 'DELETE',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': crmCoreApiKey || '',
        Authorization: `Bearer ${jwt}`,
      },
    });

    if (!resp.ok) {
      return {
        success: false,
        message: await getApiErrorMessage(resp, 'fallo al eliminar el cliente'),
        unauthorized: resp.status === 401,
      };
    }

    revalidatePath('/customers');
    return { success: true, message: 'cliente eliminado con éxito' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
