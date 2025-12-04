package domain

import (
	"time"

	"github.com/google/uuid"
)

type Role string

const (
	RoleOwner  Role = "owner"
	RoleEditor Role = "editor"
	RoleViewer Role = "viewer"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	Email     string    `json:"email" gorm:"uniqueIndex;not null"`
	Username  string    `json:"username" gorm:"uniqueIndex;not null"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DocumentPermission struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	DocumentID uuid.UUID `json:"document_id" gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null;index"`
	Role       Role      `json:"role" gorm:"type:varchar(20);not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (r Role) CanEdit() bool {
	return r == RoleOwner || r == RoleEditor
}

func (r Role) CanView() bool {
	return r == RoleOwner || r == RoleEditor || r == RoleViewer
}

