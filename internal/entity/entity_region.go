package entity

type Regions struct {
	ID          string `gorm:"column:id;primaryKey"`
	Description string `gorm:"column:description"`
}

func (c *Regions) TableName() string {
	return "regions"
}
