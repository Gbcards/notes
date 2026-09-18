package httpadapter

import (
	"time"

	domain "notes-service/internal/domain/note"
)

type createNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type updateNoteRequest struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type noteResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func toNoteResponse(n *domain.Note) noteResponse {
	return noteResponse{
		ID:        n.ID,
		Title:     n.Title,
		Content:   n.Content,
		CreatedAt: n.CreatedAt,
		UpdatedAt: n.UpdatedAt,
	}
}

func toNoteResponses(notes []*domain.Note) []noteResponse {
	responses := make([]noteResponse, 0, len(notes))
	for _, n := range notes {
		responses = append(responses, toNoteResponse(n))
	}
	return responses
}
