package entity

type InternetPackage struct {
	ID              string  `gorm:"column:id;primaryKey"`
	Name            string  `gorm:"column:name"`
	SpeedMbps       int     `gorm:"column:speed_mbps"`
	Price           float64 `gorm:"column:price"`
	InstallationFee float64 `gorm:"column:installation_fee"`
	TaxRate         float64 `gorm:"column:tax_rate"`
	IsActive        bool    `gorm:"column:is_active"`
	CreatedAt       int64   `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt       int64   `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (p *InternetPackage) TableName() string {
	return "internet_packages"
}
