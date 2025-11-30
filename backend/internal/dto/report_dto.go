package dto

import (
	"time"

	"github.com/giakiet05/lkforum/internal/model"
)

type GetReportsFilterQuery struct {
	ReporterID *string                `form:"reporter_id,omitempty"`
	TargetID   *string                `form:"target_id,omitempty"`
	TargetType model.ReportTargetType `form:"target_type,omitempty"`
	Reason     *string                `form:"reason,omitempty"`
	StartDate  *time.Time             `form:"start_date,omitempty"`
	EndDate    *time.Time             `form:"end_date,omitempty"`
	Page       int                    `form:"page,default=1"`
	PageSize   int                    `form:"page_size,default=20"`
}

type DeleteReportsRequest struct {
	ReportIDs []string `json:"report_ids" binding:"required,min=1"`
}

type ReportResponse struct {
	ID          string                 `json:"id"`
	ReporterID  string                 `json:"reporter_id"`
	TargetID    string                 `json:"target_id"`
	TargetType  model.ReportTargetType `json:"target_type"`
	Reason      string                 `json:"reason"`
	Description *string                `json:"description,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

func FromReport(r *model.Report) *ReportResponse {
	return &ReportResponse{
		ID:          r.ID.Hex(),
		ReporterID:  r.ReporterID.Hex(),
		TargetID:    r.TargetID.Hex(),
		TargetType:  r.TargetType,
		Reason:      r.Reason,
		Description: r.Description,
		CreatedAt:   r.CreatedAt,
	}
}

func FromReports(reports []model.Report) []*ReportResponse {
	responses := make([]*ReportResponse, len(reports))
	for i, rep := range reports {
		responses[i] = FromReport(&rep)
	}
	return responses
}
