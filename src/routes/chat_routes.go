package routes

import (
	"github.com/badaccuracyid/softeng_backend/src/controllers"
	"github.com/badaccuracyid/softeng_backend/src/database"
	"github.com/badaccuracyid/softeng_backend/src/model"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type ChatRoutes struct {
	baseRouter     *gin.RouterGroup
	chatController controllers.ChatController
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewChatRoutes(router *gin.Engine) (*ChatRoutes, error) {
	postgresDatabase, err := database.GetPostgresDatabase()
	if err != nil {
		panic(err)
	}

	chatService := controllers.NewChatController(postgresDatabase)
	baseRouter := router.Group("/api/v1/chats")

	return &ChatRoutes{
		baseRouter:     baseRouter,
		chatController: chatService,
	}, nil
}

func InitializeChatRoutes(router *gin.Engine) *gin.Engine {
	chatRoutes, err := NewChatRoutes(router)
	if err != nil {
		panic(err)
	}

	chatRoutes.registerRoutes()
	return router
}

func (c *ChatRoutes) registerRoutes() {
	c.baseRouter.POST("/create", c.createConversation)
	c.baseRouter.POST("/message", c.sendMessage)
	c.baseRouter.GET("/get/:id", c.getConversation)
	c.baseRouter.GET("/ws/:id", c.handleWebSocket)
}

// createConversation handles the POST /api/v1/chats/create request
// @Summary Create a new chat
// @Description Create a new chat with the input payload
// @Tags chats
// @Accept  json
// @Produce  json
// @Param chat body model.CreateConversationInput true "Chat"
// @Success 201 {object} model.Conversation
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /chats/create [post]
func (c *ChatRoutes) createConversation(ctx *gin.Context) {
	c.chatController.SetContext(ctx)
	payload := model.CreateConversationInput{}
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	conversation, err := c.chatController.CreateConversation(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, conversation)
}

// sendMessage handles the POST /api/v1/chats/message request
// @Summary Send a message
// @Description Send a message with the input payload
// @Tags chats
// @Accept  json
// @Produce  json
// @Param message body model.SendMessageInput true "Message"
// @Success 201 {object} model.Message
// @Failure 400 {string} string
// @Failure 500 {string} string
// @Router /chats/message [post]
func (c *ChatRoutes) sendMessage(ctx *gin.Context) {
	c.chatController.SetContext(ctx)
	payload := model.SendMessageInput{}
	err := ctx.BindJSON(&payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	message, err := c.chatController.SendMessage(payload)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	ctx.JSON(http.StatusCreated, message)
}

// getConversation handles the GET /api/v1/chats/get/:id request
// @Summary Get a chat by ID
// @Description Get a chat by ID
// @Tags chats
// @Accept  json
// @Produce  json
// @Param id path string true "Chat ID"
// @Success 200 {object} model.Message
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /chats/{id} [get]
func (c *ChatRoutes) getConversation(ctx *gin.Context) {
	c.chatController.SetContext(ctx)
	id := ctx.Param("id")
	conversation, err := c.chatController.GetConversation(id)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	if conversation == nil {
		ctx.JSON(http.StatusNotFound, "Conversation not found")
		return
	}

	ctx.JSON(http.StatusOK, conversation)
}

// handleWebSocket handles the GET /api/v1/chats/ws/:id request
// @Summary Handle a websocket connection
// @Description Handle a websocket connection
// @Tags chats
// @Accept  json
// @Produce  json
// @Param id path string true "Chat ID"
// @Success 200 {string} string
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 500 {string} string
// @Router /chats/ws/{id} [get]
func (c *ChatRoutes) handleWebSocket(ctx *gin.Context) {
	conversationID := ctx.Param("id")
	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, "Failed to set websocket upgrade")
		return
	}

	messageChannel, doneChannel, err := c.chatController.NewMessageSubscription(conversationID)
	if err != nil {
		err := conn.Close()
		if err != nil {
			return
		}
		ctx.JSON(http.StatusInternalServerError, "Failed to subscribe to messages")
		return
	}

	defer close(doneChannel)

	go func() {
		for {
			message, ok := <-messageChannel
			if !ok {
				err := conn.Close()
				if err != nil {
					return
				}
				return
			}
			err := conn.WriteJSON(message)
			if err != nil {
				return
			}
		}
	}()

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			doneChannel <- struct{}{}
			conn.Close()
			return
		}
	}

}
