package domain

// Plan はプラン・価格マスターエンティティ
type Plan struct {
	BaseEntity
	Name        string       `json:"name" db:"name"`
	Description *string      `json:"description,omitempty" db:"description"`
	Price       float64      `json:"price" db:"price"`
	BillingCycle BillingCycle `json:"billing_cycle" db:"billing_cycle"`
	Features    []string     `json:"features,omitempty" db:"features"` // PostgreSQLの配列型として保存
	IsActive    bool         `json:"is_active" db:"is_active"`
}

