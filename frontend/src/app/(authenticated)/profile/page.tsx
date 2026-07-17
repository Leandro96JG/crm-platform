'use server';
import { ProfileForm } from '@/app/components/users/profile-form';
import Topbar from '@/app/components/common/topbar';
import { getCurrentUser } from '@/app/libs/session';
import { getUserByID } from '@/app/services/user';
import { getInitials } from '@/app/utils/user-display';
import { redirect } from 'next/navigation';

export default async function Page() {
  const session = await getCurrentUser();

  if (!session) {
    redirect('/login');
  }

  const { data: user } = await getUserByID(session.user_id);

  if (!user) {
    redirect('/login');
  }

  return (
    <>
      <Topbar
        title="Mi perfil"
        subtitle="Gestione sus datos de acceso"
        initials={getInitials(user.first_name ?? session.username)}
      />
      <div className="cut-divider" />
      <main className="px-8 pb-10">
        <div className="max-w-2xl rounded-lg bg-card p-6 shadow-card">
          <ProfileForm user={user} />
        </div>
      </main>
    </>
  );
}
