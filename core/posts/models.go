package posts

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Represent business model of Post
type Post struct {
	ID     uuid.UUID `json:"id" xml:"id" gorm:"column:id;type:char(36);primary_key;"`
	UserID uuid.UUID `json:"userId" xml:"userid" gorm:"column:user_id;type:char(36);index;not null"`
	Title  string    `json:"title" xml:"title" gorm:"column:title;not null"`
	Body   string    `json:"body" xml:"body" gorm:"column:body;not null"`

	CreatedAt time.Time `json:"-" xml:"-" gorm:"column:created_at;index"`
	UpdatedAt time.Time `json:"-" xml:"-" gorm:"column:updated_at"`
}

func (p *Post) BeforeCreate(tx *gorm.DB) (err error) {
	p.ID = uuid.New()
	return
}

// MimeType represent mime types of sort.
type MimeType string

// Mime types.
const (
	MimeTypesXML  MimeType = "application/xml"
	MimeTypesJSON MimeType = "application/json"
)
