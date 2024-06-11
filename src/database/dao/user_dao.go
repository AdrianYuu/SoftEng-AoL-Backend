package dao

import (
	"github.com/badaccuracyid/softeng_backend/src/model"
	"gorm.io/gorm"
)

type UserDAO struct {
	DB *gorm.DB
}

func NewUserDAO(db *gorm.DB) *UserDAO {
	return &UserDAO{
		DB: db,
	}
}

func (dao *UserDAO) CreateUser(user *model.User) error {
	return dao.DB.Create(user).Error
}

func (dao *UserDAO) GetUserByID(id string) (*model.User, error) {
	user := &model.User{}
	err := dao.DB.First(user, "id = ?", id).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (dao *UserDAO) GetUserByEmail(email string) (*model.User, error) {
	user := &model.User{}
	err := dao.DB.First(user, "email = ?", email).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (dao *UserDAO) GetUsersByID(ids []string) ([]*model.User, error) {
	var users []*model.User
	err := dao.DB.Find(&users, "id IN ?", ids).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (dao *UserDAO) UpdateUser(user *model.User) error {
	return dao.DB.Save(user).Error
}

func (dao *UserDAO) DeleteUser(id string) error {
	return dao.DB.Delete(&model.User{}, "id = ?", id).Error
}
