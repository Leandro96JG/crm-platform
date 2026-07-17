package application

import (
	"context"
	"fmt"

	"github.com/icrxz/crm-api-core/internal/domain"
)

type printingService struct {
	printRepo  domain.PrintJobRepository
	orderRepo  domain.OrderRepository
	orderSvc   OrderService
}

type PrintingService interface {
	GetPrintQueue(ctx context.Context, jobType string, status []string) (domain.PagingResult[domain.PrintJob], error)
	UpdatePrintJobStatus(ctx context.Context, jobID string, status domain.PrintJobStatus, updatedBy string) error
	CreateManualPrintJob(ctx context.Context, job domain.PrintJob) (string, error)
}

func NewPrintingService(
	printRepo domain.PrintJobRepository,
	orderRepo domain.OrderRepository,
	orderSvc OrderService,
) PrintingService {
	return &printingService{
		printRepo: printRepo,
		orderRepo: orderRepo,
		orderSvc:  orderSvc,
	}
}

func (s *printingService) GetPrintQueue(ctx context.Context, jobType string, status []string) (domain.PagingResult[domain.PrintJob], error) {
	filters := domain.PrintJobFilters{
		JobType: []string{jobType},
		Status:  status,
		PagingFilter: domain.PagingFilter{
			Limit:  100,
			Offset: 0,
			SortBy: "queue_position",
		},
	}

	return s.printRepo.Search(ctx, filters)
}

func (s *printingService) UpdatePrintJobStatus(ctx context.Context, jobID string, status domain.PrintJobStatus, updatedBy string) error {
	job, err := s.printRepo.GetByID(ctx, jobID)
	if err != nil {
		return err
	}

	validTransition := isValidTransition(job.Status, status)
	if !validTransition {
		return domain.NewValidationError(
			fmt.Sprintf("invalid status transition from %s to %s", job.Status, status),
			map[string]any{"job_id": jobID, "current_status": job.Status, "new_status": status},
		)
	}

	return s.printRepo.UpdateStatus(ctx, jobID, status, updatedBy)
}

func (s *printingService) CreateManualPrintJob(ctx context.Context, job domain.PrintJob) (string, error) {
	return s.printRepo.Create(ctx, job)
}

func isValidTransition(current, new domain.PrintJobStatus) bool {
	transitions := map[domain.PrintJobStatus][]domain.PrintJobStatus{
		domain.PrintJobQueued:   {domain.PrintJobPrinting, domain.PrintJobCutting, domain.PrintJobFailed},
		domain.PrintJobPrinting: {domain.PrintJobPrinted, domain.PrintJobFailed},
		domain.PrintJobPrinted:  {domain.PrintJobCutting, domain.PrintJobFailed},
		domain.PrintJobCutting:  {domain.PrintJobCut, domain.PrintJobFailed},
		domain.PrintJobCut:      {},
		domain.PrintJobFailed:   {domain.PrintJobQueued},
	}

	allowed, ok := transitions[current]
	if !ok {
		return false
	}

	for _, a := range allowed {
		if a == new {
			return true
		}
	}
	return false
}
