package model

type User struct {
	ID             uint    `json:"id" gorm:"primaryKey"`
	Email          string  `json:"email" gorm:"uniqueIndex;not null"`
	Username       string  `json:"username"`
	DisplayName    string  `json:"displayName"`
	ProfilePicture *string `gorm:"foreignKey:OwnerID"`

	// associations
	Conversations []*Conversation `json:"conversations" gorm:"many2many:user_conversations;"`
}
