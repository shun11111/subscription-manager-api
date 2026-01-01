package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"subscription-manager-api/internal/domain"
)

// PlanRepository はプランリポジトリ
type PlanRepository struct {
	db *pgxpool.Pool
}

// NewPlanRepository は新しいプランリポジトリを作成
func NewPlanRepository(db *pgxpool.Pool) *PlanRepository {
	return &PlanRepository{db: db}
}

// Create はプランを作成
func (r *PlanRepository) Create(ctx context.Context, plan *domain.Plan) error {
	plan.BeforeCreate()
	query := `
		INSERT INTO plans (id, name, description, price, billing_cycle, features, is_active, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		plan.ID, plan.Name, plan.Description,
		plan.Price, plan.BillingCycle, plan.Features,
		plan.IsActive, plan.CreatedAt, plan.UpdatedAt,
	)
	return err
}

// FindByID はIDでプランを取得
func (r *PlanRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Plan, error) {
	var plan domain.Plan
	query := `
		SELECT id, name, description, price, billing_cycle, COALESCE(features, ARRAY[]::TEXT[]), is_active, created_at, updated_at
		FROM plans WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&plan.ID, &plan.Name, &plan.Description, &plan.Price,
		&plan.BillingCycle, &plan.Features, &plan.IsActive,
		&plan.CreatedAt, &plan.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &plan, nil
}

// FindAll は全プラン一覧を取得
func (r *PlanRepository) FindAll(ctx context.Context) ([]*domain.Plan, error) {
	query := `
		SELECT id, name, description, price, billing_cycle, COALESCE(features, ARRAY[]::TEXT[]), is_active, created_at, updated_at
		FROM plans ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var plans []*domain.Plan
	for rows.Next() {
		var plan domain.Plan
		if err := rows.Scan(
			&plan.ID, &plan.Name, &plan.Description, &plan.Price,
			&plan.BillingCycle, &plan.Features, &plan.IsActive,
			&plan.CreatedAt, &plan.UpdatedAt,
		); err != nil {
			return nil, err
		}
		plans = append(plans, &plan)
	}
	return plans, rows.Err()
}

// Update はプランを更新
func (r *PlanRepository) Update(ctx context.Context, plan *domain.Plan) error {
	plan.BeforeUpdate()
	query := `
		UPDATE plans
		SET name = $1, description = $2, price = $3, billing_cycle = $4, features = $5, is_active = $6, updated_at = $7
		WHERE id = $8
	`
	result, err := r.db.Exec(ctx, query,
		plan.Name, plan.Description, plan.Price, plan.BillingCycle,
		plan.Features, plan.IsActive, plan.UpdatedAt,
		plan.ID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Delete はプランを削除
func (r *PlanRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM plans WHERE id = $1`
	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

