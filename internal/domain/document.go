package domain

import (
	"time"

	"github.com/google/uuid"
)

type DocumentType string

const (
	DocumentTypeText     DocumentType = "text"
	DocumentTypeNote     DocumentType = "note"
	DocumentTypeWhiteboard DocumentType = "whiteboard"
	DocumentTypeTask     DocumentType = "task"
)

type Document struct {
	ID          uuid.UUID   `json:"id" gorm:"type:uuid;primary_key"`
	Title       string      `json:"title" gorm:"not null"`
	Content     string      `json:"content" gorm:"type:text"`
	Type        DocumentType `json:"type" gorm:"type:varchar(20);not null"`
	OwnerID     uuid.UUID   `json:"owner_id" gorm:"type:uuid;not null;index"`
	IsPublic    bool        `json:"is_public" gorm:"default:false"`
	ShareToken  string      `json:"share_token" gorm:"uniqueIndex"`
	Version     int64       `json:"version" gorm:"default:0"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

type DocumentVersion struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	DocumentID uuid.UUID `json:"document_id" gorm:"type:uuid;not null;index"`
	Version    int64     `json:"version" gorm:"not null"`
	Content    string    `json:"content" gorm:"type:text"`
	CreatedBy  uuid.UUID `json:"created_by" gorm:"type:uuid;not null"`
	CreatedAt  time.Time `json:"created_at"`
}

type Activity struct {
	ID         uuid.UUID `json:"id" gorm:"type:uuid;primary_key"`
	DocumentID uuid.UUID `json:"document_id" gorm:"type:uuid;not null;index"`
	UserID     uuid.UUID `json:"user_id" gorm:"type:uuid;not null"`
	Action     string    `json:"action" gorm:"type:varchar(50);not null"`
	Details    string    `json:"details" gorm:"type:text"`
	CreatedAt  time.Time `json:"created_at"`
}

