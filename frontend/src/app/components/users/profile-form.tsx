'use client';
import { useSnackbar } from '@/app/context/SnackbarProvider';
import { changePassword, isEmailTaken, updateUser } from '@/app/services/user';
import { User } from '@/app/types/user';
import { roleLabels } from '@/app/utils/roles';
import { signOut } from 'next-auth/react';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { Button } from '../common/button';
import { ErrorMessage } from '../common/error-message';
import { TextInput } from '../common/text-input/text-input';

function ReadOnlyField({
  label,
  value,
  className,
}: {
  label: string;
  value: string;
  className?: string;
}) {
  return (
    <div className={className}>
      <p className="mb-1 block text-xs font-semibold uppercase tracking-wider text-ink-faint">{label}</p>
      <p className="block w-full text-sm text-ink">{value}</p>
    </div>
  );
}

const PASSWORD_MIN_LENGTH = 8;
const PASSWORD_SPECIAL_CHAR_REGEX = /[^A-Za-z0-9]/;
const PASSWORD_NUMBER_REGEX = /[0-9]/;

function validateNewPassword(password: string): string | undefined {
  if (password.length < PASSWORD_MIN_LENGTH) {
    return `La contraseña debe tener al menos ${PASSWORD_MIN_LENGTH} caracteres`;
  }
  if (!PASSWORD_NUMBER_REGEX.test(password)) {
    return 'La contraseña debe contener al menos 1 número';
  }
  if (!PASSWORD_SPECIAL_CHAR_REGEX.test(password)) {
    return 'La contraseña debe contener al menos 1 carácter especial';
  }
  return undefined;
}

interface ProfileFormProps {
  user: User;
}

