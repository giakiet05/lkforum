package controller

import (
	"net/http"
	"strconv"

	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

type ReportController struct {
	reportService service.ReportService
}

func NewReportController(reportService service.ReportService) *ReportController {
	return &ReportController{reportService: reportService}
}

func (rc *ReportController) GetReportByID(ctx *gin.Context) {
	reportID := ctx.Param("report_id")
	if reportID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Report ID is required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	report, err := rc.reportService.GetReportByID(reportID, authUser.(auth.AuthUser).ID, false)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report retrieved successfully", dto.FromReport(report))
}

func (rc *ReportController) GetReportByReporterID(ctx *gin.Context) {
	pageStr := ctx.DefaultQuery("page", "1")
	pageSizeStr := ctx.DefaultQuery("page_size", "10")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	responses, err := rc.reportService.GetReportByReporterID(authUser.(auth.AuthUser).ID, page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Reports retrieved successfully", responses)
}

func (rc *ReportController) GetReportByIDAdmin(ctx *gin.Context) {
	reportID := ctx.Param("report_id")
	if reportID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Report ID is required", apperror.ErrBadRequest.Code)
		return
	}

	report, err := rc.reportService.GetReportByID(reportID, "", true)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report retrieved successfully", dto.FromReport(report))
}

func (rc *ReportController) GetReportsFilter(ctx *gin.Context) {
	var query *dto.GetReportsFilterQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	responses, err := rc.reportService.GetReportsFilter(query)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Reports retrieved successfully", responses)
}

func (rc *ReportController) DeleteReportByID(ctx *gin.Context) {
	reportID := ctx.Param("report_id")
	if reportID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Report ID is required", apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := rc.reportService.DeleteReport(reportID, authUser.(auth.AuthUser).ID, false)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report deleted successfully", gin.H{"id": reportID})
}

func (rc *ReportController) DeleteReportByIDAdmin(ctx *gin.Context) {
	reportID := ctx.Param("report_id")
	if reportID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Report ID is required", apperror.ErrBadRequest.Code)
		return
	}

	err := rc.reportService.DeleteReport(reportID, "", true)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Report deleted successfully", gin.H{"id": reportID})
}

func (rc *ReportController) DeleteReportsByID(ctx *gin.Context) {
	var req *dto.DeleteReportsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := rc.reportService.DeleteReports(req.ReportIDs, authUser.(auth.AuthUser).ID, false)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Reports deleted successfully", gin.H{})
}

func (rc *ReportController) DeleteReportsByIDAdmin(ctx *gin.Context) {
	var req *dto.DeleteReportsRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	err := rc.reportService.DeleteReports(req.ReportIDs, "", true)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Reports deleted successfully", gin.H{})
}
