package controllers

import (
	"github.com/badaccuracyid/softeng_backend/src/database/dao"
	"github.com/badaccuracyid/softeng_backend/src/model"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserController interface {
	SetContext(ctx *gin.Context)
	GetUserByID(id string) (*model.User, error)
	GetUsersByID(ids []string) ([]*model.User, error)
	UpdateUser(user model.User) (*model.User, error)
	CreateUser(user model.User) (*model.User, error)
	Login(email string, password string) (*model.User, error)
}

type userController struct {
	ctx     *gin.Context
	userDAO *dao.UserDAO
}

func NewUserService(db *gorm.DB) UserController {
	return &userController{
		userDAO: dao.NewUserDAO(db),
	}
}

func (s *userController) SetContext(ctx *gin.Context) {
	s.ctx = ctx
}

func (s *userController) GetUserByID(id string) (*model.User, error) {
	user, err := s.userDAO.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userController) GetUsersByID(ids []string) ([]*model.User, error) {
	users, err := s.userDAO.GetUsersByID(ids)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (s *userController) CreateUser(user model.User) (*model.User, error) {
	err := s.userDAO.CreateUser(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userController) UpdateUser(user model.User) (*model.User, error) {
	err := s.userDAO.UpdateUser(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *userController) Login(email string, password string) (*model.User, error) {
	user, err := s.userDAO.GetUserByEmail(email)
	if err != nil {
		return nil, err
	}

	if user.Password != password {
		return nil, nil
	}

	return user, nil
}
