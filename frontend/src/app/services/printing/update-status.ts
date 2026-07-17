'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { getCurrentUser } from '@/app/libs/session';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export async function updatePrintJobStatus(
  jobID: string,
  status: string
): Promise<ServiceResponse<null>> {
  try {
    if (!jobID) {
      return { success: false, message: 'ID del trabajo no informado' };
    }

    const jwt = (await cookies()).get('jwt')?.value;
    const session = await getCurrentUser();
    const url = `${crmCoreEndpoint}/crm/core/api/v1/print-jobs/${jobID}/status`;

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
          'fallo al actualizar el estado del trabajo'
        ),
        unauthorized: resp.status === 401,
      };
    }

    revalidatePath('/printing');
    revalidatePath('/home');
    return { success: true, message: 'estado actualizado con éxito' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
