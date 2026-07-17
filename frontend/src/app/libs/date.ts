import { format, Locale, parseISO } from 'date-fns';
import { es } from 'date-fns/locale/es';

export const PATTERN_DEFAULT = 'dd/MM/yyyy HH:mm';
export const API_PATTERN_DEFAULT = 'yyyy-MM-ddTHH:mm:ss.SSSSSS';
export const ONLY_DATE_PATTERN = 'dd/MM/yyyy';

const locales: Record<string, Locale> = { es };

function getFormat(date: string, pattern = PATTERN_DEFAULT) {
  return format(parseISO(date), pattern, { locale: locales["es"] });
}

export function parseDateTime(date: string, pattern: string = PATTERN_DEFAULT) {
  if (!date || date === '') return '';

  try {
    return getFormat(date, pattern);
  } catch (error) {
    return '';
  }
}

export function timeElapsed(startDate: Date, endDate: Date): string {
  const diffInMs: number = endDate.getTime() - startDate.getTime();

  const diffInHours: number = diffInMs / (1000 * 60 * 60);

  if (diffInHours < 24) {
    return `${diffInHours.toFixed(0)} horas`;
  }

  const diffInDays: number = diffInMs / (1000 * 60 * 60 * 24);

  return `${diffInDays.toFixed(0)} días`;
}
