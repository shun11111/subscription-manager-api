package handler

import (
	"net/http"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"subscription-manager-api/internal/interface/http/middleware"
	"subscription-manager-api/internal/usecase"
)

// SubscriptionHandler はサブスクリプションハンドラ
type SubscriptionHandler struct {
	subscriptionService *usecase.SubscriptionService
}

// NewSubscriptionHandler は新しいサブスクリプションハンドラを作成
func NewSubscriptionHandler(subscriptionService *usecase.SubscriptionService) *SubscriptionHandler {
	return &SubscriptionHandler{
		subscriptionService: subscriptionService,
	}
}

// ListSubscriptions はサブスクリプション一覧を取得
func (h *SubscriptionHandler) ListSubscriptions(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	subscriptions, err := h.subscriptionService.ListSubscriptions(c.Request().Context(), userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusOK, subscriptions)
}

// CreateSubscriptionRequest はサブスクリプション作成リクエスト
type CreateSubscriptionRequest struct {
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	BillingCycle    string  `json:"billing_cycle"`
	NextBillingDate string  `json:"next_billing_date"`
	Description     *string `json:"description,omitempty"`
}

// CreateSubscription はサブスクリプションを作成
func (h *SubscriptionHandler) CreateSubscription(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	var req CreateSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	serviceReq := &usecase.CreateSubscriptionRequest{
		Name:            req.Name,
		Price:           req.Price,
		BillingCycle:    req.BillingCycle,
		NextBillingDate: req.NextBillingDate,
		Description:     req.Description,
	}

	subscription, err := h.subscriptionService.CreateSubscription(c.Request().Context(), userID, serviceReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}

	return c.JSON(http.StatusCreated, subscription)
}

// UpdateSubscriptionRequest はサブスクリプション更新リクエスト
type UpdateSubscriptionRequest struct {
	Name            string  `json:"name"`
	Price           float64 `json:"price"`
	BillingCycle    string  `json:"billing_cycle"`
	NextBillingDate string  `json:"next_billing_date"`
	Description     *string `json:"description,omitempty"`
}

// GetSubscription は単一サブスクリプションを取得
func (h *SubscriptionHandler) GetSubscription(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid id"})
	}

	sub, err := h.subscriptionService.GetSubscription(c.Request().Context(), userID, id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	if sub == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "not found"})
	}

	return c.JSON(http.StatusOK, sub)
}

// UpdateSubscription はサブスクリプションを更新
func (h *SubscriptionHandler) UpdateSubscription(c echo.Context) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"message": "Unauthorized"})
	}

	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid id"})
	}

	var req UpdateSubscriptionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "Invalid request"})
	}

	serviceReq := &usecase.UpdateSubscriptionRequest{
		Name:            req.Name,
		Price:           req.Price,
		BillingCycle:    req.BillingCycle,
		NextBillingDate: req.NextBillingDate,
		Description:     req.Description,
	}

	sub, err := h.subscriptionService.UpdateSubscription(c.Request().Context(), userID, id, serviceReq)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"message": err.Error()})
	}
	if sub == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"message": "not found"})
	}

	return c.JSON(http.StatusOK, sub)
}

