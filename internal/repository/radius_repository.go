package repository

import (
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type RadiusRepository struct {
	Log *logrus.Logger
}

func NewRadiusRepository(log *logrus.Logger) *RadiusRepository {
	return &RadiusRepository{
		Log: log,
	}
}

func (r *RadiusRepository) CreateOrUpdateCheck(db *gorm.DB, check *entity.RadCheck) error {
	var existing entity.RadCheck
	err := db.Where("username = ? AND attribute = ?", check.Username, check.Attribute).First(&existing).Error
	if err == nil {
		existing.Value = check.Value
		existing.Op = check.Op
		return db.Save(&existing).Error
	}
	return db.Create(check).Error
}

func (r *RadiusRepository) DeleteCheck(db *gorm.DB, username string) error {
	return db.Where("username = ?", username).Delete(&entity.RadCheck{}).Error
}

func (r *RadiusRepository) CreateOrUpdateReply(db *gorm.DB, reply *entity.RadReply) error {
	var existing entity.RadReply
	err := db.Where("username = ? AND attribute = ?", reply.Username, reply.Attribute).First(&existing).Error
	if err == nil {
		existing.Value = reply.Value
		existing.Op = reply.Op
		return db.Save(&existing).Error
	}
	return db.Create(reply).Error
}

func (r *RadiusRepository) DeleteReply(db *gorm.DB, username string) error {
	return db.Where("username = ?", username).Delete(&entity.RadReply{}).Error
}

func (r *RadiusRepository) FindActiveSessions(db *gorm.DB) ([]entity.RadAcct, error) {
	var sessions []entity.RadAcct
	if err := db.Where("acctstoptime IS NULL").Find(&sessions).Error; err != nil {
		return nil, err
	}
	return sessions, nil
}
