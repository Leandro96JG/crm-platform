import { UserRole } from '@/app/types/user';

export function getInitials(name?: string | null): string {
  if (!name) {
    return 'U';
  }
  const parts = name.trim().split(/\s+/).filter(Boolean);
  if (parts.length === 0) {
    return 'U';
  }
  if (parts.length === 1) {
    return parts[0].slice(0, 2).toUpperCase();
  }
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

export const roleLabels: Record<string, string> = {
  [UserRole.THAVANNA_ADMIN]: 'Administrador',
  [UserRole.ADMIN]: 'Administrador',
  [UserRole.OPERATOR]: 'Operador',
  [UserRole.ADMIN_OPERATOR]: 'Admin / Operador',
};

export function getRoleLabel(role?: string | null): string {
  if (!role) {
    return 'Usuario';
  }
  return roleLabels[role] ?? role;
}
