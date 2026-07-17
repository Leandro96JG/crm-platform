'use client';

import {
  AtSymbolIcon,
  ExclamationCircleIcon,
  KeyIcon,
} from '@heroicons/react/24/outline';
import { ArrowRightIcon } from '@heroicons/react/20/solid';
import { useActionState } from 'react';
import { useFormStatus } from 'react-dom';
import { login } from '../../services/authentication';

export default function LoginForm() {
  const [errorMessage, dispatch] = useActionState(login, null);

  return (
    <form action={dispatch} className="space-y-6">
      <div>
        <h1 className="text-2xl font-bold tracking-tight text-ink">
          Inicia sesión
        </h1>
        <p className="mt-1 text-sm text-ink-soft">
          Accede a tu panel de taller para continuar.
        </p>
      </div>

      <div className="space-y-4">
        <div>
          <label
            className="mb-1.5 block text-xs font-semibold uppercase tracking-wider text-ink-faint"
            htmlFor="email"
          >
            Correo electrónico
          </label>
          <div className="relative">
            <AtSymbolIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-ink-faint peer-focus:text-cut" />
            <input
              className="peer block w-full rounded-md border border-line bg-card py-2.5 pl-10 pr-3 text-sm text-ink outline-none transition-colors placeholder:text-ink-faint focus:border-cut focus:ring-1 focus:ring-cut"
              id="email"
              type="email"
              name="email"
              placeholder="tu@correo.com"
              required
            />
          </div>
        </div>

        <div>
          <label
            className="mb-1.5 block text-xs font-semibold uppercase tracking-wider text-ink-faint"
            htmlFor="password"
          >
            Contraseña
          </label>
          <div className="relative">
            <KeyIcon className="pointer-events-none absolute left-3 top-1/2 h-[18px] w-[18px] -translate-y-1/2 text-ink-faint peer-focus:text-cut" />
            <input
              className="peer block w-full rounded-md border border-line bg-card py-2.5 pl-10 pr-3 text-sm text-ink outline-none transition-colors placeholder:text-ink-faint focus:border-cut focus:ring-1 focus:ring-cut"
              id="password"
              type="password"
              name="password"
              placeholder="Ingrese su contraseña"
              required
              minLength={6}
            />
          </div>
        </div>
      </div>

      <LoginButton />

      <div className="flex h-6 items-center gap-1.5" aria-live="polite" aria-atomic="true">
        {errorMessage && (
          <>
            <ExclamationCircleIcon className="h-5 w-5 text-st-danger-fg" />
            <p className="text-sm text-st-danger-fg">{errorMessage}</p>
          </>
        )}
      </div>
    </form>
  );
}

function LoginButton() {
  const { pending } = useFormStatus();

  return (
    <button
      type="submit"
      disabled={pending}
      aria-disabled={pending}
      className="flex w-full items-center justify-center gap-2 rounded-md bg-cut px-4 py-2.5 text-sm font-semibold text-white transition-colors hover:bg-cut-dark focus-visible:outline focus-visible:outline-2 focus-visible:outline-offset-2 focus-visible:outline-cut disabled:cursor-not-allowed disabled:opacity-60"
    >
      {pending ? (
        <svg
          className="h-5 w-5 animate-spin"
          xmlns="http://www.w3.org/2000/svg"
          fill="none"
          viewBox="0 0 24 24"
        >
          <circle
            className="opacity-25"
            cx="12"
            cy="12"
            r="10"
            stroke="currentColor"
            strokeWidth="4"
          />
          <path
            className="opacity-75"
            fill="currentColor"
            d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4z"
          />
        </svg>
      ) : (
        <>
          Iniciar sesión
          <ArrowRightIcon className="h-5 w-5" />
        </>
      )}
    </button>
  );
}
