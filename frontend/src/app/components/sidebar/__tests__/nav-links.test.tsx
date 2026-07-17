import { render, screen } from '@testing-library/react';
import { usePathname } from 'next/navigation';
import { UserRole } from '@/app/types/user';
import NavLinks from '../nav-links';

jest.mock('next/navigation', () => ({
  usePathname: jest.fn(),
}));

beforeEach(() => {
  jest.clearAllMocks();
  (usePathname as jest.Mock).mockReturnValue('/home');
});

describe('NavLinks', () => {
  it('should render links available to every role', () => {
    render(<NavLinks userRole={UserRole.OPERATOR} />);

    expect(screen.getByText('Home')).toBeInTheDocument();
    expect(screen.getByText('Planchas')).toBeInTheDocument();
    expect(screen.getByText('Pedidos')).toBeInTheDocument();
    expect(screen.getByText('Producción')).toBeInTheDocument();
    expect(screen.getByText('Mi Perfil')).toBeInTheDocument();
  });

  it('should hide admin-only links for an operator', () => {
    render(<NavLinks userRole={UserRole.OPERATOR} />);

    expect(screen.queryByText('Usuarios')).not.toBeInTheDocument();
  });

  it('should show admin-only links for an admin', () => {
    render(<NavLinks userRole={UserRole.ADMIN} />);

    expect(screen.getByText('Usuarios')).toBeInTheDocument();
  });

  it('should link the profile item to /profile for every role', () => {
    render(<NavLinks userRole={UserRole.OPERATOR} />);

    expect(screen.getByText('Mi Perfil').closest('a')).toHaveAttribute(
      'href',
      '/profile'
    );
  });
});
