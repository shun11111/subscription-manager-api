package usecase

import (
	"context"
	"time"

	"github.com/google/uuid"
	"subscription-manager-api/internal/domain"
	"subscription-manager-api/internal/infrastructure/persistence"
)

// SubscriptionService はサブスクリプションユースケース
type SubscriptionService struct {
	subscriptionRepo *persistence.SubscriptionRepository
}

// NewSubscriptionService は新しいサブスクリプションユースケースを作成
func NewSubscriptionService(subscriptionRepo *persistence.SubscriptionRepository) *SubscriptionService {
	return &SubscriptionService{
		subscriptionRepo: subscriptionRepo,
	}
}

// ListSubscriptions はユーザーのサブスクリプション一覧を取得
func (s *SubscriptionService) ListSubscriptions(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error) {
	return s.subscriptionRepo.FindByUserID(ctx, userID)
}

// CreateSubscription はサブスクリプションを作成
func (s *SubscriptionService) CreateSubscription(ctx context.Context, userID uuid.UUID, req *CreateSubscriptionRequest) (*domain.Subscription, error) {
	subscription := &domain.Subscription{
		UserID:          userID.String(),
		Name:            req.Name,
		Price:           req.Price,
		BillingCycle:    domain.BillingCycle(req.BillingCycle),
		NextBillingDate: req.NextBillingDate,
		Description:     req.Description,
	}

	if err := s.subscriptionRepo.Create(ctx, subscription); err != nil {
		return nil, err
	}

	return subscription, nil
}

// CreateSubscriptionRequest はサブスクリプション作成リクエスト
type CreateSubscriptionRequest struct {
	Name            string
	Price           float64
	BillingCycle    string
	NextBillingDate string
	Description     *string
}

// UpdateSubscriptionRequest はサブスクリプション更新リクエスト
type UpdateSubscriptionRequest struct {
	Name            string
	Price           float64
	BillingCycle    string
	NextBillingDate string
	Description     *string
}

// GetSubscription は単一サブスクリプションを取得
func (s *SubscriptionService) GetSubscription(ctx context.Context, userID uuid.UUID, id uuid.UUID) (*domain.Subscription, error) {
	sub, err := s.subscriptionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, nil
	}
	// 所有者チェック
	if sub.UserID != userID.String() {
		return nil, nil
	}
	return sub, nil
}

// UpdateSubscription はサブスクリプションを更新
func (s *SubscriptionService) UpdateSubscription(ctx context.Context, userID uuid.UUID, id uuid.UUID, req *UpdateSubscriptionRequest) (*domain.Subscription, error) {
	sub, err := s.subscriptionRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if sub == nil {
		return nil, nil
	}
	if sub.UserID != userID.String() {
		return nil, nil
	}

	sub.Name = req.Name
	sub.Price = req.Price
	sub.BillingCycle = domain.BillingCycle(req.BillingCycle)
	sub.NextBillingDate = req.NextBillingDate
	sub.Description = req.Description

	if err := s.subscriptionRepo.Update(ctx, sub); err != nil {
		return nil, err
	}
	return sub, nil
}

// CalculateNextBillingDate は次回課金日を計算
func CalculateNextBillingDate(billingCycle domain.BillingCycle, currentDate time.Time) time.Time {
	switch billingCycle {
	case domain.BillingCycleMonthly:
		return currentDate.AddDate(0, 1, 0)
	case domain.BillingCycleYearly:
		return currentDate.AddDate(1, 0, 0)
	default:
		return currentDate.AddDate(0, 1, 0) // デフォルトは1ヶ月後
	}
}

