package model

type Conversation struct {
	ID       string     `json:"id" gorm:"primaryKey"`
	Title    string     `json:"title" gorm:"not null"`
	Members  []*User    `json:"members" gorm:"many2many:user_conversations;"`
	Messages []*Message `json:"messages" gorm:"foreignKey:ConversationID"`
}

type Message struct {
	ID             string             `json:"id" gorm:"primaryKey"`
	SenderID       string             `json:"sender_id" gorm:"not null"`
	Sender         User               `json:"sender" gorm:"not null;foreignKey:SenderID"`
	ConversationID string             `json:"conversation_id" gorm:"not null"`
	Conversation   Conversation       `json:"conversation" gorm:"foreignKey:ConversationID"`
	Content        string             `json:"content"`
	ContentType    MessageContentType `json:"contentType" gorm:"not null"`
}

type MessageContentType string

const (
	MessageContentTypeText  MessageContentType = "TEXT"
	MessageContentTypeImage MessageContentType = "IMAGE"
)

var AllMessageContentType = []MessageContentType{
	MessageContentTypeText,
	MessageContentTypeImage,
}

func (e MessageContentType) IsValid() bool {
	switch e {
	case MessageContentTypeText, MessageContentTypeImage:
		return true
	}
	return false
}

func (e MessageContentType) String() string {
	return string(e)
}
