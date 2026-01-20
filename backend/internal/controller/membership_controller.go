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

type MembershipController struct {
	membershipService service.MembershipService
}

func NewMembershipController(membershipService service.MembershipService) *MembershipController {
	return &MembershipController{membershipService: membershipService}
}

func (m *MembershipController) CreateMembership(ctx *gin.Context) {
	var req *dto.CreateMembershipRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	membership, err := m.membershipService.CreateMembership(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusCreated, "Membership created successfully", membership)
}

func (m *MembershipController) GetMembershipByID(ctx *gin.Context) {
	membershipID := ctx.Param("membership_id")
	if membershipID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Membership ID is required", apperror.ErrBadRequest.Code)
		return
	}

	membership, err := m.membershipService.GetMembershipByID(membershipID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Membership retrieved successfully", membership)
}

func (m *MembershipController) GetMembershipByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "User ID is required", apperror.ErrBadRequest.Code)
		return
	}

	memberships, err := m.membershipService.GetMembershipsByUserID(userID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Memberships retrieved successfully", memberships)
}

func (m *MembershipController) GetAllMemberships(ctx *gin.Context) {
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

	response, err := m.membershipService.GetAllMemberships(page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Memberships retrieved successfully", response)
}

func (m *MembershipController) GetMembershipByCommunityID(ctx *gin.Context) {
	communityID := ctx.Param("community_id")
	if communityID == "" {
		dto.SendError(ctx, http.StatusBadRequest, "Community ID is required", apperror.ErrBadRequest.Code)
		return
	}

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

	response, err := m.membershipService.GetMembershipByCommunityID(communityID, page, pageSize)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Memberships retrieved successfully", response)
}

func (m *MembershipController) DeleteMembership(ctx *gin.Context) {
	var req *dto.DeleteMembershipRequest
	if err := ctx.ShouldBind(&req); err != nil {
		dto.SendError(ctx, http.StatusBadRequest, apperror.Message(apperror.ErrBadRequest), apperror.ErrBadRequest.Code)
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		dto.SendError(ctx, http.StatusForbidden, apperror.ErrForbidden.Message, apperror.ErrForbidden.Code)
		return
	}

	err := m.membershipService.DeleteMembership(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		dto.SendError(ctx, apperror.StatusFromError(err), apperror.Message(err), apperror.Code(err))
		return
	}

	dto.SendSuccess(ctx, http.StatusOK, "Membership deleted successfully", gin.H{"community_id": req.CommunityID, "user_id": authUser.(auth.AuthUser).ID})
}
