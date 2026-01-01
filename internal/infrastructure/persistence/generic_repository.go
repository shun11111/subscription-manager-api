package persistence

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Entity は全てのエンティティが実装すべきインターフェース
type Entity interface {
	GetID() uuid.UUID
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
	SetCreatedAt(time.Time)
	SetUpdatedAt(time.Time)
}

// Repository はジェネリックなリポジトリ
type Repository[T Entity] struct {
	db        *pgxpool.Pool
	tableName string
}

// NewRepository は新しいリポジトリを作成
func NewRepository[T Entity](db *pgxpool.Pool, tableName string) *Repository[T] {
	return &Repository[T]{
		db:        db,
		tableName: tableName,
	}
}

// FindByID はIDでエンティティを取得
func (r *Repository[T]) FindByID(ctx context.Context, id uuid.UUID) (*T, error) {
	var entity T
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", r.tableName)
	err := r.db.QueryRow(ctx, query, id).Scan(&entity)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &entity, nil
}

// FindByUserID はユーザーIDでエンティティを取得（ユーザー所有のリソース用）
func (r *Repository[T]) FindByUserID(ctx context.Context, userID uuid.UUID) ([]*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s WHERE user_id = $1 ORDER BY created_at DESC", r.tableName)
	rows, err := r.db.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*T
	for rows.Next() {
		var entity T
		if err := rows.Scan(&entity); err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}
	return entities, rows.Err()
}

// FindAll は全てのエンティティを取得（管理者用）
func (r *Repository[T]) FindAll(ctx context.Context) ([]*T, error) {
	query := fmt.Sprintf("SELECT * FROM %s ORDER BY created_at DESC", r.tableName)
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entities []*T
	for rows.Next() {
		var entity T
		if err := rows.Scan(&entity); err != nil {
			return nil, err
		}
		entities = append(entities, &entity)
	}
	return entities, rows.Err()
}

// Create はエンティティを作成
func (r *Repository[T]) Create(ctx context.Context, entity *T) error {
	// エンティティがBeforeCreateメソッドを持っている場合、呼び出す
	if bc, ok := any(entity).(interface{ BeforeCreate() }); ok {
		bc.BeforeCreate()
	}

	// リフレクションを使わず、各リポジトリで個別に実装する方が実用的
	// ここでは基本的な実装のみ提供
	return fmt.Errorf("Create method must be implemented in specific repository")
}

// Update はエンティティを更新
func (r *Repository[T]) Update(ctx context.Context, entity *T) error {
	if bc, ok := any(entity).(interface{ BeforeUpdate() }); ok {
		bc.BeforeUpdate()
	}
	return fmt.Errorf("Update method must be implemented in specific repository")
}

// Delete はエンティティを削除
func (r *Repository[T]) Delete(ctx context.Context, id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", r.tableName)
	_, err := r.db.Exec(ctx, query, id)
	return err
}

// DeleteByUserID はユーザーIDでエンティティを削除
func (r *Repository[T]) DeleteByUserID(ctx context.Context, userID uuid.UUID, id uuid.UUID) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1 AND user_id = $2", r.tableName)
	result, err := r.db.Exec(ctx, query, id, userID)
	if err != nil {
		return err
	}
	if result.RowsAffected() == 0 {
		return sql.ErrNoRows
	}
	return nil
}

