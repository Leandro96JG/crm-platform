package database

import (
	"time"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type PrintJobDTO struct {
	JobID         string     `db:"job_id"`
	OrderItemID   string     `db:"order_item_id"`
	JobType       string     `db:"job_type"`
	Status        string     `db:"status"`
	QueuePosition int        `db:"queue_position"`
	FilePath      string     `db:"file_path"`
	Notes         string     `db:"notes"`
	Copies        int        `db:"copies"`
	StartedAt     *time.Time `db:"started_at"`
	CompletedAt   *time.Time `db:"completed_at"`
	FailedReason  string     `db:"failed_reason"`
	CreatedBy     string     `db:"created_by"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}

func mapPrintJobToDTO(j domain.PrintJob) PrintJobDTO {
	return PrintJobDTO{
		JobID:         j.JobID,
		OrderItemID:   j.OrderItemID,
		JobType:       string(j.JobType),
		Status:        string(j.Status),
		QueuePosition: j.QueuePosition,
		FilePath:      j.FilePath,
		Notes:         j.Notes,
		Copies:        j.Copies,
		StartedAt:     j.StartedAt,
		CompletedAt:   j.CompletedAt,
		FailedReason:  j.FailedReason,
		CreatedBy:     j.CreatedBy,
		CreatedAt:     j.CreatedAt,
		UpdatedAt:     j.UpdatedAt,
	}
}

func mapDTOToPrintJob(dto PrintJobDTO) domain.PrintJob {
	return domain.PrintJob{
		JobID:         dto.JobID,
		OrderItemID:   dto.OrderItemID,
		JobType:       domain.PrintJobType(dto.JobType),
		Status:        domain.PrintJobStatus(dto.Status),
		QueuePosition: dto.QueuePosition,
		FilePath:      dto.FilePath,
		Notes:         dto.Notes,
		Copies:        dto.Copies,
		StartedAt:     dto.StartedAt,
		CompletedAt:   dto.CompletedAt,
		FailedReason:  dto.FailedReason,
		CreatedBy:     dto.CreatedBy,
		CreatedAt:     dto.CreatedAt,
		UpdatedAt:     dto.UpdatedAt,
	}
}

func mapDTOsToPrintJobs(dtos []PrintJobDTO) []domain.PrintJob {
	jobs := make([]domain.PrintJob, 0, len(dtos))
	for _, dto := range dtos {
		jobs = append(jobs, mapDTOToPrintJob(dto))
	}
	return jobs
}
