package comments

import (
	"time"

	"github.com/Tamplier2911/gorest/core/posts"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represent business model for Comment
type Comment struct {
	ID     uuid.UUID `json:"id" xml:"id" gorm:"column:id;type:char(36);primary_key;"`
	PostID uuid.UUID `json:"postId" xml:"postId" gorm:"column:post_id;type:char(36);index;not null"`

	Email string `json:"email" xml:"email" gorm:"column:email;index;not null"`
	Name  string `json:"name" xml:"name" gorm:"column:name;not null"`
	Body  string `json:"body" xml:"body" gorm:"column:body;not null"`

	CreatedAt time.Time      `json:"-" xml:"-" gorm:"column:created_at;index"`
	UpdatedAt time.Time      `json:"-" xml:"-" gorm:"column:updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" xml:"-" gorm:"column:deleted_at;index"`

	Post posts.Post `json:"-" xml:"-" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:PostID"`
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
