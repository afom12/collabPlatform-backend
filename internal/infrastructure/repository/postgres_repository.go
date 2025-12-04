package repository

import (
	"github.com/collab-platform/backend/internal/domain"
	"github.com/collab-platform/backend/internal/usecase"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PostgresAuthRepository struct {
	db *gorm.DB
}

func NewPostgresAuthRepository(db *gorm.DB) usecase.AuthRepository {
	return &PostgresAuthRepository{db: db}
}

func (r *PostgresAuthRepository) CreateUser(user *domain.User) error {
	return r.db.Create(user).Error
}

func (r *PostgresAuthRepository) GetUserByEmail(email string) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *PostgresAuthRepository) GetUserByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	err := r.db.Where("id = ?", id).First(&user).Error
	return &user, err
}

type PostgresDocumentRepository struct {
	db *gorm.DB
}

func NewPostgresDocumentRepository(db *gorm.DB) usecase.DocumentRepository {
	return &PostgresDocumentRepository{db: db}
}

func (r *PostgresDocumentRepository) CreateDocument(doc *domain.Document) error {
	return r.db.Create(doc).Error
}

func (r *PostgresDocumentRepository) GetDocumentByID(id uuid.UUID) (*domain.Document, error) {
	var doc domain.Document
	err := r.db.Where("id = ?", id).First(&doc).Error
	return &doc, err
}

func (r *PostgresDocumentRepository) UpdateDocument(doc *domain.Document) error {
	return r.db.Save(doc).Error
}

func (r *PostgresDocumentRepository) DeleteDocument(id uuid.UUID) error {
	return r.db.Delete(&domain.Document{}, id).Error
}

func (r *PostgresDocumentRepository) GetUserDocuments(userID uuid.UUID) ([]*domain.Document, error) {
	var docs []*domain.Document
	err := r.db.Where("owner_id = ?", userID).
		Or("id IN (SELECT document_id FROM document_permissions WHERE user_id = ?)", userID).
		Find(&docs).Error
	return docs, err
}

func (r *PostgresDocumentRepository) CreatePermission(perm *domain.DocumentPermission) error {
	return r.db.Create(perm).Error
}

func (r *PostgresDocumentRepository) GetPermission(userID, docID uuid.UUID) (*domain.DocumentPermission, error) {
	var perm domain.DocumentPermission
	err := r.db.Where("user_id = ? AND document_id = ?", userID, docID).First(&perm).Error
	return &perm, err
}

func (r *PostgresDocumentRepository) UpdatePermission(perm *domain.DocumentPermission) error {
	return r.db.Save(perm).Error
}

func (r *PostgresDocumentRepository) GetDocumentPermissions(docID uuid.UUID) ([]*domain.DocumentPermission, error) {
	var perms []*domain.DocumentPermission
	err := r.db.Where("document_id = ?", docID).Find(&perms).Error
	return perms, err
}

func (r *PostgresDocumentRepository) CreateVersion(version *domain.DocumentVersion) error {
	return r.db.Create(version).Error
}

func (r *PostgresDocumentRepository) GetDocumentVersions(docID uuid.UUID, limit int) ([]*domain.DocumentVersion, error) {
	var versions []*domain.DocumentVersion
	query := r.db.Where("document_id = ?", docID).Order("version DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&versions).Error
	return versions, err
}

func (r *PostgresDocumentRepository) CreateActivity(activity *domain.Activity) error {
	return r.db.Create(activity).Error
}

func (r *PostgresDocumentRepository) GetDocumentActivities(docID uuid.UUID, limit int) ([]*domain.Activity, error) {
	var activities []*domain.Activity
	query := r.db.Where("document_id = ?", docID).Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&activities).Error
	return activities, err
}

type PostgresCollaborationRepository struct {
	db *gorm.DB
}

func NewPostgresCollaborationRepository(db *gorm.DB) usecase.CollaborationRepository {
	return &PostgresCollaborationRepository{db: db}
}

func (r *PostgresCollaborationRepository) GetDocumentByID(id uuid.UUID) (*domain.Document, error) {
	var doc domain.Document
	err := r.db.Where("id = ?", id).First(&doc).Error
	return &doc, err
}

func (r *PostgresCollaborationRepository) UpdateDocument(doc *domain.Document) error {
	return r.db.Save(doc).Error
}

func (r *PostgresCollaborationRepository) GetPermission(userID, docID uuid.UUID) (*domain.DocumentPermission, error) {
	var perm domain.DocumentPermission
	err := r.db.Where("user_id = ? AND document_id = ?", userID, docID).First(&perm).Error
	return &perm, err
}

func (r *PostgresCollaborationRepository) CreateActivity(activity *domain.Activity) error {
	return r.db.Create(activity).Error
}

