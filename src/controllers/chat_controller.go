package controllers

import (
	"errors"
	"github.com/badaccuracyid/softeng_backend/src/database/dao"
	"github.com/badaccuracyid/softeng_backend/src/model"
	"github.com/badaccuracyid/softeng_backend/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ChatController interface {
	CreateConversation(input model.CreateConversationInput) (*model.Conversation, error)
	GetConversation(id string) (*model.Conversation, error)
	DeleteConversation(id string) error

	SendMessage(input model.SendMessageInput) (*model.Message, error)

	AddUserToConversation(conversationID string, userID string) (*model.Conversation, error)
	RemoveUserFromConversation(conversationID string, userID string) (*model.Conversation, error)

	NewMessageSubscription(conversationID string) (<-chan *model.Message, chan<- struct{}, error)
}

type chatController struct {
	ctx     *gin.Context
	chatDAO *dao.ChatDAO
}

func NewChatController(chatDAO *dao.ChatDAO) ChatController {
	return &chatController{
		chatDAO: chatDAO,
	}
}

func (s *chatController) CreateConversation(input model.CreateConversationInput) (*model.Conversation, error) {
	userId := utils.GetCurrentUserID(s.ctx)
	if userId == "" {
		return nil, errors.New("user not found")
	}

	userController := NewUserService(s.chatDAO.DB)
	userController.SetContext(s.ctx)

	user, err := userController.GetUserByID(userId)
	if err != nil {
		return nil, err
	}

	participantUser, err := userController.GetUsersByID(input.MemberIds)
	if err != nil {
		return nil, err
	}

	// make sure that the user is in the conversation
	if !membersContainsUser(participantUser, user) {
		participantUser = append(participantUser, user)
	}

	conversation := &model.Conversation{
		ID:      uuid.New().String(),
		Title:   input.Title,
		Members: participantUser,
	}

	if err := s.chatDAO.CreateConversation(conversation); err != nil {
		return nil, err
	}

	return conversation, nil
}

func (s *chatController) GetConversation(id string) (*model.Conversation, error) {
	return s.chatDAO.GetConversationByID(id)
}

func (s *chatController) DeleteConversation(id string) error {
	return s.chatDAO.DeleteConversation(id)
}

func (s *chatController) SendMessage(input model.SendMessageInput) (*model.Message, error) {
	return nil, nil
}

func (s *chatController) AddUserToConversation(conversationID string, userID string) (*model.Conversation, error) {
	return nil, nil
}

func (s *chatController) RemoveUserFromConversation(conversationID string, userID string) (*model.Conversation, error) {
	return nil, nil
}

func (s *chatController) NewMessageSubscription(conversationID string) (<-chan *model.Message, chan<- struct{}, error) {
	return nil, nil, nil
}

func membersContainsUser(members []*model.User, user *model.User) bool {
	if members == nil {
		return false
	}

	for _, member := range members {
		if member.ID == user.ID {
			return true
		}
	}

	return false
}
