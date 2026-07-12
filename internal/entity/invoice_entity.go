package entity

type Invoice struct {
	ID              string   `gorm:"column:id;primaryKey"`
	CustomerID      string   `gorm:"column:customer_id"`
	DueDate         int64    `gorm:"column:due_date"`
	PeriodMonth     int      `gorm:"column:period_month"`
	PeriodYear      int      `gorm:"column:period_year"`
	Amount          float64  `gorm:"column:amount"`
	TaxAmount       float64  `gorm:"column:tax_amount"`
	InstallationFee float64  `gorm:"column:installation_fee"`
	TotalAmount     float64  `gorm:"column:total_amount"`
	Status          string   `gorm:"column:status"` // pending, paid, owed, expired, cancelled
	SnapToken       string   `gorm:"column:snap_token"`
	PaidAt          *int64   `gorm:"column:paid_at"`
	CreatedAt       int64    `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt       int64    `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	Customer        Customer `gorm:"foreignKey:customer_id;references:id"`
}

func (i *Invoice) TableName() string {
	return "invoices"
}
