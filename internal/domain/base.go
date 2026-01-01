package domain

import (
	"time"

	"github.com/google/uuid"
)

// BaseEntity は全てのエンティティが持つ共通フィールド
type BaseEntity struct {
	ID        uuid.UUID `json:"id" db:"id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// GetID はIDを取得
func (e *BaseEntity) GetID() uuid.UUID {
	return e.ID
}

// GetCreatedAt は作成日時を取得
func (e *BaseEntity) GetCreatedAt() time.Time {
	return e.CreatedAt
}

// GetUpdatedAt は更新日時を取得
func (e *BaseEntity) GetUpdatedAt() time.Time {
	return e.UpdatedAt
}

// SetCreatedAt は作成日時を設定
func (e *BaseEntity) SetCreatedAt(t time.Time) {
	e.CreatedAt = t
}

// SetUpdatedAt は更新日時を設定
func (e *BaseEntity) SetUpdatedAt(t time.Time) {
	e.UpdatedAt = t
}

// BeforeCreate は作成前に呼ばれるフック
func (e *BaseEntity) BeforeCreate() {
	now := time.Now()
	if e.ID == uuid.Nil {
		e.ID = uuid.New()
	}
	if e.CreatedAt.IsZero() {
		e.CreatedAt = now
	}
	if e.UpdatedAt.IsZero() {
		e.UpdatedAt = now
	}
}

// BeforeUpdate は更新前に呼ばれるフック
func (e *BaseEntity) BeforeUpdate() {
	e.UpdatedAt = time.Now()
}

