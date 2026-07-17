'use client';

import { Button } from '@/app/components/common/button';
import { ErrorMessage } from '@/app/components/common/error-message';
import { useSnackbar } from '@/app/context/SnackbarProvider';
import { createOrder } from '@/app/services/orders';
import { calculatePrice } from '@/app/services/planchas';
import { Plancha, StickerMaterial } from '@/app/types/plancha';
import { Customer } from '@/app/types/customer';
import { formatCurrency } from '@/app/utils/status';
import { useRouter } from 'next/navigation';
import { useState } from 'react';
import { signOut } from 'next-auth/react';

interface OrderFormProps {
  planchas: Plancha[];
  materials: StickerMaterial[];
  customers: Customer[];
}

interface ItemRow {
  key: number;
  plancha_id: string;
  material_id: string;
  sheet_quantity: number;
  unit_price: number;
  calculating: boolean;
}

const emptyItem = (key: number): ItemRow => ({
  key,
  plancha_id: '',
  material_id: '',
  sheet_quantity: 1,
  unit_price: 0,
  calculating: false,
});

const inputClass =
  'w-full rounded-md border border-line bg-card px-3 py-2 text-[13px] text-ink outline-none transition-colors focus:border-cut-dark';
const labelClass = 'mb-1 block text-[12px] font-medium text-ink-soft';

