package dao

import (
	"fmt"

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
	err := dao.DB.Preload("Members").Preload("Messages.Sender").First(&conversation, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return conversation, nil
}

func (dao *ChatDAO) GetConversationsForUser(userID string) ([]*model.Conversation, error) {
	conversations := []*model.Conversation{}
	err := dao.DB.
		Joins("JOIN user_conversations ON user_conversations.conversation_id = conversations.id").
		Joins("JOIN users ON users.id = user_conversations.user_id").
		Preload("Members").
		Preload("Messages.Sender").
		Where("users.id = ?", userID).
		Find(&conversations).Error
	if err != nil {
		return nil, err
	}

	fmt.Println(conversations, userID)

	return conversations, nil
}

func (dao *ChatDAO) UpdateConversation(conversation *model.Conversation) error {
	return dao.DB.Save(conversation).Error
}

func (dao *ChatDAO) DeleteConversation(id string) error {
	return dao.DB.Delete(&model.Conversation{}, "id = ?", id).Error
}

func (dao *ChatDAO) CreateMessage(message *model.Message) error {
	// Create the message
	if err := dao.DB.Create(message).Error; err != nil {
		return err
	}

	// Preload the Sender (User)
	if err := dao.DB.Preload("Sender").First(message, "id = ?", message.ID).Error; err != nil {
		return err
	}

	return nil
}
