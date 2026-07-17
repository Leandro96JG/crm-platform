package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type PrintJobRepository interface {
	Create(ctx context.Context, job PrintJob) (string, error)
	GetByID(ctx context.Context, jobID string) (*PrintJob, error)
	Search(ctx context.Context, filters PrintJobFilters) (PagingResult[PrintJob], error)
	UpdateStatus(ctx context.Context, jobID string, status PrintJobStatus, updatedBy string) error
}

type PrintJobStatus string

const (
	PrintJobQueued   PrintJobStatus = "queued"
	PrintJobPrinting PrintJobStatus = "printing"
	PrintJobPrinted  PrintJobStatus = "printed"
	PrintJobCutting  PrintJobStatus = "cutting"
	PrintJobCut      PrintJobStatus = "cut"
	PrintJobFailed   PrintJobStatus = "failed"
)

type PrintJobType string

const (
	PrintJobTypePrint PrintJobType = "print"
	PrintJobTypeCut   PrintJobType = "cut"
)

type PrintJob struct {
	JobID         string
	OrderItemID   string
	JobType       PrintJobType
	Status        PrintJobStatus
	QueuePosition int
	FilePath      string
	Notes         string
	Copies        int
	StartedAt     *time.Time
	CompletedAt   *time.Time
	FailedReason  string
	CreatedBy     string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type PrintJobFilters struct {
	OrderItemID []string
	JobType     []string
	Status      []string
	PagingFilter
}

func NewPrintJob(
	orderItemID string,
	jobType PrintJobType,
	filePath string,
	notes string,
	copies int,
	queuePosition int,
	createdBy string,
) (PrintJob, error) {
	jobID, err := uuid.NewUUID()
	if err != nil {
		return PrintJob{}, err
	}

	return PrintJob{
		JobID:         jobID.String(),
		OrderItemID:   orderItemID,
		JobType:       jobType,
		Status:        PrintJobQueued,
		QueuePosition: queuePosition,
		FilePath:      filePath,
		Notes:         notes,
		Copies:        copies,
		CreatedBy:     createdBy,
		CreatedAt:     time.Now().UTC(),
		UpdatedAt:     time.Now().UTC(),
	}, nil
}

func (j *PrintJob) Start() {
	now := time.Now().UTC()
	j.Status = PrintJobPrinting
	j.StartedAt = &now
	j.UpdatedAt = now
}

func (j *PrintJob) Complete() {
	now := time.Now().UTC()
	j.Status = PrintJobCut
	j.CompletedAt = &now
	j.UpdatedAt = now
}

func (j *PrintJob) Fail(reason string) {
	now := time.Now().UTC()
	j.Status = PrintJobFailed
	j.FailedReason = reason
	j.UpdatedAt = now
}
