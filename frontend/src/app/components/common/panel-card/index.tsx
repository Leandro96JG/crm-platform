import Link from 'next/link';
import { ReactNode } from 'react';

interface PanelCardProps {
  title: string;
  subtitle?: string;
  action?: { label: string; href: string };
  children: ReactNode;
  bodyPadding?: boolean;
}

export default function PanelCard({
  title,
  subtitle,
  action,
  children,
  bodyPadding = false,
}: PanelCardProps) {
  return (
    <div className="overflow-hidden rounded-lg bg-card shadow-card">
      <div className="flex items-center justify-between px-5 pb-3.5 pt-4">
        <div>
          <h2 className="text-[15px] font-bold text-ink">{title}</h2>
          {subtitle && (
            <p className="mt-0.5 text-[11.5px] text-ink-soft">{subtitle}</p>
          )}
        </div>
        {action && (
          <Link
            href={action.href}
            className="text-xs font-semibold text-cut-dark hover:underline"
          >
            {action.label} →
          </Link>
        )}
      </div>
      <div className={bodyPadding ? 'px-5 pb-5' : ''}>{children}</div>
    </div>
  );
}
