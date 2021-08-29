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
} // @name Base

// Represent business model of User
type User struct {
	Base

	Username  string   `json:"username" xml:"username" gorm:"column:username;index;not null"`
	Email     string   `json:"email" xml:"email" gorm:"column:email;index;not null"`
	UserRole  UserRole `json:"userRole" xml:"userrole" gorm:"column:user_role;not null"`
	AvatarURL string   `json:"avatarUrl" xml:"avatarurl" gorm:"column:avatar_url;"`

	// one-to-many relations
	Post         []Post         `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;ForeignKey:UserID"`
	Comment      []Comment      `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;ForeignKey:UserID"`
	AuthProvider []AuthProvider `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;ForeignKey:UserID"`
} // @name User

type AuthProvider struct {
	Base

	// fk
	UserID uuid.UUID `json:"userId" xml:"userid" gorm:"column:user_id;type:char(36);index;not null"`

	ProviderUID      string           `json:"providerUid" xml:"providerUid" gorm:"column:provider_uid;index;not null"`
	RefreshToken     string           `json:"refreshToken" xml:"refreshtoken" gorm:"column:refresh_token;type:char(255);index;not null"`
	AuthProviderType AuthProviderType `json:"authProviderType" xml:"authproviderType" gorm:"column:auth_provider_type;not null"`
} // @name AuthProvider

// Represent business model of Post
type Post struct {
	Base

	// fk
	UserID uuid.UUID `json:"userId" xml:"userid" gorm:"column:user_id;type:char(36);index;not null"`

	Title string `json:"title" xml:"title" gorm:"column:title;not null"`
	Body  string `json:"body" xml:"body" gorm:"column:body;not null"`

	// one-to-many relation
	Comment []Comment `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;ForeignKey:PostID"`
} // @name Post

// Represent business model for Comment
type Comment struct {
	Base

	// fk
	PostID uuid.UUID `json:"postId" xml:"postId" gorm:"column:post_id;type:char(36);index;not null"`
	UserID uuid.UUID `json:"userId" xml:"userId" gorm:"column:user_id;type:char(36);index;not null"`

	Name string `json:"name" xml:"name" gorm:"column:name;not null"`
	Body string `json:"body" xml:"body" gorm:"column:body;not null"`
} // @name Comment

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

// AuthProviderType represent auth providers.
type AuthProviderType string

// Auth providers.
const (
	AuthProviderTypeGoogle   AuthProviderType = "google"
	AuthProviderTypeFacebook AuthProviderType = "facebook"
	AuthProviderTypeGithub   AuthProviderType = "github"
)

// origin
// https://console.cloud.google.com/apis/credentials?folder=&hl=RU&organizationId=&project=gorest-323217
// decliner
// https://myaccount.google.com/permissions?continue=https%3A%2F%2Fmyaccount.google.com%2Fsecurity

// origin
// https://developers.facebook.com/apps/208618197804435/settings/basic/
// decliner
// https://www.facebook.com/settings?tab=applications&ref=settings

// origin
// https://github.com/settings/developers
// declinder
// https://github.com/settings/applications

// TODO: remove
// err := echo.FormFieldBinder(c).
// 	MustString("postId", &a).
// 	MustString("name", &b).
// 	MustString("body", &z).
// 	BindError()

// TODO: remove auth routes mocks
// https://stackoverflow.com/questions/43240970/how-to-mock-http-client-do-method
// https://dou.ua/lenta/articles/golang-httptest/
// https://dev.to/clavinjune/mocking-http-call-in-golang-15i5
// https://olegcodes.medium.com/mocking-http-services-in-go-6b76215a81c9
