'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { getCurrentUser } from '@/app/libs/session';
import { CreateCustomerDTO } from '@/app/types/customer';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export async function createCustomer(
  input: CreateCustomerDTO
): Promise<ServiceResponse<{ customer_id: string }>> {
  try {
    if (!input.name?.trim()) {
      return { success: false, message: 'El nombre es obligatorio' };
    }

    const jwt = (await cookies()).get('jwt')?.value;
    const session = await getCurrentUser();
    const url = `${crmCoreEndpoint}/crm/core/api/v1/customers`;

    const resp = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': crmCoreApiKey || '',
        Authorization: `Bearer ${jwt}`,
      },
      body: JSON.stringify({
        name: input.name.trim(),
        phone: input.phone || '',
        email: input.email || '',
        document: input.document || '',
        address: input.address || '',
        notes: input.notes || '',
        created_by: session?.username || '',
      }),
    });

    if (!resp.ok) {
      return {
        success: false,
        message: await getApiErrorMessage(resp, 'fallo al crear el cliente'),
        unauthorized: resp.status === 401,
      };
    }

    const data = (await resp.json()) as { customer_id: string };

    revalidatePath('/customers');
    return { success: true, data, message: 'cliente creado con éxito' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
