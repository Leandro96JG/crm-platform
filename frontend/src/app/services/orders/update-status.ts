'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { getCurrentUser } from '@/app/libs/session';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export async function updateOrderStatus(
  orderID: string,
  status: string
): Promise<ServiceResponse<null>> {
  try {
    if (!orderID) {
      return { success: false, message: 'ID del pedido no informado' };
    }

    const jwt = (await cookies()).get('jwt')?.value;
    const session = await getCurrentUser();
    const url = `${crmCoreEndpoint}/crm/core/api/v1/orders/${orderID}/status`;

    const resp = await fetch(url, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': crmCoreApiKey || '',
        Authorization: `Bearer ${jwt}`,
      },
      body: JSON.stringify({
        status,
        updated_by: session?.username || '',
      }),
    });

    if (!resp.ok) {
      return {
        success: false,
        message: await getApiErrorMessage(
          resp,
          'fallo al actualizar el estado del pedido'
        ),
        unauthorized: resp.status === 401,
      };
    }

    revalidatePath('/orders');
    revalidatePath('/home');
    return { success: true, message: 'estado actualizado con éxito' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
