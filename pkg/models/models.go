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

// Represent business model of Post
type Post struct {
	Base

	UserID uuid.UUID `json:"userId" xml:"userid" gorm:"column:user_id;type:char(36);index;not null"`
	Title  string    `json:"title" xml:"title" gorm:"column:title;not null"`
	Body   string    `json:"body" xml:"body" gorm:"column:body;not null"`
}

// Represent business model for Comment
type Comment struct {
	Base

	PostID uuid.UUID `json:"postId" xml:"postId" gorm:"column:post_id;type:char(36);index;not null"`
	Email  string    `json:"email" xml:"email" gorm:"column:email;index;not null"`
	Name   string    `json:"name" xml:"name" gorm:"column:name;not null"`
	Body   string    `json:"body" xml:"body" gorm:"column:body;not null"`

	Post Post `json:"-" xml:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PostID"`
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
