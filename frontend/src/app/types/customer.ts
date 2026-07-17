export interface Customer {
  customer_id: string;
  name: string;
  phone: string;
  email: string;
  document: string;
  address: string;
  notes: string;
  is_active: boolean;
  created_by: string;
  created_at: string;
  updated_by: string;
  updated_at: string;
}

export interface CreateCustomerDTO {
  name: string;
  phone: string;
  email: string;
  document: string;
  address: string;
  notes: string;
}

export interface UpdateCustomerDTO {
  name?: string;
  phone?: string;
  email?: string;
  document?: string;
  address?: string;
  notes?: string;
  is_active?: boolean;
}
