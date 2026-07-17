'use client';
import { UserRole } from '@/app/types/user';
import { adminRoles } from '@/app/utils/roles';
import {
  GridIcon,
  OrderBoxIcon,
  GearIcon,
} from '@/app/components/common/icons';
import {
  HomeIcon,
  UserGroupIcon,
  UserCircleIcon,
} from '@heroicons/react/24/outline';
import clsx from 'clsx';
import Link from 'next/link';
import { usePathname } from 'next/navigation';
import { ComponentType, SVGProps } from 'react';

type IconType = ComponentType<SVGProps<SVGSVGElement>>;

const links: {
  name: string;
  href: string;
  icon: IconType;
  onlyAdmin: boolean;
}[] = [
  { name: 'Home', href: '/home', icon: HomeIcon, onlyAdmin: false },
  { name: 'Planchas', href: '/stickers', icon: GridIcon, onlyAdmin: false },
  { name: 'Pedidos', href: '/orders', icon: OrderBoxIcon, onlyAdmin: false },
  { name: 'Producción', href: '/printing', icon: GearIcon, onlyAdmin: false },
  { name: 'Usuarios', href: '/users', icon: UserGroupIcon, onlyAdmin: true },
  { name: 'Mi Perfil', href: '/profile', icon: UserCircleIcon, onlyAdmin: false },
];

export default function NavLinks({ userRole }: { userRole: UserRole }) {
  const pathname = usePathname();

  return (
    <>
      {links.map((link) => {
        if (link.onlyAdmin && !adminRoles.includes(userRole)) {
          return null;
        }

        const LinkIcon = link.icon;
        const active = pathname === link.href;
        return (
          <Link
            key={link.name}
            href={link.href}
            className={clsx(
              'relative flex items-center justify-center gap-3 rounded-sm px-2.5 py-2.5 text-[13.5px] font-medium transition-colors lg:justify-start',
              active
                ? 'bg-sidebar-soft text-white'
                : 'text-sidebar-text hover:bg-sidebar-soft'
            )}
          >
            {active && (
              <span className="absolute -left-3 bottom-2 top-2 w-[3px] rounded-r bg-cut lg:-left-4" />
            )}
            <LinkIcon className="h-[17px] w-[17px] flex-none opacity-90" />
            <span className="hidden lg:block">{link.name}</span>
          </Link>
        );
      })}
    </>
  );
}
