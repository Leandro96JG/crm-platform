package database

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/icrxz/crm-api-core/internal/domain"
	"github.com/jmoiron/sqlx"
)

type printJobRepository struct {
	client *sqlx.DB
}

func NewPrintJobRepository(client *sqlx.DB) domain.PrintJobRepository {
	return &printJobRepository{client: client}
}

func (r *printJobRepository) Create(ctx context.Context, job domain.PrintJob) (string, error) {
	dto := mapPrintJobToDTO(job)

	_, err := executor(ctx, r.client).NamedExecContext(
		ctx,
		`INSERT INTO print_jobs
		(job_id, order_item_id, job_type, status, queue_position, file_path, notes, copies, started_at, completed_at, failed_reason, created_by, created_at, updated_at)
		VALUES
		(:job_id, :order_item_id, :job_type, :status, :queue_position, :file_path, :notes, :copies, :started_at, :completed_at, :failed_reason, :created_by, :created_at, :updated_at)`,
		dto,
	)
	if err != nil {
		return "", err
	}

	return job.JobID, nil
}

func (r *printJobRepository) GetByID(ctx context.Context, jobID string) (*domain.PrintJob, error) {
	if jobID == "" {
		return nil, domain.NewValidationError("jobID is required", nil)
	}

	var dto PrintJobDTO
	err := executor(ctx, r.client).GetContext(ctx, &dto, "SELECT * FROM print_jobs WHERE job_id=$1", jobID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, domain.NewNotFoundError("no print job found with this id", map[string]any{"job_id": jobID})
		}
		return nil, err
	}

	job := mapDTOToPrintJob(dto)
	return &job, nil
}

func (r *printJobRepository) Search(ctx context.Context, filters domain.PrintJobFilters) (domain.PagingResult[domain.PrintJob], error) {
	whereQuery := []string{"1=1"}
	whereArgs := make([]any, 0)

	whereQuery, whereArgs = prepareInQuery(filters.OrderItemID, whereQuery, whereArgs, "order_item_id")
	whereQuery, whereArgs = prepareInQuery(filters.JobType, whereQuery, whereArgs, "job_type")
	whereQuery, whereArgs = prepareInQuery(filters.Status, whereQuery, whereArgs, "status")

	limitQuery := fmt.Sprintf("LIMIT $%d OFFSET $%d", len(whereArgs)+1, len(whereArgs)+2)
	limitArgs := append(whereArgs, filters.Limit, filters.Offset)

	orderBy := buildOrderBy(filters.SortBy, filters.SortOrder, validPrintJobSortColumns)

	query := fmt.Sprintf("SELECT * FROM print_jobs WHERE %s ORDER BY queue_position ASC, %s %s", joinWhere(whereQuery), orderBy, limitQuery)
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM print_jobs WHERE %s", joinWhere(whereQuery))

	var dtos []PrintJobDTO
	err := executor(ctx, r.client).SelectContext(ctx, &dtos, query, limitArgs...)
	if err != nil {
		return domain.PagingResult[domain.PrintJob]{}, err
	}

	var count int
	err = executor(ctx, r.client).GetContext(ctx, &count, countQuery, whereArgs...)
	if err != nil {
		return domain.PagingResult[domain.PrintJob]{}, err
	}

	return domain.PagingResult[domain.PrintJob]{
		Result: mapDTOsToPrintJobs(dtos),
		Paging: domain.Paging{Total: count, Limit: filters.Limit, Offset: filters.Offset},
	}, nil
}

func (r *printJobRepository) UpdateStatus(ctx context.Context, jobID string, status domain.PrintJobStatus, updatedBy string) error {
	_, err := executor(ctx, r.client).ExecContext(ctx,
		"UPDATE print_jobs SET status=$1, updated_at=NOW() WHERE job_id=$2",
		string(status), jobID,
	)
	return err
}

var validPrintJobSortColumns = map[string]bool{
	"created_at":     true,
	"updated_at":     true,
	"queue_position": true,
	"status":         true,
	"job_type":       true,
}
