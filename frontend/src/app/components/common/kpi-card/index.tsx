import clsx from 'clsx';
import { ReactNode } from 'react';

type Tone = 'cut' | 'warn' | 'teal' | 'info';

const toneClasses: Record<Tone, string> = {
  cut: 'bg-st-danger-bg text-cut-dark',
  warn: 'bg-st-warn-bg text-st-warn-fg',
  teal: 'bg-st-prod-bg text-st-prod-fg',
  info: 'bg-st-info-bg text-st-info-fg',
};

type DeltaKind = 'up' | 'flat' | 'attn';

const deltaClasses: Record<DeltaKind, string> = {
  up: 'text-st-ok-fg',
  flat: 'text-ink-faint',
  attn: 'text-st-warn-fg',
};

interface KpiCardProps {
  label: string;
  value: ReactNode;
  icon: ReactNode;
  tone: Tone;
  delta?: string;
  deltaKind?: DeltaKind;
}

export default function KpiCard({
  label,
  value,
  icon,
  tone,
  delta,
  deltaKind = 'flat',
}: KpiCardProps) {
  return (
    <div className="register-mark relative overflow-hidden rounded-lg bg-card px-5 pb-4 pt-5 shadow-card">
      <div className="mb-3.5 flex items-center gap-2.5">
        <div
          className={clsx(
            'flex h-[34px] w-[34px] flex-none items-center justify-center rounded-[10px] [&>svg]:h-[17px] [&>svg]:w-[17px]',
            toneClasses[tone]
          )}
        >
          {icon}
        </div>
        <div className="text-[12.5px] font-medium text-ink-soft">{label}</div>
      </div>
      <div className="mb-2 font-mono text-3xl font-bold leading-none tracking-tight">
        {value}
      </div>
      {delta && (
        <span
          className={clsx(
            'inline-flex items-center gap-1 text-[11.5px] font-semibold',
            deltaClasses[deltaKind]
          )}
        >
          {delta}
        </span>
      )}
    </div>
  );
}
