package usecase

import (
	"errors"
	"time"

	"github.com/collab-platform/backend/internal/domain"
	"github.com/google/uuid"
)

var (
	ErrInvalidOperation = errors.New("invalid operation")
	ErrConflict         = errors.New("operation conflict")
)

type CollaborationRepository interface {
	GetDocumentByID(id uuid.UUID) (*domain.Document, error)
	UpdateDocument(doc *domain.Document) error
	GetPermission(userID, docID uuid.UUID) (*domain.DocumentPermission, error)
	CreateActivity(activity *domain.Activity) error
}

type CollaborationUsecase struct {
	repo CollaborationRepository
}

func NewCollaborationUsecase(repo CollaborationRepository) *CollaborationUsecase {
	return &CollaborationUsecase{repo: repo}
}

// ApplyOperation applies a CRDT-based operation to a document
func (c *CollaborationUsecase) ApplyOperation(userID, docID uuid.UUID, op domain.Operation) (*domain.Document, error) {
	// Check permission
	perm, err := c.repo.GetPermission(userID, docID)
	if err != nil {
		return nil, errors.New("permission denied")
	}

	if !perm.Role.CanEdit() {
		return nil, errors.New("read-only access")
	}

	doc, err := c.repo.GetDocumentByID(docID)
	if err != nil {
		return nil, err
	}

	// Apply CRDT operation
	newContent := c.applyCRDTOperation(doc.Content, op)
	doc.Content = newContent
	doc.Version++
	doc.UpdatedAt = time.Now()

	if err := c.repo.UpdateDocument(doc); err != nil {
		return nil, err
	}

	// Log activity
	activity := &domain.Activity{
		ID:         uuid.New(),
		DocumentID: docID,
		UserID:     userID,
		Action:     "edit",
		Details:    op.Type,
		CreatedAt:  time.Now(),
	}
	c.repo.CreateActivity(activity)

	return doc, nil
}

// applyCRDTOperation applies a CRDT operation to content
// This is a simplified CRDT implementation - in production, you'd use a more sophisticated approach
func (c *CollaborationUsecase) applyCRDTOperation(content string, op domain.Operation) string {
	switch op.Type {
	case "insert":
		if op.Position < 0 {
			op.Position = 0
		}
		if op.Position > len(content) {
			op.Position = len(content)
		}
		return content[:op.Position] + op.Content + content[op.Position:]
	case "delete":
		if op.Position < 0 || op.Position >= len(content) {
			return content
		}
		end := op.Position + op.Length
		if end > len(content) {
			end = len(content)
		}
		return content[:op.Position] + content[end:]
	default:
		return content
	}
}

// TransformOperation transforms an operation against another operation (OT-like)
func (c *CollaborationUsecase) TransformOperation(op1, op2 domain.Operation) domain.Operation {
	// Simplified OT - in production, use proper OT algorithm
	if op1.Type == "insert" && op2.Type == "insert" {
		if op1.Position <= op2.Position {
			// op1 happens before op2, no change needed
			return op1
		} else {
			// op1 happens after op2, adjust position
			op1.Position += len(op2.Content)
			return op1
		}
	}

	if op1.Type == "insert" && op2.Type == "delete" {
		if op1.Position <= op2.Position {
			return op1
		} else if op1.Position > op2.Position+op2.Length {
			op1.Position -= op2.Length
			return op1
		}
		// Conflict case - simplified handling
		return op1
	}

	return op1
}

