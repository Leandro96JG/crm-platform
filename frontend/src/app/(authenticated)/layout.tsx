import Snackbar from '../components/common/snackbar';
import SideNav from '../components/sidebar/sidenav';
import { SnackbarProvider } from '../context/SnackbarProvider';
import { getCurrentUser } from '../libs/session';
import { redirect } from 'next/navigation';

interface LayoutProps {
  children: React.ReactNode;
}

export default async function Layout({ children }: LayoutProps) {
  const user = await getCurrentUser();
  if (!user) {
    redirect('/login');
  }

  const displayName = user.name ?? user.username ?? 'Usuario';

  return (
    <SnackbarProvider>
      <div className="grid min-h-screen grid-cols-[76px_1fr] bg-paper lg:grid-cols-[248px_1fr]">
        <SideNav userRole={user.role} displayName={displayName} role={user.role} />

        <div className="flex min-w-0 flex-col">{children}</div>

        <Snackbar />
      </div>
    </SnackbarProvider>
  );
}
