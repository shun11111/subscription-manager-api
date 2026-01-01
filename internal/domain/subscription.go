package domain

// BillingCycle は課金サイクル
type BillingCycle string

const (
	BillingCycleMonthly BillingCycle = "monthly"
	BillingCycleYearly  BillingCycle = "yearly"
)

// Subscription はサブスクリプションエンティティ
type Subscription struct {
	BaseEntity
	UserID         string       `json:"user_id" db:"user_id"`
	Name           string       `json:"name" db:"name"`
	Price          float64      `json:"price" db:"price"`
	BillingCycle   BillingCycle `json:"billing_cycle" db:"billing_cycle"`
	NextBillingDate string       `json:"next_billing_date" db:"next_billing_date"` // date形式で保存
	Description    *string       `json:"description,omitempty" db:"description"`
}

