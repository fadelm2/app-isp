package entity

type Router struct {
	ID        string `gorm:"column:id;primaryKey"`
	Name      string `gorm:"column:name"`
	Host      string `gorm:"column:host"`
	Port      int    `gorm:"column:port"`
	Username  string `gorm:"column:username"`
	Password  string `gorm:"column:password"`
	Status    string `gorm:"column:status"` // online, offline
	CreatedAt int64  `gorm:"column:created_at;autoCreateTime:milli"`
	UpdatedAt int64  `gorm:"column:updated_at;autoCreateTime:milli;autoUpdateTime:milli"`
}

func (r *Router) TableName() string {
	return "routers"
}
