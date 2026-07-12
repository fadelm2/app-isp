package entity

type CustomerHistory struct {
	ID         string   `gorm:"column:id;primaryKey"`
	CustomerID string   `gorm:"column:customer_id"`
	Action     string   `gorm:"column:action"` // register, activate, suspend, unsuspend, terminate
	Notes      string   `gorm:"column:notes"`
	CreatedBy  string   `gorm:"column:created_by"`
	CreatedAt  int64    `gorm:"column:created_at;autoCreateTime:milli"`
	Customer   Customer `gorm:"foreignKey:customer_id;references:id"`
	User       User     `gorm:"foreignKey:created_by;references:id"`
}

func (h *CustomerHistory) TableName() string {
	return "customer_histories"
}
