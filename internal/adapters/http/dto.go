package httpadapter

import (
	"time"

	domain "notes-service/internal/domain/note"
)

type createNoteRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type updateNoteRequest struct {
	Title   string   `json:"title"`
	Content string   `json:"content"`
	Tags    []string `json:"tags"`
}

type noteResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Tags      []string  `json:"tags"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// toNoteResponse is the single owner of the "tags is never null" JSON
// guarantee: a nil slice (e.g. from legacy documents) becomes [].
func toNoteResponse(n *domain.Note) noteResponse {
	tags := make([]string, len(n.Tags))
	copy(tags, n.Tags)
	return noteResponse{
		ID:        n.ID,
		Title:     n.Title,
		Content:   n.Content,
		Tags:      tags,
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
