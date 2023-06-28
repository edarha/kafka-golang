package repositories

import (
	"context"

	"github.com/edarha/kafka-golang/internals/models"
	"gorm.io/gorm"
)

type Student interface {
	Create(ctx context.Context, entity *models.Student) error
}

type studentRepo struct {
	db *gorm.DB
}

func NewStudentRepo(db *gorm.DB) Student {
	return &studentRepo{
		db: db,
	}
}

func (r *studentRepo) Create(ctx context.Context, entity *models.Student) error {
	return r.db.WithContext(ctx).Model(&models.Student{}).Create(entity).Error
}