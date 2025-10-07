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
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: err.Error()})
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Code: "FORBIDDEN", Error: "Invalid jwt token"})
		return
	}

	membership, err := m.membershipService.CreateMembership(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: ""})
		return
	}

	ctx.JSON(http.StatusCreated, dto.SuccessResponse{
		ID:      membership.ID.Hex(),
		Message: "Create membership successfully",
	})
}

func (m *MembershipController) GetMembershipByID(ctx *gin.Context) {
	membershipID := ctx.Param("membership_id")
	if membershipID == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: membershipID})
		return
	}

	membership, err := m.membershipService.GetMembershipByID(membershipID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: ""})
		return
	}

	ctx.JSON(http.StatusOK, membership)
}

func (m *MembershipController) GetMembershipByUserID(ctx *gin.Context) {
	userID := ctx.Param("user_id")
	if userID == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: userID})
		return
	}

	memberships, err := m.membershipService.GetMembershipByUserID(userID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: ""})
		return
	}

	ctx.JSON(http.StatusOK, memberships)
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

	memberships, err := m.membershipService.GetAllMemberships(page, pageSize)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: ""})
		return
	}

	ctx.JSON(http.StatusOK, memberships)
}

func (m *MembershipController) GetMembershipByCommunityID(ctx *gin.Context) {
	communityID := ctx.Param("community_id")
	if communityID == "" {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: communityID})
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
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: ""})
		return
	}

	ctx.JSON(http.StatusOK, response)
}

func (m *MembershipController) DeleteMembership(ctx *gin.Context) {
	var req *dto.DeleteMembershipRequest
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, dto.ErrorResponse{Code: "INVALID_REQUEST", Error: err.Error()})
		return
	}

	authUser, exists := ctx.Get("authUser")
	if !exists {
		ctx.JSON(http.StatusForbidden, dto.ErrorResponse{Code: "FORBIDDEN", Error: "Invalid jwt token"})
		return
	}

	err := m.membershipService.DeleteMembership(req, authUser.(auth.AuthUser).ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, dto.ErrorResponse{Code: apperror.Code(err), Error: ""})
		return
	}

	ctx.JSON(http.StatusNoContent, dto.SuccessResponse{
		ID:      req.CommunityID,
		Message: "Delete membership successfully",
	})
}
