package user

import (
	"errors"
	"github.com/vietanhduong/ota-server/pkg/mysql"
	"github.com/vietanhduong/ota-server/pkg/mysql/models"
	"gorm.io/gorm"
)

type repository struct {
	*mysql.DB
}

func NewRepository(db *mysql.DB) *repository {
	return &repository{db}
}

func (r *repository) FindById(userId uint) (*models.User, error) {
	var model *models.User
	err := r.Where("id = ?", userId).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return model, err
}

func (r *repository) FindByEmail(email string, active bool) (*models.User, error) {
	var model *models.User
	err := r.Where("email = ? AND active = ?", email, active).First(&model).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return model, err
}
