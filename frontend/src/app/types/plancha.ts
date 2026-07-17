export interface Plancha {
  plancha_id: string;
  name: string;
  description: string;
  category: string;
  subcategory: string;
  layout_file_url: string;
  preview_image_url: string;
  notes: string;
  is_active: boolean;
  created_by: string;
  created_at: string;
  updated_by: string;
  updated_at: string;
}

export interface StickerMaterial {
  material_id: string;
  name: string;
  description: string;
  material_type: string;
  finish: string;
  is_cuttable: boolean;
  is_printable: boolean;
  base_cost: number;
  is_active: boolean;
  created_at: string;
}

export interface PlanchaPrice {
  price_id: string;
  plancha_id: string;
  material_id: string;
  base_price: number;
  min_quantity: number;
  bulk_discount: BulkDiscount[];
  is_active: boolean;
}

export interface BulkDiscount {
  min_quantity: number;
  unit_price: number;
}

export interface CreatePlanchaDTO {
  name: string;
  description: string;
  category: string;
  subcategory: string;
  layout_file_url: string;
  notes: string;
  created_by: string;
}
