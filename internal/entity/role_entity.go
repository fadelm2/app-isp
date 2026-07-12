package entity

import "gorm.io/gorm"

type Role struct {
	gorm.Model
	ID          uint   `gorm:"column:id;primaryKey"`
	Name        string `gorm:"column:name"`
	Description string `gorm:"column:description"`
}

func (u *Role) TableName() string {
	return "roles"
}
