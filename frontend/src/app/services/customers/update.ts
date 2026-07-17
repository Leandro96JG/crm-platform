'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { getCurrentUser } from '@/app/libs/session';
import { UpdateCustomerDTO } from '@/app/types/customer';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export async function updateCustomer(
  customerID: string,
  input: UpdateCustomerDTO
): Promise<ServiceResponse<null>> {
  try {
    if (!customerID) {
      return { success: false, message: 'ID del cliente no informado' };
    }

    const jwt = (await cookies()).get('jwt')?.value;
    const session = await getCurrentUser();
    const url = `${crmCoreEndpoint}/crm/core/api/v1/customers/${customerID}`;

    const resp = await fetch(url, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': crmCoreApiKey || '',
        Authorization: `Bearer ${jwt}`,
      },
      body: JSON.stringify({
        ...input,
        updated_by: session?.username || '',
      }),
    });

    if (!resp.ok) {
      return {
        success: false,
        message: await getApiErrorMessage(
          resp,
          'fallo al actualizar el cliente'
        ),
        unauthorized: resp.status === 401,
      };
    }

    revalidatePath('/customers');
    return { success: true, message: 'cliente actualizado con éxito' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
