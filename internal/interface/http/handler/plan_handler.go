package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"subscription-manager-api/internal/usecase"
)

// PlanHandler はプランハンドラ
type PlanHandler struct {
	planService *usecase.PlanService
}

// NewPlanHandler は新しいプランハンドラを作成
func NewPlanHandler(planService *usecase.PlanService) *PlanHandler {
	return &PlanHandler{
		planService: planService,
	}
}

// ListPlans はプラン一覧を取得
func (h *PlanHandler) ListPlans(c echo.Context) error {
	plans, err := h.planService.ListPlans(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, plans)
}

// CreatePlanRequest はプラン作成リクエスト
type CreatePlanRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Price       float64  `json:"price"`
	BillingCycle string  `json:"billing_cycle"`
	Features    []string `json:"features,omitempty"`
	IsActive    bool     `json:"is_active"`
}

// CreatePlan はプランを作成
func (h *PlanHandler) CreatePlan(c echo.Context) error {
	var req CreatePlanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	serviceReq := &usecase.CreatePlanRequest{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		BillingCycle: req.BillingCycle,
		Features:     req.Features,
		IsActive:     req.IsActive,
	}

	plan, err := h.planService.CreatePlan(c.Request().Context(), serviceReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, plan)
}

// UpdatePlanRequest はプラン更新リクエスト
type UpdatePlanRequest struct {
	Name        string   `json:"name"`
	Description *string  `json:"description,omitempty"`
	Price       float64  `json:"price"`
	BillingCycle string  `json:"billing_cycle"`
	Features    []string `json:"features,omitempty"`
	IsActive    bool     `json:"is_active"`
}

// GetPlan は単一プランを取得
func (h *PlanHandler) GetPlan(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid id"})
	}

	plan, err := h.planService.GetPlan(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	if plan == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "not found"})
	}

	return c.JSON(http.StatusOK, plan)
}

// UpdatePlan はプランを更新
func (h *PlanHandler) UpdatePlan(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid id"})
	}

	var req UpdatePlanRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	serviceReq := &usecase.UpdatePlanRequest{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		BillingCycle: req.BillingCycle,
		Features:     req.Features,
		IsActive:     req.IsActive,
	}

	plan, err := h.planService.UpdatePlan(c.Request().Context(), id, serviceReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	if plan == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "not found"})
	}

	return c.JSON(http.StatusOK, plan)
}

// DeletePlan はプランを削除
func (h *PlanHandler) DeletePlan(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid id"})
	}

	if err := h.planService.DeletePlan(c.Request().Context(), id); err != nil {
		// pgx.ErrNoRowsの場合は404を返す
		if err.Error() == "no rows in result set" {
			return c.JSON(http.StatusNotFound, map[string]string{"message": "not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}

