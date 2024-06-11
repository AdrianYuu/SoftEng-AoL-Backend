package model

type User struct {
	ID             string  `json:"id" gorm:"primaryKey"`
	Email          string  `json:"email" gorm:"uniqueIndex;not null"`
	Password       string  `json:"password" gorm:"not null"`
	Username       string  `json:"username"`
	DisplayName    string  `json:"displayName"`
	ProfilePicture *string `gorm:"foreignKey:OwnerID"`

	// associations
	Conversations []*Conversation `json:"conversations" gorm:"many2many:user_conversations;"`
}

type CreateUserInput struct {
	Email       string `json:"email"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
}
