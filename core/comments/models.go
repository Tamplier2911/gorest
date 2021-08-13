package comments

import (
	"time"

	"github.com/Tamplier2911/gorest/pkg/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represent business model of User
type User struct {
	BaID uuid.UUID `json:"id" xml:"id" gorm:"column:id;type:char(36);primary_key;"`

	Name        string   `json:"name" xml:"name" gorm:"column:name;not null"`
	Username    string   `json:"username" xml:"username" gorm:"column:username;index;not null"`
	Email       string   `json:"email" xml:"email" gorm:"column:email;index;not null"`
	UserRole    UserRole `json:"userRole" xml:"userrole" gorm:"column:user_role;not null"`
	PhoneNumber string   `json:"phoneNumber" xml:"phonenumber" gorm:"column:phone_number;index"`
	AvatarURL   string   `json:"avatarUrl" xml:"avatarurl" gorm:"column:avatar_url;"`

	CreatedAt time.Time      `json:"-" xml:"-" gorm:"column:created_at;index"`
	UpdatedAt time.Time      `json:"-" xml:"-" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" xml:"-" gorm:"column:deleted_at;index"`
}

// UserRole represent user roles.
type UserRole string

// Mime types.
const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleModerator UserRole = "moderator"
	UserRoleUser      UserRole = "user"
)

// Represent business model for Comment
type Comment struct {
	ID uuid.UUID `json:"id" xml:"id" gorm:"column:id;type:char(36);primary_key;"`

	PostID uuid.UUID `json:"postId" xml:"postId" gorm:"column:post_id;type:char(36);index;not null"`
	UserID uuid.UUID `json:"userId" xml:"userId" gorm:"column:user_id;type:char(36);index;not null"`

	Name string `json:"name" xml:"name" gorm:"column:name;not null"`
	Body string `json:"body" xml:"body" gorm:"column:body;not null"`

	CreatedAt time.Time      `json:"-" xml:"-" gorm:"column:created_at;index"`
	UpdatedAt time.Time      `json:"-" xml:"-" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" xml:"-" gorm:"column:deleted_at;index"`

	Post models.Post `json:"-" xml:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PostID"`
	User User        `json:"-" xml:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
}

func (c *Comment) BeforeCreate(tx *gorm.DB) (err error) {
	c.ID = uuid.New()
	return
}

// MimeType represent mime types of sort.
type MimeType string

// Mime types.
const (
	MimeTypesXML  MimeType = "application/xml"
	MimeTypesJSON MimeType = "application/json"
)
