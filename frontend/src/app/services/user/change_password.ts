'use server';
import { getApiErrorMessage } from '@/app/libs/api-error';
import { ServiceResponse } from '@/app/types/service';
import { cookies } from 'next/headers';
import { crmCoreApiKey, crmCoreEndpoint } from '.';

export type ChangePassword = {
  old_password: string;
  new_password: string;
};

export async function changePassword(
  userID: string,
  payload: ChangePassword
): Promise<ServiceResponse<null>> {
  try {
    if (userID == '') {
      return {
        success: false,
        message: 'ID del usuario no informado',
      };
    }

    const jwt = (await cookies()).get('jwt')?.value;
    const url = `${crmCoreEndpoint}/crm/core/api/v1/users/${userID}/password`;

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
        : 'fallo en el cambio de contraseña';

      const errorMessage = await getApiErrorMessage(resp, errorMessageDefault);

      return {
        success: false,
        message: errorMessage,
        unauthorized: unauthorized,
      };
    }

    return {
      success: true,
      message: 'contraseña actualizada con éxito',
    };
  } catch (error) {
    return {
      success: false,
      message: 'algo salió mal, contacte al soporte!',
    };
  }
}
