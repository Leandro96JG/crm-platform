export interface Order {
  order_id: string;
  order_number: string;
  customer_id: string;
  status: OrderStatus;
  source: string;
  ai_handled: boolean;
  assigned_to: string;
  notes: string;
  total: number;
  urgency: string;
  completed_at: string | null;
  created_by: string;
  created_at: string;
  updated_by: string;
  updated_at: string;
  items: OrderItem[];
}

export type OrderStatus =
  | 'pending'
  | 'approved'
  | 'in_production'
  | 'ready'
  | 'delivered'
  | 'cancelled';

export interface OrderItem {
  item_id: string;
  plancha_id: string;
  material_id: string;
  sheet_quantity: number;
  unit_price: number;
  subtotal: number;
  custom_design_file: string;
  custom_design_notes: string;
}

export interface CreateOrderDTO {
  customer_id: string;
  source: string;
  notes: string;
  urgency: string;
  created_by: string;
  items: CreateOrderItemDTO[];
}

export interface CreateOrderItemDTO {
  plancha_id: string;
  material_id: string;
  sheet_quantity: number;
  unit_price: number;
  custom_design_file: string;
  custom_design_notes: string;
}
