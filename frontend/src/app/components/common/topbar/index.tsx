import ThemeToggle from '@/app/components/common/theme-toggle';

interface TopbarProps {
  title: string;
  subtitle?: string;
  initials: string;
}

export default function Topbar({ title, subtitle, initials }: TopbarProps) {
  return (
    <header className="flex items-center justify-between gap-6 px-8 py-5">
      <div>
        <h1 className="text-[21px] font-bold tracking-tight text-ink">{title}</h1>
        {subtitle && (
          <p className="mt-0.5 text-[12.5px] text-ink-soft">{subtitle}</p>
        )}
      </div>

      <div className="flex max-w-[380px] flex-1 items-center gap-2.5 rounded-md border border-line bg-card px-3.5 py-2.5">
        <svg
          viewBox="0 0 24 24"
          fill="none"
          stroke="currentColor"
          strokeWidth={2}
          strokeLinecap="round"
          strokeLinejoin="round"
          className="h-4 w-4 flex-none text-ink-faint"
        >
          <circle cx="11" cy="11" r="7" />
          <path d="M21 21l-4.3-4.3" />
        </svg>
        <input
          type="text"
          placeholder="Buscar pedido, cliente o plancha…"
          className="w-full border-none bg-transparent text-[13px] text-ink outline-none placeholder:text-ink-faint focus:ring-0"
        />
      </div>

      <div className="flex items-center gap-4">
        <ThemeToggle />
        <button
          type="button"
          aria-label="Notificaciones"
          className="relative flex h-9 w-9 items-center justify-center rounded-sm border border-line bg-card text-ink-soft"
        >
          <svg
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            strokeWidth={2}
            strokeLinecap="round"
            strokeLinejoin="round"
            className="h-4 w-4"
          >
            <path d="M18 8a6 6 0 00-12 0c0 7-3 9-3 9h18s-3-2-3-9" />
            <path d="M13.7 21a2 2 0 01-3.4 0" />
          </svg>
          <span className="absolute right-[7px] top-[7px] h-[7px] w-[7px] rounded-full border-[1.5px] border-card bg-cut" />
        </button>
        <div className="flex h-[38px] w-[38px] items-center justify-center rounded-full bg-ink text-[13px] font-bold text-paper">
          {initials}
        </div>
      </div>
    </header>
  );
}
