package repository

import (
	"github.com/sirupsen/logrus"
	"golang-clean-architecture/internal/entity"
	"gorm.io/gorm"
)

type PackageRepository struct {
	Repository[entity.InternetPackage]
	Log *logrus.Logger
}

func NewPackageRepository(log *logrus.Logger) *PackageRepository {
	return &PackageRepository{
		Log: log,
	}
}

func (r *PackageRepository) FindAll(db *gorm.DB) ([]entity.InternetPackage, error) {
	var packages []entity.InternetPackage
	if err := db.Find(&packages).Error; err != nil {
		return nil, err
	}
	return packages, nil
}
