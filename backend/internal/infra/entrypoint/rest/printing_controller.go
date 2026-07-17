package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/icrxz/crm-api-core/internal/application"
	"github.com/icrxz/crm-api-core/internal/domain"
)

type PrintingController struct {
	printingService application.PrintingService
}

type PrintJobDTO struct {
	JobID         string  `json:"job_id"`
	OrderItemID   string  `json:"order_item_id"`
	JobType       string  `json:"job_type"`
	Status        string  `json:"status"`
	QueuePosition int     `json:"queue_position"`
	FilePath      string  `json:"file_path"`
	Notes         string  `json:"notes"`
	Copies        int     `json:"copies"`
	CreatedBy     string  `json:"created_by"`
}

type UpdatePrintJobStatusDTO struct {
	Status    string `json:"status" validate:"required"`
	UpdatedBy string `json:"updated_by" validate:"required"`
}

func NewPrintingController(printingService application.PrintingService) PrintingController {
	return PrintingController{
		printingService: printingService,
	}
}

func (ctrl *PrintingController) GetPrintQueue(ctx *gin.Context) {
	jobType := ctx.DefaultQuery("job_type", "print")
	statusFilter := ctx.QueryArray("status")

	if len(statusFilter) == 0 {
		statusFilter = []string{"queued", "printing", "cutting"}
	}

	jobs, err := ctrl.printingService.GetPrintQueue(ctx.Request.Context(), jobType, statusFilter)
	if err != nil {
		ctx.Error(err)
		return
	}

	dtos := make([]PrintJobDTO, len(jobs.Result))
	for i, j := range jobs.Result {
		dtos[i] = PrintJobDTO{
			JobID:         j.JobID,
			OrderItemID:   j.OrderItemID,
			JobType:       string(j.JobType),
			Status:        string(j.Status),
			QueuePosition: j.QueuePosition,
			FilePath:      j.FilePath,
			Notes:         j.Notes,
			Copies:        j.Copies,
			CreatedBy:     j.CreatedBy,
		}
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": dtos,
		"paging": jobs.Paging,
	})
}

func (ctrl *PrintingController) UpdatePrintJobStatus(ctx *gin.Context) {
	jobID := ctx.Param("jobID")
	if jobID == "" {
		ctx.Error(domain.NewValidationError("job_id is required", nil))
		return
	}

	var dto UpdatePrintJobStatusDTO
	if err := ctx.BindJSON(&dto); err != nil {
		ctx.Error(err)
		return
	}

	err := ctrl.printingService.UpdatePrintJobStatus(
		ctx.Request.Context(),
		jobID,
		domain.PrintJobStatus(dto.Status),
		dto.UpdatedBy,
	)
	if err != nil {
		ctx.Error(err)
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