export default function OrderForm({
  planchas,
  materials,
  customers,
}: OrderFormProps) {
  const router = useRouter();
  const { showSnackbar } = useSnackbar();

  const [customerId, setCustomerId] = useState('');
  const [urgency, setUrgency] = useState('normal');
  const [notes, setNotes] = useState('');
  const [items, setItems] = useState<ItemRow[]>([emptyItem(0)]);
  const [nextKey, setNextKey] = useState(1);
  const [error, setError] = useState<string | null>(null);
  const [saving, setSaving] = useState(false);

  const total = items.reduce(
    (sum, it) => sum + it.unit_price * it.sheet_quantity,
    0
  );

  function addItem() {
    setItems((prev) => [...prev, emptyItem(nextKey)]);
    setNextKey((k) => k + 1);
  }

  function removeItem(key: number) {
    setItems((prev) => prev.filter((it) => it.key !== key));
  }

  async function recalc(row: ItemRow) {
    if (!row.plancha_id || !row.material_id || row.sheet_quantity < 1) {
      return;
    }
    setItems((prev) =>
      prev.map((it) =>
        it.key === row.key ? { ...it, calculating: true } : it
      )
    );
    const resp = await calculatePrice(
      row.plancha_id,
      row.material_id,
      row.sheet_quantity
    );
    setItems((prev) =>
      prev.map((it) => {
        if (it.key !== row.key) return it;
        const unit =
          resp.success && resp.data
            ? resp.data.total / row.sheet_quantity
            : 0;
        return { ...it, unit_price: unit, calculating: false };
      })
    );
    if (!resp.success) {
      setError(resp.message);
    }
  }

  function updateItem(key: number, patch: Partial<ItemRow>) {
    setItems((prev) => {
      const updated = prev.map((it) =>
        it.key === key ? { ...it, ...patch } : it
      );
      const row = updated.find((it) => it.key === key);
      if (row) recalc(row);
      return updated;
    });
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError(null);

    if (!customerId.trim()) {
      setError('El cliente es obligatorio');
      return;
    }
    const validItems = items.filter(
      (it) => it.plancha_id && it.material_id && it.sheet_quantity >= 1
    );
    if (validItems.length === 0) {
      setError('Agregá al menos un ítem con plancha, material y cantidad');
      return;
    }

    setSaving(true);
    const resp = await createOrder({
      customer_id: customerId,
      source: 'manual',
      notes,
      urgency,
      items: validItems.map((it) => ({
        plancha_id: it.plancha_id,
        material_id: it.material_id,
        sheet_quantity: it.sheet_quantity,
        unit_price: it.unit_price,
        custom_design_file: '',
        custom_design_notes: '',
      })),
    });
    setSaving(false);

    if (!resp.success) {
      if (resp.unauthorized) {
        signOut({ callbackUrl: '/login' });
        return;
      }
      setError(resp.message);
      return;
    }

    showSnackbar(
      `Pedido ${resp.data?.order_number ?? ''} creado`,
      'success'
    );
    router.push('/orders');
    router.refresh();
  }

  return (
    <form onSubmit={handleSubmit} className="flex flex-col gap-6">
      <div className="rounded-lg bg-card p-5 shadow-card">
        <h2 className="mb-4 text-[15px] font-bold text-ink">Datos del pedido</h2>
        <div className="grid grid-cols-1 gap-4 md:grid-cols-2">
          <div>
            <label className={labelClass}>Cliente *</label>
            <input
              className={inputClass}
              list="customers-list"
              value={customerId}
              onChange={(e) => setCustomerId(e.target.value)}
              placeholder="Elegí un cliente o escribí uno nuevo"
            />
            <datalist id="customers-list">
              {customers.map((c) => (
                <option key={c.customer_id} value={c.name}>
                  {c.phone ? `${c.name} — ${c.phone}` : c.name}
                </option>
              ))}
            </datalist>
          </div>
          <div>
            <label className={labelClass}>Urgencia</label>
            <select
              className={inputClass}
              value={urgency}
              onChange={(e) => setUrgency(e.target.value)}
            >
              <option value="normal">Normal</option>
              <option value="urgent">Urgente</option>
            </select>
          </div>
          <div className="md:col-span-2">
            <label className={labelClass}>Notas</label>
            <textarea
              className={inputClass}
              value={notes}
              onChange={(e) => setNotes(e.target.value)}
              rows={2}
              placeholder="Observaciones del pedido (opcional)"
            />
          </div>
        </div>
      </div>

      <div className="rounded-lg bg-card p-5 shadow-card">
        <div className="mb-4 flex items-center justify-between">
          <h2 className="text-[15px] font-bold text-ink">Ítems</h2>
          <button
            type="button"
            onClick={addItem}
            className="rounded-md border border-line px-3 py-1.5 text-[12px] font-semibold text-cut-dark transition-colors hover:bg-paper-dim"
          >
            + Agregar ítem
          </button>
        </div>

        <div className="flex flex-col gap-3">
          {items.map((it) => (
            <div
              key={it.key}
              className="grid grid-cols-1 items-end gap-3 rounded-md border border-line p-3 md:grid-cols-[1fr_1fr_90px_120px_auto]"
            >
              <div>
                <label className={labelClass}>Plancha</label>
                <select
                  className={inputClass}
                  value={it.plancha_id}
                  onChange={(e) =>
                    updateItem(it.key, { plancha_id: e.target.value })
                  }
                >
                  <option value="">Seleccionar…</option>
                  {planchas.map((p) => (
                    <option key={p.plancha_id} value={p.plancha_id}>
                      {p.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className={labelClass}>Material</label>
                <select
                  className={inputClass}
                  value={it.material_id}
                  onChange={(e) =>
                    updateItem(it.key, { material_id: e.target.value })
                  }
                >
                  <option value="">Seleccionar…</option>
                  {materials.map((m) => (
                    <option key={m.material_id} value={m.material_id}>
                      {m.name}
                    </option>
                  ))}
                </select>
              </div>
              <div>
                <label className={labelClass}>Cant.</label>
                <input
                  type="number"
                  min={1}
                  className={inputClass}
                  value={it.sheet_quantity}
                  onChange={(e) =>
                    updateItem(it.key, {
                      sheet_quantity: Math.max(1, Number(e.target.value) || 1),
                    })
                  }
                />
              </div>
              <div>
                <label className={labelClass}>Subtotal</label>
                <div className="py-2 font-mono text-[13px] font-semibold text-ink">
                  {it.calculating
                    ? '…'
                    : formatCurrency(it.unit_price * it.sheet_quantity)}
                </div>
              </div>
              <div>
                {items.length > 1 && (
                  <button
                    type="button"
                    onClick={() => removeItem(it.key)}
                    className="rounded-md border border-line px-2.5 py-2 text-[12px] text-st-danger-fg transition-colors hover:bg-paper-dim"
                    aria-label="Quitar ítem"
                  >
                    Quitar
                  </button>
                )}
              </div>
            </div>
          ))}
        </div>

        <div className="mt-4 flex items-center justify-end gap-2 border-t border-line pt-4">
          <span className="text-[13px] text-ink-soft">Total:</span>
          <span className="font-mono text-xl font-bold text-ink">
            {formatCurrency(total)}
          </span>
        </div>
      </div>

      {error && <ErrorMessage message={error} />}

      <div className="flex items-center justify-end gap-3">
        <button
          type="button"
          onClick={() => router.push('/orders')}
          className="rounded-lg border border-line px-4 py-2 text-sm font-medium text-ink-soft transition-colors hover:bg-paper-dim"
        >
          Cancelar
        </button>
        <Button type="submit" color="success" isLoading={saving}>
          Crear pedido
        </Button>
      </div>
    </form>
  );
}
