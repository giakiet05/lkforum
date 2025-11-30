package service

import (
	"github.com/giakiet05/lkforum/internal/apperror"
	"github.com/giakiet05/lkforum/internal/dto"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
)

type ReportService interface {
	GetReportByID(reportID string, requesterID string, isAdmin bool) (*model.Report, error)
	GetReportByReporterID(reporterID string, page int, pageSize int) (*dto.PaginatedReportResponse, error)
	GetReportsFilter(req *dto.GetReportsFilterQuery) (*dto.PaginatedReportResponse, error)
	DeleteReport(reportID string, requesterID string, isAdmin bool) error
	DeleteReports(reportIDs []string, requesterID string, isAdmin bool) error
}

type reportService struct {
	reportRepo repo.ReportRepo
}

func NewReportService(reportRepo repo.ReportRepo) ReportService {
	return &reportService{reportRepo: reportRepo}
}

func (r *reportService) GetReportByID(reportID string, requesterID string, isAdmin bool) (*model.Report, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if !isAdmin {
		ok, err := r.reportRepo.IsReporter(ctx, requesterID, reportID)
		if err != nil {
			return nil, apperror.ErrInternal
		}
		if !ok {
			return nil, apperror.ErrForbidden
		}
	}

	report, err := r.reportRepo.GetByID(ctx, reportID)
	if err != nil {
		return nil, err
	}

	return report, nil
}

func (r *reportService) GetReportByReporterID(reporterID string, page int, pageSize int) (*dto.PaginatedReportResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	reports, total, err := r.reportRepo.GetFilter(
		ctx, &reporterID, nil, "", nil, nil, nil, page, pageSize,
	)
	if err != nil {
		return nil, err
	}
	reportResponses := dto.FromReports(reports)

	var response = &dto.PaginatedReportResponse{
		Reports: reportResponses,
		Pagination: dto.Pagination{
			Total:    total,
			Page:     page,
			PageSize: pageSize,
		},
	}

	return response, nil
}

func (r *reportService) GetReportsFilter(req *dto.GetReportsFilterQuery) (*dto.PaginatedReportResponse, error) {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	reports, total, err := r.reportRepo.GetFilter(
		ctx, req.ReporterID, req.TargetID, req.TargetType, req.Reason, req.StartDate, req.EndDate,
		req.Page, req.PageSize,
	)
	if err != nil {
		return nil, err
	}

	reportResponses := dto.FromReports(reports)
	var response = &dto.PaginatedReportResponse{
		Reports: reportResponses,
		Pagination: dto.Pagination{
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
	}

	return response, nil
}

func (r *reportService) DeleteReport(reportID string, requesterID string, isAdmin bool) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if !isAdmin {
		ok, err := r.reportRepo.IsReporter(ctx, requesterID, reportID)
		if err != nil {
			return apperror.ErrInternal
		}
		if !ok {
			return apperror.ErrForbidden
		}
	}

	return r.reportRepo.Delete(ctx, reportID)
}

func (r *reportService) DeleteReports(reportIDs []string, requesterID string, isAdmin bool) error {
	ctx, cancel := util.NewDefaultDBContext()
	defer cancel()

	if len(reportIDs) == 0 {
		return apperror.ErrBadRequest
	}

	if !isAdmin {
		ok, err := r.reportRepo.IsReporterOfAllReports(ctx, requesterID, reportIDs)
		if err != nil {
			return apperror.ErrInternal
		}
		if !ok {
			return apperror.ErrForbidden
		}
	}

	return r.reportRepo.DeleteBatch(ctx, reportIDs)
}
