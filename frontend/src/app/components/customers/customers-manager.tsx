'use client';

import { Button } from '@/app/components/common/button';
import { ErrorMessage } from '@/app/components/common/error-message';
import Modal from '@/app/components/common/modal';
import StatusBadge from '@/app/components/common/status-badge';
import { useSnackbar } from '@/app/context/SnackbarProvider';
import {
  createCustomer,
  deleteCustomer,
  updateCustomer,
} from '@/app/services/customers';
import { Customer } from '@/app/types/customer';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { signOut } from 'next-auth/react';

interface CustomersManagerProps {
  customers: Customer[];
  total: number;
  initialSearch: string;
}

interface FormState {
  name: string;
  phone: string;
  email: string;
  document: string;
  address: string;
  notes: string;
}

const emptyForm: FormState = {
  name: '',
  phone: '',
  email: '',
  document: '',
  address: '',
  notes: '',
};

const inputClass =
  'w-full rounded-md border border-line bg-card px-3 py-2 text-[13px] text-ink outline-none transition-colors focus:border-cut-dark';
const labelClass = 'mb-1 block text-[12px] font-medium text-ink-soft';

export default function CustomersManager({
  customers,
  total,
  initialSearch,
}: CustomersManagerProps) {
  const router = useRouter();
  const { showSnackbar } = useSnackbar();

  const [search, setSearch] = useState(initialSearch);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<Customer | null>(null);
  const [form, setForm] = useState<FormState>(emptyForm);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);
  const [deletingId, setDeletingId] = useState<string | null>(null);

  function openCreate() {
    setEditing(null);
    setForm(emptyForm);
    setError(null);
    setModalOpen(true);
  }

  function openEdit(c: Customer) {
    setEditing(c);
    setForm({
      name: c.name,
      phone: c.phone,
      email: c.email,
      document: c.document,
      address: c.address,
      notes: c.notes,
    });
    setError(null);
    setModalOpen(true);
  }

  function handleSearch(e: React.FormEvent) {
    e.preventDefault();
    const params = search.trim() ? `?query=${encodeURIComponent(search.trim())}` : '';
    router.push(`/customers${params}`);
  }

  function unauthorized() {
    signOut({ callbackUrl: '/login' });
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);
    if (!form.name.trim()) {
      setError('El nombre es obligatorio');
      return;
    }
    setSaving(true);
    const resp = editing
      ? await updateCustomer(editing.customer_id, form)
      : await createCustomer(form);
    setSaving(false);

    if (!resp.success) {
      if (resp.unauthorized) return unauthorized();
      setError(resp.message);
      return;
    }

    showSnackbar(
      editing ? 'Cliente actualizado' : 'Cliente creado',
      'success'
    );
    setModalOpen(false);
    router.refresh();
  }

  async function handleDelete(c: Customer) {
    if (!confirm(`¿Eliminar al cliente "${c.name}"?`)) return;
    setDeletingId(c.customer_id);
    const resp = await deleteCustomer(c.customer_id);
    setDeletingId(null);
    if (!resp.success) {
      if (resp.unauthorized) return unauthorized();
      showSnackbar(resp.message, 'error');
      return;
    }
    showSnackbar('Cliente eliminado', 'success');
    router.refresh();
  }

  return (
    <>
      <div className="mb-4 flex flex-col gap-3 sm:flex-row sm:items-center sm:justify-between">
        <form onSubmit={handleSearch} className="flex gap-2">
          <input
            className={`${inputClass} sm:w-72`}
            value={search}
            onChange={(e) => setSearch(e.target.value)}
            placeholder="Buscar por nombre, teléfono o email"
          />
          <button
            type="submit"
            className="rounded-md border border-line px-3 py-2 text-[12px] font-semibold text-ink-soft transition-colors hover:bg-paper-dim"
          >
            Buscar
          </button>
        </form>
        <button
          type="button"
          onClick={openCreate}
          className="rounded-lg bg-cut-dark px-4 py-2 text-sm font-semibold text-white transition-colors hover:opacity-90"
        >
          + Nuevo cliente
        </button>
      </div>

      <div className="overflow-hidden rounded-lg bg-card shadow-card">
        {customers.length === 0 ? (
          <div className="p-10 text-center text-sm text-ink-soft">
            No hay clientes disponibles
          </div>
        ) : (
          <div className="overflow-x-auto">
            <table className="w-full border-collapse">
              <thead>
                <tr>
                  {['Nombre', 'Teléfono', 'Email', 'Documento', 'Estado'].map(
                    (h) => (
                      <th
                        key={h}
                        className="border-y border-line bg-paper-dim px-5 py-2.5 text-left text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint"
                      >
                        {h}
                      </th>
                    )
                  )}
                  <th className="border-y border-line bg-paper-dim px-5 py-2.5 text-right text-[10.5px] font-semibold uppercase tracking-wider text-ink-faint">
                    Acciones
                  </th>
                </tr>
              </thead>
              <tbody>
                {customers.map((c, i) => (
                  <tr
                    key={c.customer_id}
                    style={{ animationDelay: `${i * 40}ms` }}
                    className="animate-fadeIn border-b border-line transition-colors duration-150 last:border-none hover:bg-paper-dim/40"
                  >
                    <td className="px-5 py-3.5 text-[13px] font-semibold text-ink">
                      {c.name}
                    </td>
                    <td className="px-5 py-3.5 text-[13px] text-ink-soft">
                      {c.phone || '—'}
                    </td>
                    <td className="px-5 py-3.5 text-[13px] text-ink-soft">
                      {c.email || '—'}
                    </td>
                    <td className="px-5 py-3.5 text-[13px] text-ink-soft">
                      {c.document || '—'}
                    </td>
                    <td className="px-5 py-3.5">
                      <StatusBadge
                        label={c.is_active ? 'Activo' : 'Inactivo'}
                        className={
                          c.is_active
                            ? 'bg-st-ok-bg text-st-ok-fg'
                            : 'bg-paper-dim text-ink-soft'
                        }
                      />
                    </td>
                    <td className="px-5 py-3.5 text-right">
                      <div className="inline-flex gap-2">
                        <button
                          type="button"
                          onClick={() => openEdit(c)}
                          className="rounded-md border border-line px-2.5 py-1.5 text-[12px] font-medium text-ink-soft transition-colors hover:text-ink"
                        >
                          Editar
                        </button>
                        <button
                          type="button"
                          disabled={deletingId === c.customer_id}
                          onClick={() => handleDelete(c)}
                          className="rounded-md border border-line px-2.5 py-1.5 text-[12px] font-medium text-st-danger-fg transition-colors hover:bg-paper-dim disabled:opacity-50"
                        >
                          {deletingId === c.customer_id ? '…' : 'Eliminar'}
                        </button>
                      </div>
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
        <div className="border-t border-line px-5 py-3 text-[12px] text-ink-soft">
          Total: {total} clientes
        </div>
      </div>

      <Modal isOpen={modalOpen} onClose={() => setModalOpen(false)}>
        <form onSubmit={handleSubmit} className="w-[min(90vw,520px)]">
          <h2 className="mb-4 text-[16px] font-bold text-ink">
            {editing ? 'Editar cliente' : 'Nuevo cliente'}
          </h2>
          <div className="grid grid-cols-1 gap-4 sm:grid-cols-2">
            <div className="sm:col-span-2">
              <label className={labelClass}>Nombre *</label>
              <input
                className={inputClass}
                value={form.name}
                onChange={(e) => setForm({ ...form, name: e.target.value })}
                placeholder="Nombre del cliente"
              />
            </div>
            <div>
              <label className={labelClass}>Teléfono</label>
              <input
                className={inputClass}
                value={form.phone}
                onChange={(e) => setForm({ ...form, phone: e.target.value })}
              />
            </div>
            <div>
              <label className={labelClass}>Email</label>
              <input
                className={inputClass}
                value={form.email}
                onChange={(e) => setForm({ ...form, email: e.target.value })}
              />
            </div>
            <div>
              <label className={labelClass}>Documento</label>
              <input
                className={inputClass}
                value={form.document}
                onChange={(e) => setForm({ ...form, document: e.target.value })}
              />
            </div>
            <div>
              <label className={labelClass}>Dirección</label>
              <input
                className={inputClass}
                value={form.address}
                onChange={(e) => setForm({ ...form, address: e.target.value })}
              />
            </div>
            <div className="sm:col-span-2">
              <label className={labelClass}>Notas</label>
              <textarea
                className={inputClass}
                rows={2}
                value={form.notes}
                onChange={(e) => setForm({ ...form, notes: e.target.value })}
              />
            </div>
          </div>

          {error && <ErrorMessage message={error} />}

          <div className="mt-5 flex items-center justify-end gap-3">
            <button
              type="button"
              onClick={() => setModalOpen(false)}
              className="rounded-lg border border-line px-4 py-2 text-sm font-medium text-ink-soft transition-colors hover:bg-paper-dim"
            >
              Cancelar
            </button>
            <Button type="submit" color="success" isLoading={saving}>
              {editing ? 'Guardar' : 'Crear'}
            </Button>
          </div>
        </form>
      </Modal>
    </>
  );
}
