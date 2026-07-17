import clsx from 'clsx';

interface StatusBadgeProps {
  label: string;
  className?: string;
  pulse?: boolean;
}

export default function StatusBadge({
  label,
  className,
  pulse = false,
}: StatusBadgeProps) {
  return (
    <span
      className={clsx(
        'inline-block whitespace-nowrap rounded-full px-2.5 py-1 text-[11px] font-bold',
        pulse && 'animate-pulse',
        className
      )}
    >
      {label}
    </span>
  );
}
