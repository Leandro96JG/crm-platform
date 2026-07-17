'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { getCurrentUser } from '@/app/libs/session';
import { ServiceResponse } from '@/app/types/service';
import { UpdateUser } from '@/app/types/user';
import { cookies } from 'next/headers';
import { crmCoreApiKey, crmCoreEndpoint } from '.';

export async function updateUser(
  userID: string,
  update: UpdateUser
): Promise<ServiceResponse<null>> {
  try {
    if (userID == '') {
      return {
        success: false,
        message: 'ID del usuario no informado',
      };
    }

    const jwt = (await cookies()).get('jwt')?.value;
    const url = `${crmCoreEndpoint}/crm/core/api/v1/users/${userID}`;
    const session = await getCurrentUser();
    const author = session?.username || '';

    const payload: Record<string, unknown> = {
      updated_by: author,
    };

    if (update.first_name !== undefined || update.last_name !== undefined) {
      payload.name = [update.first_name, update.last_name]
        .filter(Boolean)
        .join(' ')
        .trim();
    }
    if (update.email !== undefined) {
      payload.email = update.email;
    }
    if (update.active !== undefined) {
      payload.is_active = update.active;
    }

    const resp = await fetch(url, {
      method: 'PUT',
      headers: {
        'Content-Type': 'application/json',
        'X-API-Key': crmCoreApiKey || '',
        Authorization: `Bearer ${jwt}`,
      },
      body: JSON.stringify(payload),
    });

    if (!resp.ok) {
      const unauthorized = resp.status === 401;
      const errorMessageDefault = unauthorized
        ? 'usuario no autorizado'
        : 'fallo en la actualización del usuario';

      const errorMessage = await getApiErrorMessage(resp, errorMessageDefault);

      return {
        success: false,
        message: errorMessage,
        unauthorized: unauthorized,
      };
    }

    return {
      success: true,
      message: 'usuario actualizado con éxito',
    };
  } catch (error) {
    return {
      success: false,
      message: 'algo salió mal, contacte al soporte!',
    };
  }
}
