export interface PrintJob {
  job_id: string;
  order_item_id: string;
  job_type: 'print' | 'cut';
  status: PrintJobStatus;
  queue_position: number;
  file_path: string;
  notes: string;
  copies: number;
  created_by: string;
}

export type PrintJobStatus =
  | 'queued'
  | 'printing'
  | 'printed'
  | 'cutting'
  | 'cut'
  | 'failed';
