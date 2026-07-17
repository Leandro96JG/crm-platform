import Image from 'next/image';
import LoginForm from '../components/login/login-form';

export default function LoginPage() {
  return (
    <main className="grid min-h-screen bg-paper lg:grid-cols-[1.1fr_1fr]">
      <section className="relative hidden flex-col justify-between overflow-hidden bg-ink p-12 lg:flex">
        <div className="flex items-center gap-3">
          <Image
            src="/viva-logo.png"
            alt="Viva"
            width={44}
            height={44}
            className="rounded-md"
          />
          <span className="text-xl font-bold tracking-tight text-paper">
            Viva
          </span>
        </div>

        <div className="relative z-10 max-w-sm">
          <h2 className="text-3xl font-bold leading-tight text-paper">
            Gestión de taller,
            <br />
            <span className="text-cut">simple y precisa.</span>
          </h2>
          <p className="mt-4 text-sm leading-relaxed text-paper/60">
            Pedidos, planchas y producción en un solo lugar. Corta, imprime y
            entrega sin perder el hilo.
          </p>
        </div>

        <div className="register-mark absolute -right-16 -top-16 h-64 w-64 opacity-20" />
        <div className="register-mark absolute -bottom-20 left-1/3 h-52 w-52 opacity-10" />
      </section>

      <section className="flex items-center justify-center px-6 py-12">
        <div className="w-full max-w-[380px]">
          <div className="mb-8 flex items-center gap-3 lg:hidden">
            <Image
              src="/viva-logo.png"
              alt="Viva"
              width={40}
              height={40}
              className="rounded-md"
            />
            <span className="text-lg font-bold tracking-tight text-ink">
              Viva
            </span>
          </div>
          <LoginForm />
        </div>
      </section>
    </main>
  );
}
