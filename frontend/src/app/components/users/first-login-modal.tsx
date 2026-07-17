'use client';
import { updateUser } from '@/app/services/user';
import { UpdateUser } from '@/app/types/user';
import { signOut } from 'next-auth/react';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { useActionState } from 'react';
import { Button } from '../common/button';
import { ErrorMessage } from '../common/error-message';
import Modal from '../common/modal';

interface FirstLoginModalProps {
  userId: string;
}

export default function FirstLoginModal({ userId }: FirstLoginModalProps) {
  const [state, dispatch] = useActionState(onSubmit, null);
  const { refresh } = useRouter();

  const [isLoading, setIsLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const [errorMessage, setErrorMessage] = useState<string | undefined>();

  async function onSubmit(_: any, formData: FormData) {
    setIsLoading(true);
    setErrorMessage(undefined);
    try {
      const newPassword = formData.get('new_password')?.toString();
      const confirmPassword = formData.get('confirm_password')?.toString();
      const email = formData.get('email')?.toString();

      if (newPassword != confirmPassword) {
        setErrorMessage('¡Las contraseñas deben coincidir!');
        return;
      }

      if (!newPassword || !email) {
        setErrorMessage('¡Los campos no pueden estar vacíos!');
        return;
      }

      const updateFormData: UpdateUser = {
        updated_by: '',
        email: email,
        password: newPassword,
        first_login_completed: true,
      };

      const resp = await updateUser(userId, updateFormData);
      if (!resp.success) {
        if (resp.unauthorized) {
          signOut({ callbackUrl: '/login' });
        }
        setErrorMessage(resp.message);
        return;
      }

      refresh();
      setIsOpen(false);
    } finally {
      setIsLoading(false);
    }
  }

  return (
    <Modal isOpen={isOpen} onClose={() => setIsOpen(false)} closable={false}>
      <form action={dispatch} className="px-4 py-2">
        <div>
          <h3 className="text-xl font-bold text-ink">Cambiar datos</h3>
          <p className="text-xs text-ink-soft">
            Cambie sus datos de inicio de sesión en el primer acceso
          </p>
        </div>

        <div className="mt-4 space-y-2">
          <label
            className="block text-xs font-semibold uppercase tracking-wider text-ink-faint"
            htmlFor="email"
          >
            Correo electrónico
          </label>

          <input
            className="peer block w-full rounded-md border border-line bg-card px-3 py-2.5 text-sm text-ink outline-none transition-colors placeholder:text-ink-faint focus:border-cut focus:ring-1 focus:ring-cut"
            id="email"
            type="email"
            name="email"
            placeholder="Ingrese su correo electrónico"
            required
          />
        </div>

        <div className="mt-4 space-y-2">
          <label
            className="block text-xs font-semibold uppercase tracking-wider text-ink-faint"
            htmlFor="new_password"
          >
            Nueva contraseña
          </label>

          <input
            className="peer block w-full rounded-md border border-line bg-card px-3 py-2.5 text-sm text-ink outline-none transition-colors placeholder:text-ink-faint focus:border-cut focus:ring-1 focus:ring-cut"
            id="new_password"
            type="password"
            name="new_password"
            placeholder="Ingrese su nueva contraseña"
            required
            minLength={6}
          />
        </div>

        <div className="mt-4 space-y-2">
          <label
            className="block text-xs font-semibold uppercase tracking-wider text-ink-faint"
            htmlFor="confirm_password"
          >
            Confirme su contraseña
          </label>

          <div className="relative">
            <input
              className="peer block w-full rounded-md border border-line bg-card px-3 py-2.5 text-sm text-ink outline-none transition-colors placeholder:text-ink-faint focus:border-cut focus:ring-1 focus:ring-cut"
              id="confirm_password"
              type="password"
              name="confirm_password"
              placeholder="Confirme su nueva contraseña"
              required
              minLength={6}
            />
          </div>
        </div>

        {errorMessage && <ErrorMessage message={errorMessage} />}

        <div className="mt-4">
          <Button type="submit" isLoading={isLoading} size="md">
            Confirmar
          </Button>
        </div>
      </form>
    </Modal>
  );
}
