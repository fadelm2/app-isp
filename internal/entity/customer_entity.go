package entity

type Customer struct {
	ID             string          `gorm:"column:id;primaryKey"`
	RegistrationID *string         `gorm:"column:registration_id"`
	UserID         string          `gorm:"column:user_id"`
	Status         string          `gorm:"column:status"` // active, suspended, terminated
	PackageID      string          `gorm:"column:package_id"`
	RouterID       *string         `gorm:"column:router_id"`
	PppUsername    string          `gorm:"column:ppp_username"`
	PppPassword    string          `gorm:"column:ppp_password"`
	RadiusUsername string          `gorm:"column:radius_username"`
	RadiusPassword string          `gorm:"column:radius_password"`
	DueDateDay     int             `gorm:"column:due_date_day"`
	CreatedAt      int64           `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt      int64           `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	User           User            `gorm:"foreignKey:user_id;references:id;->"`
	Package        InternetPackage `gorm:"foreignKey:package_id;references:id;->"`
	Router         *Router         `gorm:"foreignKey:router_id;references:id"`
}

func (c *Customer) TableName() string {
	return "customers"
}
