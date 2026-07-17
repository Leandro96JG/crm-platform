'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export interface CalculatePriceResult {
  plancha_id: string;
  material_id: string;
  quantity: number;
  total: number;
}

export async function calculatePrice(
  planchaID: string,
  materialID: string,
  quantity: number
): Promise<ServiceResponse<CalculatePriceResult>> {
  try {
    if (!planchaID || !materialID || quantity < 1) {
      return { success: false, message: 'datos incompletos para calcular precio' };
    }

    const jwt = (await cookies()).get('jwt')?.value;
    const url = `${crmCoreEndpoint}/crm/core/api/v1/planchas/${planchaID}/calculate-price?material_id=${materialID}&quantity=${quantity}`;

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
        message: await getApiErrorMessage(resp, 'fallo al calcular el precio'),
        unauthorized: resp.status === 401,
      };
    }

    const data = (await resp.json()) as CalculatePriceResult;
    return { success: true, data, message: '' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
