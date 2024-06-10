package controllers

import (
	"errors"
	"github.com/badaccuracyid/softeng_backend/src/database/dao"
	"github.com/badaccuracyid/softeng_backend/src/model"
	"github.com/badaccuracyid/softeng_backend/src/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sync"
)

type ChatController interface {
	SetContext(ctx *gin.Context)

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

func NewChatController(db *gorm.DB) ChatController {
	return &chatController{
		chatDAO: dao.NewChatDAO(db),
	}
}

func (s *chatController) SetContext(ctx *gin.Context) {
	s.ctx = ctx
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
	message := &model.Message{
		ID:             uuid.New().String(),
		ConversationID: input.ConversationID,
		SenderID:       input.SenderID,
		ContentType:    input.ContentType,
		Content:        input.Content,
	}

	if err := s.chatDAO.CreateMessage(message); err != nil {
		return nil, err
	}

	triggerSubscription(input.ConversationID, message)
	return message, nil
}

func (s *chatController) AddUserToConversation(conversationID string, userID string) (*model.Conversation, error) {
	userService := NewUserService(s.chatDAO.DB)
	userService.SetContext(s.ctx)

	user, err := userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	conversation, err := s.GetConversation(conversationID)
	if err != nil {
		return nil, err
	}

	if membersContainsUser(conversation.Members, user) {
		return conversation, nil
	}

	conversation.Members = append(conversation.Members, user)

	if err := s.chatDAO.DB.Association("Members").Append(user); err != nil {
		return nil, err
	}

	return conversation, nil
}

func (s *chatController) RemoveUserFromConversation(conversationID string, userID string) (*model.Conversation, error) {
	userService := NewUserService(s.chatDAO.DB)
	userService.SetContext(s.ctx)

	user, err := userService.GetUserByID(userID)
	if err != nil {
		return nil, err
	}

	conversation, err := s.GetConversation(conversationID)
	if err != nil {
		return nil, err
	}

	if !membersContainsUser(conversation.Members, user) {
		return conversation, nil
	}

	if err := s.chatDAO.DB.Association("Members").Delete(user); err != nil {
		return nil, err
	}

	return conversation, nil
}

var (
	subscriptions      = make(map[string][]*model.MessageSubscription)
	subscriptionsMutex sync.Mutex
)

func (s *chatController) NewMessageSubscription(conversationID string) (<-chan *model.Message, chan<- struct{}, error) {
	subscription := &model.MessageSubscription{
		MessageChannel: make(chan *model.Message),
		DoneChannel:    make(chan struct{}),
	}

	onSubscribe(conversationID, subscription)
	return subscription.MessageChannel, subscription.DoneChannel, nil
}

func onSubscribe(conversationId string, subscription *model.MessageSubscription) {
	subscriptionsMutex.Lock()
	defer subscriptionsMutex.Unlock()
	subscriptions[conversationId] = append(subscriptions[conversationId], subscription)
}

func triggerSubscription(conversationId string, message *model.Message) {
	subscriptionsMutex.Lock()
	defer subscriptionsMutex.Unlock()

	subscribers, found := subscriptions[conversationId]
	if found {
		for _, subscriber := range subscribers {
			select {
			case <-subscriber.DoneChannel:
				subscriber = nil
			case subscriber.MessageChannel <- message:
				// Message went through, do nothing
			}
		}
	}
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
