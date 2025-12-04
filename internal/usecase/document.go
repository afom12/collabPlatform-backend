package usecase

import (
	"errors"
	"time"

	"github.com/collab-platform/backend/internal/domain"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var (
	ErrDocumentNotFound   = errors.New("document not found")
	ErrPermissionDenied   = errors.New("permission denied")
	ErrInvalidDocumentType = errors.New("invalid document type")
)

type DocumentRepository interface {
	CreateDocument(doc *domain.Document) error
	GetDocumentByID(id uuid.UUID) (*domain.Document, error)
	UpdateDocument(doc *domain.Document) error
	DeleteDocument(id uuid.UUID) error
	GetUserDocuments(userID uuid.UUID) ([]*domain.Document, error)
	CreatePermission(perm *domain.DocumentPermission) error
	GetPermission(userID, docID uuid.UUID) (*domain.DocumentPermission, error)
	UpdatePermission(perm *domain.DocumentPermission) error
	GetDocumentPermissions(docID uuid.UUID) ([]*domain.DocumentPermission, error)
	CreateVersion(version *domain.DocumentVersion) error
	GetDocumentVersions(docID uuid.UUID, limit int) ([]*domain.DocumentVersion, error)
	CreateActivity(activity *domain.Activity) error
	GetDocumentActivities(docID uuid.UUID, limit int) ([]*domain.Activity, error)
}

type DocumentUsecase struct {
	repo DocumentRepository
}

func NewDocumentUsecase(repo DocumentRepository) *DocumentUsecase {
	return &DocumentUsecase{repo: repo}
}

func (d *DocumentUsecase) CreateDocument(userID uuid.UUID, title string, docType domain.DocumentType) (*domain.Document, error) {
	if docType != domain.DocumentTypeText && docType != domain.DocumentTypeNote &&
		docType != domain.DocumentTypeWhiteboard && docType != domain.DocumentTypeTask {
		return nil, ErrInvalidDocumentType
	}

	doc := &domain.Document{
		ID:         uuid.New(),
		Title:      title,
		Content:    "",
		Type:       docType,
		OwnerID:    userID,
		IsPublic:   false,
		ShareToken: uuid.New().String(),
		Version:    0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if err := d.repo.CreateDocument(doc); err != nil {
		return nil, err
	}

	// Create owner permission
	perm := &domain.DocumentPermission{
		ID:         uuid.New(),
		DocumentID: doc.ID,
		UserID:     userID,
		Role:       domain.RoleOwner,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := d.repo.CreatePermission(perm); err != nil {
		return nil, err
	}

	// Log activity
	activity := &domain.Activity{
		ID:         uuid.New(),
		DocumentID: doc.ID,
		UserID:     userID,
		Action:     "created",
		Details:    "Document created",
		CreatedAt:  time.Now(),
	}
	d.repo.CreateActivity(activity)

	return doc, nil
}

func (d *DocumentUsecase) GetDocument(userID, docID uuid.UUID) (*domain.Document, error) {
	doc, err := d.repo.GetDocumentByID(docID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}

	// Check permission
	perm, err := d.repo.GetPermission(userID, docID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if perm == nil && !doc.IsPublic {
		return nil, ErrPermissionDenied
	}

	return doc, nil
}

func (d *DocumentUsecase) UpdateDocument(userID, docID uuid.UUID, title, content string) (*domain.Document, error) {
	doc, err := d.repo.GetDocumentByID(docID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDocumentNotFound
		}
		return nil, err
	}

	// Check permission
	perm, err := d.repo.GetPermission(userID, docID)
	if err != nil {
		return nil, ErrPermissionDenied
	}

	if !perm.Role.CanEdit() {
		return nil, ErrPermissionDenied
	}

	doc.Title = title
	doc.Content = content
	doc.Version++
	doc.UpdatedAt = time.Now()

	if err := d.repo.UpdateDocument(doc); err != nil {
		return nil, err
	}

	// Create version snapshot
	version := &domain.DocumentVersion{
		ID:         uuid.New(),
		DocumentID: doc.ID,
		Version:    doc.Version,
		Content:    content,
		CreatedBy:  userID,
		CreatedAt:  time.Now(),
	}
	d.repo.CreateVersion(version)

	// Log activity
	activity := &domain.Activity{
		ID:         uuid.New(),
		DocumentID: doc.ID,
		UserID:     userID,
		Action:     "updated",
		Details:    "Document updated",
		CreatedAt:  time.Now(),
	}
	d.repo.CreateActivity(activity)

	return doc, nil
}

func (d *DocumentUsecase) ShareDocument(ownerID, docID uuid.UUID, userID uuid.UUID, role domain.Role) error {
	doc, err := d.repo.GetDocumentByID(docID)
	if err != nil {
		return ErrDocumentNotFound
	}

	if doc.OwnerID != ownerID {
		return ErrPermissionDenied
	}

	perm := &domain.DocumentPermission{
		ID:         uuid.New(),
		DocumentID: docID,
		UserID:     userID,
		Role:       role,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return d.repo.CreatePermission(perm)
}

func (d *DocumentUsecase) GetUserDocuments(userID uuid.UUID) ([]*domain.Document, error) {
	return d.repo.GetUserDocuments(userID)
}

func (d *DocumentUsecase) GetDocumentVersions(docID uuid.UUID, limit int) ([]*domain.DocumentVersion, error) {
	return d.repo.GetDocumentVersions(docID, limit)
}

func (d *DocumentUsecase) GetDocumentActivities(docID uuid.UUID, limit int) ([]*domain.Activity, error) {
	return d.repo.GetDocumentActivities(docID, limit)
}

