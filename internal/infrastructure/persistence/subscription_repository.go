package persistence

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"subscription-manager-api/internal/domain"
)

// SubscriptionRepository はサブスクリプションリポジトリ
type SubscriptionRepository struct {
	db *pgxpool.Pool
}

// NewSubscriptionRepository は新しいサブスクリプションリポジトリを作成
func NewSubscriptionRepository(db *pgxpool.Pool) *SubscriptionRepository {
	return &SubscriptionRepository{db: db}
}

// Create はサブスクリプションを作成
func (r *SubscriptionRepository) Create(ctx context.Context, subscription *domain.Subscription) error {
	subscription.BeforeCreate()
	query := `
		INSERT INTO subscriptions (id, user_id, name, price, billing_cycle, next_billing_date, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	_, err := r.db.Exec(ctx, query,
		subscription.ID, subscription.UserID, subscription.Name,
		subscription.Price, subscription.BillingCycle, subscription.NextBillingDate,
		subscription.Description, subscription.CreatedAt, subscription.UpdatedAt,
	)
	return err
}

// FindByID はIDでサブスクリプションを取得
func (r *SubscriptionRepository) FindByID(ctx context.Context, id uuid.UUID) (*domain.Subscription, error) {
	var sub domain.Subscription
	query := `
		SELECT id, user_id, name, price, billing_cycle, next_billing_date::text, description, created_at, updated_at
		FROM subscriptions WHERE id = $1
	`
	err := r.db.QueryRow(ctx, query, id).Scan(
		&sub.ID, &sub.UserID, &sub.Name, &sub.Price,
		&sub.BillingCycle, &sub.NextBillingDate, &sub.Description,
		&sub.CreatedAt, &sub.UpdatedAt,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &sub, nil
}

// FindByUserID はユーザーIDでサブスクリプション一覧を取得
func (r *SubscriptionRepository) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*domain.Subscription, error) {
	query := `
		SELECT id, user_id, name, price, billing_cycle, next_billing_date::text, description, created_at, updated_at
		FROM subscriptions WHERE user_id = $1 ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscriptions []*domain.Subscription
	for rows.Next() {
		var sub domain.Subscription
		if err := rows.Scan(
			&sub.ID, &sub.UserID, &sub.Name, &sub.Price,
			&sub.BillingCycle, &sub.NextBillingDate, &sub.Description,
			&sub.CreatedAt, &sub.UpdatedAt,
		); err != nil {
			return nil, err
		}
		subscriptions = append(subscriptions, &sub)
	}
	return subscriptions, rows.Err()
}

// Update はサブスクリプションを更新
func (r *SubscriptionRepository) Update(ctx context.Context, subscription *domain.Subscription) error {
	subscription.BeforeUpdate()
	query := `
		UPDATE subscriptions
		SET name = $1, price = $2, billing_cycle = $3, next_billing_date = $4, description = $5, updated_at = $6
		WHERE id = $7 AND user_id = $8
	`
	result, err := r.db.Exec(ctx, query,
		subscription.Name, subscription.Price, subscription.BillingCycle,
		subscription.NextBillingDate, subscription.Description, subscription.UpdatedAt,
		subscription.ID, subscription.UserID,
	)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

// Delete はサブスクリプションを削除
func (r *SubscriptionRepository) Delete(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	query := `DELETE FROM subscriptions WHERE id = $1 AND user_id = $2`
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return pgx.ErrNoRows
	}
	return nil
}

