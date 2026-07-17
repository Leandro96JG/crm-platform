'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { getCurrentUser } from '@/app/libs/session';
import { CreateOrderItemDTO } from '@/app/types/order';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { revalidatePath } from 'next/cache';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

interface CreateOrderInput {
  customer_id: string;
  source: string;
  notes: string;
  urgency: string;
  items: CreateOrderItemDTO[];
}

export async function createOrder(
  input: CreateOrderInput
): Promise<ServiceResponse<{ order_id: string; order_number: string }>> {
  try {
    if (!input.customer_id?.trim()) {
      return { success: false, message: 'El cliente es obligatorio' };
    }
    if (!input.items || input.items.length === 0) {
      return { success: false, message: 'El pedido debe tener al menos un ítem' };
    }

    const jwt = (await cookies()).get('jwt')?.value;
    const session = await getCurrentUser();
    const url = `${crmCoreEndpoint}/crm/core/api/v1/orders`;

    const resp = await fetch(url, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': crmCoreApiKey || '',
        Authorization: `Bearer ${jwt}`,
      },
      body: JSON.stringify({
        customer_id: input.customer_id.trim(),
        source: input.source || 'manual',
        notes: input.notes || '',
        urgency: input.urgency || 'normal',
        created_by: session?.username || '',
        items: input.items,
      }),
    });

    if (!resp.ok) {
      return {
        success: false,
        message: await getApiErrorMessage(resp, 'fallo al crear el pedido'),
        unauthorized: resp.status === 401,
      };
    }

    const data = (await resp.json()) as {
      order_id: string;
      order_number: string;
    };

    revalidatePath('/orders');
    revalidatePath('/home');
    return { success: true, data, message: 'pedido creado con éxito' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
