package usecase

import (
	"context"

	"github.com/google/uuid"
	"subscription-manager-api/internal/domain"
	"subscription-manager-api/internal/infrastructure/persistence"
)

// PlanService はプランユースケース
type PlanService struct {
	planRepo *persistence.PlanRepository
}

// NewPlanService は新しいプランユースケースを作成
func NewPlanService(planRepo *persistence.PlanRepository) *PlanService {
	return &PlanService{
		planRepo: planRepo,
	}
}

// ListPlans はプラン一覧を取得
func (s *PlanService) ListPlans(ctx context.Context) ([]*domain.Plan, error) {
	return s.planRepo.FindAll(ctx)
}

// CreatePlanRequest はプラン作成リクエスト
type CreatePlanRequest struct {
	Name        string
	Description *string
	Price       float64
	BillingCycle string
	Features    []string
	IsActive    bool
}

// CreatePlan はプランを作成
func (s *PlanService) CreatePlan(ctx context.Context, req *CreatePlanRequest) (*domain.Plan, error) {
	plan := &domain.Plan{
		Name:         req.Name,
		Description:  req.Description,
		Price:        req.Price,
		BillingCycle: domain.BillingCycle(req.BillingCycle),
		Features:     req.Features,
		IsActive:     req.IsActive,
	}

	if err := s.planRepo.Create(ctx, plan); err != nil {
		return nil, err
	}

	return plan, nil
}

// GetPlan は単一プランを取得
func (s *PlanService) GetPlan(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	plan, err := s.planRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return plan, nil
}

// UpdatePlanRequest はプラン更新リクエスト
type UpdatePlanRequest struct {
	Name        string
	Description *string
	Price       float64
	BillingCycle string
	Features    []string
	IsActive    bool
}

// UpdatePlan はプランを更新
func (s *PlanService) UpdatePlan(ctx context.Context, id uuid.UUID, req *UpdatePlanRequest) (*domain.Plan, error) {
	plan, err := s.planRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if plan == nil {
		return nil, nil
	}

	plan.Name = req.Name
	plan.Description = req.Description
	plan.Price = req.Price
	plan.BillingCycle = domain.BillingCycle(req.BillingCycle)
	plan.Features = req.Features
	plan.IsActive = req.IsActive

	if err := s.planRepo.Update(ctx, plan); err != nil {
		return nil, err
	}
	return plan, nil
}

// DeletePlan はプランを削除
func (s *PlanService) DeletePlan(ctx context.Context, id uuid.UUID) error {
	return s.planRepo.Delete(ctx, id)
}

