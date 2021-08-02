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

func (r *repository) FindByIds(userIds []int, active bool) (map[int]*models.User, error) {
	result := make(map[int]*models.User)
	if userIds == nil || len(userIds) == 0 {
		return result, nil
	}

	var users []*models.User
	err := r.Where("id IN ? AND active = ?", userIds, active).Find(&users).Error
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		result[int(user.ID)] = user
	}

	return result, nil
}

func (r *repository) FindByEmails(emails []string, active bool) (map[string]*models.User, error) {
	result := make(map[string]*models.User)
	if emails == nil || len(emails) == 0 {
		return result, nil
	}

	var users []*models.User
	err := r.Where("email IN ? AND active = ?", emails, active).Find(&users).Error
	if err != nil {
		return nil, err
	}

	for _, user := range users {
		result[user.Email] = user
	}

	return result, nil
}
