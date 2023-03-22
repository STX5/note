package service

import (
	"context"
	"note/cmd/note/dal/db"
	"note/kitex_gen/note"
)

type UpdateNoteService struct {
	ctx context.Context
}

// NewUpdateNoteService new UpdateNoteService
func NewUpdateNoteService(ctx context.Context) *UpdateNoteService {
	return &UpdateNoteService{ctx: ctx}
}

// UpdateNote update note info
func (s *UpdateNoteService) UpdateNote(req *note.UpdateNoteRequest) error {
	return db.UpdateNote(s.ctx, req.NoteId, req.UserId, req.Title, req.Content)
}
