CREATE TABLE IF NOT EXISTS print_jobs (
    job_id TEXT PRIMARY KEY,
    order_item_id TEXT NOT NULL REFERENCES order_items(item_id),
    job_type TEXT NOT NULL DEFAULT 'print',
    status TEXT NOT NULL DEFAULT 'queued',
    queue_position INT NOT NULL DEFAULT 0,
    file_path TEXT DEFAULT '',
    notes TEXT DEFAULT '',
    copies INT NOT NULL DEFAULT 1,
    started_at TIMESTAMP,
    completed_at TIMESTAMP,
    failed_reason TEXT DEFAULT '',
    created_by TEXT NOT NULL DEFAULT '',
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    updated_at TIMESTAMP NOT NULL DEFAULT now()
);

CREATE INDEX IF NOT EXISTS idx_print_jobs_status ON print_jobs (status);
CREATE INDEX IF NOT EXISTS idx_print_jobs_type ON print_jobs (job_type);
CREATE INDEX IF NOT EXISTS idx_print_jobs_order ON print_jobs (order_item_id);
