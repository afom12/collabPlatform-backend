package handlers

import (
	"net/http"

	"github.com/collab-platform/backend/internal/domain"
	"github.com/collab-platform/backend/internal/usecase"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DocumentHandler struct {
	docUsecase *usecase.DocumentUsecase
}

func NewDocumentHandler(docUsecase *usecase.DocumentUsecase) *DocumentHandler {
	return &DocumentHandler{docUsecase: docUsecase}
}

type CreateDocumentRequest struct {
	Title string              `json:"title" binding:"required" example:"My First Document"`
	Type  domain.DocumentType `json:"type" binding:"required" example:"text" enums:"text,note,whiteboard,task"`
}

type UpdateDocumentRequest struct {
	Title   string `json:"title" example:"Updated Document Title"`
	Content string `json:"content" example:"Document content here..."`
}

type ShareDocumentRequest struct {
	UserID string      `json:"user_id" binding:"required" example:"550e8400-e29b-41d4-a716-446655440000"`
	Role   domain.Role `json:"role" binding:"required" example:"editor" enums:"owner,editor,viewer"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"error message"`
}

type SuccessMessageResponse struct {
	Message string `json:"message" example:"Document shared successfully"`
}

// CreateDocument godoc
// @Summary      Create a new document
// @Description  Create a new document (text, note, whiteboard, or task)
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        request  body      CreateDocumentRequest  true  "Document creation details"
// @Success      201      {object}  domain.Document
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /documents [post]
func (h *DocumentHandler) CreateDocument(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	var req CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := h.docUsecase.CreateDocument(userID, req.Title, req.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, doc)
}

// GetDocument godoc
// @Summary      Get document by ID
// @Description  Retrieve a specific document by its ID
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Document ID"
// @Success      200  {object}  domain.Document
// @Failure      400   {object}  ErrorResponse
// @Failure      401   {object}  ErrorResponse
// @Failure      403   {object}  ErrorResponse
// @Failure      404   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /documents/{id} [get]
func (h *DocumentHandler) GetDocument(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	doc, err := h.docUsecase.GetDocument(userID, docID)
	if err != nil {
		if err == usecase.ErrDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == usecase.ErrPermissionDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// UpdateDocument godoc
// @Summary      Update document
// @Description  Update document title and/or content
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                 true  "Document ID"
// @Param        request  body      UpdateDocumentRequest  true  "Document update details"
// @Success      200      {object}  domain.Document
// @Failure      400       {object}  ErrorResponse
// @Failure      401       {object}  ErrorResponse
// @Failure      403       {object}  ErrorResponse
// @Failure      404       {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Router       /documents/{id} [put]
func (h *DocumentHandler) UpdateDocument(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	doc, err := h.docUsecase.UpdateDocument(userID, docID, req.Title, req.Content)
	if err != nil {
		if err == usecase.ErrDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == usecase.ErrPermissionDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, doc)
}

// ListDocuments godoc
// @Summary      List user documents
// @Description  Get all documents accessible by the authenticated user
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   domain.Document
// @Failure      401   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /documents [get]
func (h *DocumentHandler) ListDocuments(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	docs, err := h.docUsecase.GetUserDocuments(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, docs)
}

// ShareDocument godoc
// @Summary      Share document with user
// @Description  Share a document with another user with specified role (owner, editor, viewer)
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id       path      string                true  "Document ID"
// @Param        request  body      ShareDocumentRequest  true  "Share details"
// @Success      200      {object}  SuccessMessageResponse
// @Failure      400       {object}  ErrorResponse
// @Failure      401       {object}  ErrorResponse
// @Failure      403       {object}  ErrorResponse
// @Failure      404       {object}  ErrorResponse
// @Failure      500       {object}  ErrorResponse
// @Router       /documents/{id}/share [post]
func (h *DocumentHandler) ShareDocument(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	ownerID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	var req ShareDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	shareUserID, err := uuid.Parse(req.UserID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid share user ID"})
		return
	}

	if err := h.docUsecase.ShareDocument(ownerID, docID, shareUserID, req.Role); err != nil {
		if err == usecase.ErrDocumentNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		if err == usecase.ErrPermissionDenied {
			c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document shared successfully"})
}

// GetVersions godoc
// @Summary      Get document versions
// @Description  Get version history for a document
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Document ID"
// @Success      200  {array}   domain.DocumentVersion
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /documents/{id}/versions [get]
func (h *DocumentHandler) GetVersions(c *gin.Context) {
	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	versions, err := h.docUsecase.GetDocumentVersions(docID, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, versions)
}

// GetActivities godoc
// @Summary      Get document activities
// @Description  Get activity feed for a document (who edited what and when)
// @Tags         documents
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      string  true  "Document ID"
// @Success      200  {array}   domain.Activity
// @Failure      400   {object}  ErrorResponse
// @Failure      500   {object}  ErrorResponse
// @Router       /documents/{id}/activities [get]
func (h *DocumentHandler) GetActivities(c *gin.Context) {
	docID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid document ID"})
		return
	}

	activities, err := h.docUsecase.GetDocumentActivities(docID, 100)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, activities)
}
