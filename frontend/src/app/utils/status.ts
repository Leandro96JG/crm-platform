import { OrderStatus } from '@/app/types/order';
import { PrintJobStatus } from '@/app/types/print_job';

export interface StatusMeta {
  label: string;
  badge: string;
}

export const orderStatusMeta: Record<OrderStatus, StatusMeta> = {
  pending: { label: 'Pendiente', badge: 'bg-st-warn-bg text-st-warn-fg' },
  approved: { label: 'Aprobado', badge: 'bg-st-info-bg text-st-info-fg' },
  in_production: {
    label: 'En producción',
    badge: 'bg-st-prod-bg text-st-prod-fg',
  },
  ready: { label: 'Listo', badge: 'bg-st-ok-bg text-st-ok-fg' },
  delivered: { label: 'Entregado', badge: 'bg-st-done-bg text-st-done-fg' },
  cancelled: { label: 'Cancelado', badge: 'bg-st-danger-bg text-st-danger-fg' },
};

export function getOrderStatusMeta(status: string): StatusMeta {
  return (
    orderStatusMeta[status as OrderStatus] ?? {
      label: status,
      badge: 'bg-paper-dim text-ink-soft',
    }
  );
}

export const printJobStatusMeta: Record<PrintJobStatus, StatusMeta> = {
  queued: { label: 'En cola', badge: 'bg-paper-dim text-ink-soft' },
  printing: { label: 'Imprimiendo', badge: 'bg-st-prod-bg text-st-prod-fg' },
  printed: { label: 'Impreso', badge: 'bg-st-ok-bg text-st-ok-fg' },
  cutting: { label: 'Cortando', badge: 'bg-st-prod-bg text-st-prod-fg' },
  cut: { label: 'Cortado', badge: 'bg-st-ok-bg text-st-ok-fg' },
  failed: { label: 'Fallido', badge: 'bg-st-danger-bg text-st-danger-fg' },
};

export function getPrintJobStatusMeta(status: string): StatusMeta {
  return (
    printJobStatusMeta[status as PrintJobStatus] ?? {
      label: status,
      badge: 'bg-paper-dim text-ink-soft',
    }
  );
}

export const orderStatusTransitions: Record<string, string[]> = {
  pending: ['approved', 'cancelled'],
  approved: ['in_production', 'cancelled'],
  in_production: ['ready', 'cancelled'],
  ready: ['delivered'],
  delivered: [],
  cancelled: [],
};

export const printJobStatusTransitions: Record<string, string[]> = {
  queued: ['printing', 'cutting', 'failed'],
  printing: ['printed', 'failed'],
  printed: ['cutting', 'failed'],
  cutting: ['cut', 'failed'],
  cut: [],
  failed: ['queued'],
};

export function getNextOrderStatuses(status: string): string[] {
  return orderStatusTransitions[status] ?? [];
}

export function getNextPrintJobStatuses(status: string): string[] {
  return printJobStatusTransitions[status] ?? [];
}

const currencyFormatter = new Intl.NumberFormat('es-AR', {
  style: 'currency',
  currency: 'ARS',
});

export function formatCurrency(value: number): string {
  return currencyFormatter.format(value ?? 0);
}

const dateTimeFormatter = new Intl.DateTimeFormat('es-AR', {
  day: '2-digit',
  month: '2-digit',
  year: 'numeric',
  hour: '2-digit',
  minute: '2-digit',
});

export function formatDateTime(value: string | null | undefined): string {
  if (!value) return '-';
  const date = new Date(value);
  if (isNaN(date.getTime())) return '-';
  return dateTimeFormatter.format(date);
}
