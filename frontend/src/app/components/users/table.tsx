'use client';
import { UserListItem } from '@/app/types/user-list-item';
import { UserRole } from '@/app/types/user';
import { SearchResponse } from '@/app/types/search_response';
import { roleLabels } from '@/app/utils/roles';
import StatusBadge from '@/app/components/common/status-badge';
import { EyeIcon } from '@heroicons/react/24/outline';
import { Pagination } from '@heroui/pagination';
import { useRouter } from 'next/navigation';

interface UsersTableProps {
  users?: SearchResponse<UserListItem>;
  initialPage?: number;
}

export default function UsersTable({
  users,
  initialPage = 1,
}: UsersTableProps) {
  const router = useRouter();

  function handleRowClick(userID: string) {
    router.push(`/users/${userID}`);
  }

  function handleChangePage(value: number) {
    router.push(`?page=${value}`);
  }

  return (
    <div className="w-full">
      <div className="overflow-hidden rounded-lg bg-card shadow-card">
        <div className="overflow-x-auto">
          <table className="w-full border-collapse">
            <thead>
              <tr>
                {['Nombre', 'Usuario', 'Correo', 'Cargo', 'Estado', ''].map(
                  (h, i) => (
                    <th
                      key={i}
                      className="border-y border-line bg-paper-dim px-5 py-2.5 text-left text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint"
                    >
                      {h}
                    </th>
                  )
                )}
              </tr>
            </thead>
            <tbody>
              {users?.result
                .filter((user) => user.role != UserRole.THAVANNA_ADMIN)
                .map((user, i) => (
                  <tr
                    key={user.user_id}
                    style={{ animationDelay: `${i * 40}ms` }}
                    className="animate-fadeIn border-b border-line transition-colors duration-150 last:border-none hover:bg-paper-dim/40"
                  >
                    <td className="px-5 py-3.5 text-[13px] font-semibold text-ink">
                      {`${user.first_name} ${user.last_name}`}
                    </td>
                    <td className="px-5 py-3.5 font-mono text-[12.5px] text-ink-soft">
                      {user.username}
                    </td>
                    <td className="px-5 py-3.5 text-[13px] text-ink-soft">
                      {user.email}
                    </td>
                    <td className="px-5 py-3.5 text-[13px] text-ink">
                      {roleLabels[user.role]}
                    </td>
                    <td className="px-5 py-3.5">
                      {user.active ? (
                        <StatusBadge
                          label="Activo"
                          className="bg-st-ok-bg text-st-ok-fg"
                        />
                      ) : (
                        <StatusBadge
                          label="Inactivo"
                          className="bg-paper-dim text-ink-soft"
                        />
                      )}
                    </td>
                    <td className="px-5 py-3.5">
                      <button
                        disabled={!user.active}
                        aria-label="Ver usuario"
                        className="text-ink-soft transition-colors hover:text-cut disabled:opacity-40"
                        onClick={() => handleRowClick(user.user_id)}
                      >
                        <EyeIcon className="h-5 w-5" />
                      </button>
                    </td>
                  </tr>
                ))}
            </tbody>
          </table>
        </div>
      </div>

      <div className="mt-4">
        <Pagination
          onChange={handleChangePage}
          siblings={3}
          showControls
          total={Math.ceil(
            Number((users?.paging.total || 1) / (users?.paging.limit || 1))
          )}
          page={Number(initialPage || 1)}
        />
      </div>
    </div>
  );
}