export function ProfileForm({ user }: ProfileFormProps) {
  const { showSnackbar } = useSnackbar();
  const { refresh } = useRouter();

  const [isEditingInfo, setIsEditingInfo] = useState(false);
  const [isSavingInfo, setIsSavingInfo] = useState(false);
  const [infoError, setInfoError] = useState<string | undefined>();

  const [isSavingPassword, setIsSavingPassword] = useState(false);
  const [passwordError, setPasswordError] = useState<string | undefined>();
  const [passwordFormKey, setPasswordFormKey] = useState(0);

  async function onSubmitInfo(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    setInfoError(undefined);
    setIsSavingInfo(true);
    try {
      const firstName = formData.get('first_name')?.toString().trim();
      const lastName = formData.get('last_name')?.toString().trim();
      const email = formData.get('email')?.toString().trim();

      if (!firstName || !lastName || !email) {
        setInfoError('Los campos no pueden estar vacíos');
        return;
      }

      if (email !== user.email) {
        const emailTaken = await isEmailTaken(email, user.user_id);
        if (emailTaken) {
          setInfoError('Este correo electrónico ya está en uso por otro usuario');
          return;
        }
      }

      const resp = await updateUser(user.user_id, {
        first_name: firstName,
        last_name: lastName,
        email: email,
        updated_by: user.username,
      });

      if (!resp.success) {
        if (resp.unauthorized) {
          signOut({ callbackUrl: '/login' });
          return;
        }
        setInfoError(resp.message);
        return;
      }

      showSnackbar('Datos actualizados con éxito', 'success');
      setIsEditingInfo(false);
      refresh();
    } finally {
      setIsSavingInfo(false);
    }
  }

  async function onSubmitPassword(event: React.FormEvent<HTMLFormElement>) {
    event.preventDefault();
    const formData = new FormData(event.currentTarget);
    setPasswordError(undefined);
    setIsSavingPassword(true);
    try {
      const currentPassword = formData.get('current_password')?.toString();
      const newPassword = formData.get('new_password')?.toString();
      const confirmPassword = formData.get('confirm_password')?.toString();

      if (!currentPassword || !newPassword || !confirmPassword) {
        setPasswordError('Los campos no pueden estar vacíos');
        return;
      }

      if (newPassword !== confirmPassword) {
        setPasswordError('La nueva contraseña y la confirmación deben coincidir');
        return;
      }

      const complexityError = validateNewPassword(newPassword);
      if (complexityError) {
        setPasswordError(complexityError);
        return;
      }

      const resp = await changePassword(user.user_id, {
        old_password: currentPassword,
        new_password: newPassword,
      });

      if (!resp.success) {
        if (resp.unauthorized) {
          signOut({ callbackUrl: '/login' });
          return;
        }
        setPasswordError(resp.message);
        return;
      }

      showSnackbar('Contraseña actualizada con éxito', 'success');
      setPasswordFormKey((key) => key + 1);
    } finally {
      setIsSavingPassword(false);
    }
  }

  return (
    <div className="flex flex-col gap-6">
      <form
        onSubmit={onSubmitInfo}
        className="rounded-lg border border-line p-6"
      >
        <div className="mb-4 flex items-start justify-between">
          <div>
            <h2 className="text-lg font-semibold text-ink">Mis datos</h2>
            <p className="text-xs text-ink-soft">
              {isEditingInfo
                ? 'Actualice su nombre y correo de acceso'
                : 'Sus datos de acceso y registro'}
            </p>
          </div>

          {!isEditingInfo && (
            <Button
              type="button"
              size="sm"
              onClick={() => setIsEditingInfo(true)}
            >
              Editar
            </Button>
          )}
        </div>

        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <ReadOnlyField label="Usuario" value={user.username} />
          <ReadOnlyField label="Perfil" value={roleLabels[user.role]} />

          {isEditingInfo ? (
            <>
              <TextInput
                name="first_name"
                label="Nombre"
                placeholder="Digite su nombre"
                defaultValue={user.first_name}
                required
              />
              <TextInput
                name="last_name"
                label="Apellido"
                placeholder="Digite su apellido"
                defaultValue={user.last_name}
                required
              />
              <TextInput
                name="email"
                label="Correo electrónico"
                type="email"
                placeholder="Digite su correo electrónico"
                defaultValue={user.email}
                className="md:col-span-2"
                required
              />
            </>
          ) : (
            <>
              <ReadOnlyField label="Nombre" value={user.first_name} />
              <ReadOnlyField label="Apellido" value={user.last_name} />
              <ReadOnlyField
                label="Correo electrónico"
                value={user.email}
                className="md:col-span-2"
              />
            </>
          )}
        </div>

        {infoError && <ErrorMessage message={infoError} />}

        {isEditingInfo && (
          <div className="mt-4 flex gap-3">
            <Button type="submit" size="md" isLoading={isSavingInfo}>
              Guardar datos
            </Button>
            <Button
              type="button"
              size="md"
              color="warning"
              onClick={() => {
                setInfoError(undefined);
                setIsEditingInfo(false);
              }}
            >
              Cancelar
            </Button>
          </div>
        )}
      </form>

      <form
        key={passwordFormKey}
        onSubmit={onSubmitPassword}
        className="rounded-lg border border-line p-6"
      >
        <h2 className="text-lg font-semibold text-ink">Cambiar contraseña</h2>
        <p className="mb-4 text-xs text-ink-soft">
          Mínimo de {PASSWORD_MIN_LENGTH} caracteres, 1 número y 1 carácter
          especial
        </p>

        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <TextInput
            name="current_password"
            label="Contraseña actual"
            type="password"
            placeholder="Digite su contraseña actual"
            className="md:col-span-2"
            required
          />
          <TextInput
            name="new_password"
            label="Nueva contraseña"
            type="password"
            placeholder="Digite su nueva contraseña"
            required
          />
          <TextInput
            name="confirm_password"
            label="Confirme la nueva contraseña"
            type="password"
            placeholder="Confirme su nueva contraseña"
            required
          />
        </div>

        {passwordError && <ErrorMessage message={passwordError} />}

        <div className="mt-4">
          <Button type="submit" size="md" isLoading={isSavingPassword}>
            Cambiar contraseña
          </Button>
        </div>
      </form>
    </div>
  );
}
