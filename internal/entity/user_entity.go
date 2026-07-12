package entity

// User is a struct that represents a user entity
type User struct {
	ID          string    `gorm:"column:id;primaryKey"`
	Password    string    `gorm:"column:password"`
	Name        string    `gorm:"column:username"`
	Email       string    `gorm:"column:email"`
	RoleID      string    `gorm:"column:role_id"`
	CompanyName string    `gorm:"column:company_name"`
	Token       string    `gorm:"-"`
	RegionId    string    `gorm:"column:region_id"`
	CreatedAt   int64     `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt   int64     `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
	Contacts    []Contact `gorm:"foreignKey:user_id;references:id"`
}

func (u *User) TableName() string {
	return "users"
}
