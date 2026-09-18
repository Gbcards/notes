package httpadapter

import (
	"context"
	"errors"
	"net/http"

	domain "notes-service/internal/domain/note"

	"github.com/gin-gonic/gin"
)

// NoteCreator, NoteLister, NoteUpdater, and NoteDeleter are the inbound ports
// this adapter depends on. They are exported so the DI wiring (internal/platform)
// can bind a use case's concrete type to the interface via fx.As.
type NoteCreator interface {
	Execute(ctx context.Context, title, content string) (*domain.Note, error)
}

type NoteLister interface {
	Execute(ctx context.Context) ([]*domain.Note, error)
}

type NoteUpdater interface {
	Execute(ctx context.Context, id, title, content string) (*domain.Note, error)
}

type NoteDeleter interface {
	Execute(ctx context.Context, id string) error
}

type NoteHandler struct {
	create NoteCreator
	list   NoteLister
	update NoteUpdater
	delete NoteDeleter
}

func NewNoteHandler(create NoteCreator, list NoteLister, update NoteUpdater, delete NoteDeleter) *NoteHandler {
	return &NoteHandler{create: create, list: list, update: update, delete: delete}
}

func (h *NoteHandler) Create(c *gin.Context) {
	var req createNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	n, err := h.create.Execute(c.Request.Context(), req.Title, req.Content)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toNoteResponse(n))
}

func (h *NoteHandler) List(c *gin.Context) {
	notes, err := h.list.Execute(c.Request.Context())
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toNoteResponses(notes))
}

func (h *NoteHandler) Update(c *gin.Context) {
	var req updateNoteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	n, err := h.update.Execute(c.Request.Context(), c.Param("id"), req.Title, req.Content)
	if err != nil {
		writeError(c, err)
		return
	}
	c.JSON(http.StatusOK, toNoteResponse(n))
}

func (h *NoteHandler) Delete(c *gin.Context) {
	if err := h.delete.Execute(c.Request.Context(), c.Param("id")); err != nil {
		writeError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func writeError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, domain.ErrEmptyContent):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, domain.ErrNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
