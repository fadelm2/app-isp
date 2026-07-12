package repository

import (
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type RouterRepository struct {
	Repository[entity.Router]
	Log *logrus.Logger
}

func NewRouterRepository(log *logrus.Logger) *RouterRepository {
	return &RouterRepository{
		Log: log,
	}
}

func (r *RouterRepository) FindAll(db *gorm.DB) ([]entity.Router, error) {
	var routers []entity.Router
	if err := db.Find(&routers).Error; err != nil {
		return nil, err
	}
	return routers, nil
}
