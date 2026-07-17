'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { StickerMaterial } from '@/app/types/plancha';
import { SearchResponse } from '@/app/types/search_response';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { crmCoreEndpoint, crmCoreApiKey } from './index';

export async function fetchMaterials(): Promise<
  ServiceResponse<SearchResponse<StickerMaterial>>
> {
  try {
    const jwt = (await cookies()).get('jwt')?.value;
    const url = `${crmCoreEndpoint}/crm/core/api/v1/materials`;

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
          'fallo en la búsqueda de los materiales'
        ),
        unauthorized: resp.status === 401,
      };
    }

    const respData = (await resp.json()) as SearchResponse<StickerMaterial>;
    return { success: true, data: respData, message: '' };
  } catch (ex) {
    console.error(ex);
    return { success: false, message: 'algo salió mal' };
  }
}
