package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represent base model
type Base struct {
	ID uuid.UUID `json:"id" xml:"id" gorm:"column:id;type:char(36);primary_key;"`

	CreatedAt time.Time      `json:"-" xml:"-" gorm:"column:created_at;index"`
	UpdatedAt time.Time      `json:"-" xml:"-" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" xml:"-"  gorm:"column:deleted_at;index"`
}

// Represent business model of User
type User struct {
	Base

	Name        string   `json:"name" xml:"name" gorm:"column:name;not null"`
	Username    string   `json:"username" xml:"username" gorm:"column:username;index;not null"`
	Email       string   `json:"email" xml:"email" gorm:"column:email;index;not null"`
	UserRole    UserRole `json:"userRole" xml:"userrole" gorm:"column:user_role;not null"`
	PhoneNumber string   `json:"phoneNumber" xml:"phonenumber" gorm:"column:phone_number;index"`
	AvatarURL   string   `json:"avatarUrl" xml:"avatarurl" gorm:"column:avatar_url;"`
}

// Represent business model of Post
type Post struct {
	Base

	UserID uuid.UUID `json:"userId" xml:"userid" gorm:"column:user_id;type:char(36);index;not null"`
	Title  string    `json:"title" xml:"title" gorm:"column:title;not null"`
	Body   string    `json:"body" xml:"body" gorm:"column:body;not null"`

	User User `json:"-" xml:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
}

// Represent business model for Comment
type Comment struct {
	Base

	PostID uuid.UUID `json:"postId" xml:"postId" gorm:"column:post_id;type:char(36);index;not null"`
	UserID uuid.UUID `json:"userId" xml:"userId" gorm:"column:user_id;type:char(36);index;not null"`

	Name string `json:"name" xml:"name" gorm:"column:name;not null"`
	Body string `json:"body" xml:"body" gorm:"column:body;not null"`

	Post Post `json:"-" xml:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PostID"`
	User User `json:"-" xml:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:UserID"`
}

func (b *Base) BeforeCreate(tx *gorm.DB) (err error) {
	b.ID = uuid.New()
	return
}

// MimeType represent mime types of sort.
type MimeType string

// Mime types.
const (
	MimeTypesXML  MimeType = "application/xml"
	MimeTypesJSON MimeType = "application/json"
)

// UserRole represent user roles.
type UserRole string

// User roles.
const (
	UserRoleAdmin     UserRole = "admin"
	UserRoleModerator UserRole = "moderator"
	UserRoleUser      UserRole = "user"
)
