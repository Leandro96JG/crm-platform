'use client';
import { UserRole } from '@/app/types/user';
import { getInitials, getRoleLabel } from '@/app/utils/user-display';
import { ArrowLeftStartOnRectangleIcon } from '@heroicons/react/24/outline';
import { signOut } from 'next-auth/react';
import Image from 'next/image';
import Link from 'next/link';
import NavLinks from './nav-links';

interface SideNavProps {
  userRole: UserRole;
  displayName: string;
  role: string;
}

export default function SideNav({ userRole, displayName, role }: SideNavProps) {
  return (
    <aside className="sticky top-0 flex h-screen flex-col bg-sidebar px-3 py-5 text-sidebar-text lg:px-4">
      <Link
        href="/home"
        className="mb-6 flex items-center justify-center gap-2.5 px-2 lg:justify-start"
      >
        <div className="relative h-9 w-9 flex-none lg:h-10 lg:w-10">
          <Image
            src="/viva-logo.png"
            alt="Viva"
            fill
            sizes="40px"
            className="object-contain"
            priority
          />
        </div>
        <span className="hidden text-[15.5px] font-bold leading-tight text-white lg:block">
          Viva
        </span>
      </Link>

      <div className="mb-2 hidden px-2.5 text-[10.5px] uppercase tracking-wider text-ink-faint lg:block">
        Menú
      </div>

      <nav className="mb-auto flex flex-col gap-0.5">
        <NavLinks userRole={userRole} />
      </nav>

      <div className="mt-3.5 border-t border-white/10 pt-3.5">
        <div className="flex items-center justify-center gap-2.5 rounded-sm p-2 lg:justify-start">
          <div className="flex h-8 w-8 flex-none items-center justify-center rounded-full bg-teal text-[12.5px] font-bold text-white">
            {getInitials(displayName)}
          </div>
          <div className="hidden min-w-0 flex-1 lg:block">
            <div className="truncate text-[13px] font-semibold text-white">
              {displayName}
            </div>
            <div className="text-[11px] text-ink-faint">
              {getRoleLabel(role)}
            </div>
          </div>
        </div>
        <button
          type="button"
          onClick={() => signOut({ callbackUrl: '/login' })}
          className="mt-1.5 flex w-full items-center justify-center gap-2 rounded-sm px-2.5 py-2 text-[12.5px] text-ink-faint transition-colors hover:bg-sidebar-soft hover:text-white lg:justify-start"
        >
          <ArrowLeftStartOnRectangleIcon className="h-4 w-4" />
          <span className="hidden lg:block">Cerrar sesión</span>
        </button>
      </div>
    </aside>
  );
}
