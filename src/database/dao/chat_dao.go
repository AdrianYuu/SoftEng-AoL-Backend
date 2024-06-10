package dao

import (
	"github.com/badaccuracyid/softeng_backend/src/model"
	"gorm.io/gorm"
)

type ChatDAO struct {
	DB *gorm.DB
}

func NewChatDAO(db *gorm.DB) *ChatDAO {
	return &ChatDAO{
		DB: db,
	}
}

func (dao *ChatDAO) CreateConversation(conversation *model.Conversation) error {
	return dao.DB.Create(conversation).Error
}

func (dao *ChatDAO) GetConversationByID(id string) (*model.Conversation, error) {
	conversation := &model.Conversation{}
	err := dao.DB.First(&conversation, id).Error
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func (dao *ChatDAO) UpdateConversation(conversation *model.Conversation) error {
	return dao.DB.Save(conversation).Error
}

func (dao *ChatDAO) DeleteConversation(id string) error {
	return dao.DB.Delete(&model.Conversation{}, id).Error
}
