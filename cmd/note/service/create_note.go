package service

import (
	"context"
	"note/cmd/note/dal/db"
	"note/kitex_gen/note"
)

type CreateNoteService struct {
	ctx context.Context
}

// NewCreateNoteService new CreateNoteService
func NewCreateNoteService(ctx context.Context) *CreateNoteService {
	return &CreateNoteService{ctx: ctx}
}

// CreateNote create note info
func (s *CreateNoteService) CreateNote(req *note.CreateNoteRequest) error {
	noteModel := &db.Note{
		UserID:  req.UserId,
		Title:   req.Title,
		Content: req.Content,
	}
	return db.CreateNote(s.ctx, []*db.Note{noteModel})
}
