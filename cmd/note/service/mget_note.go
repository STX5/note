package service

import (
	"context"
	"note/cmd/note/dal/db"
	"note/cmd/note/pack"
	"note/cmd/note/rpc"
	"note/kitex_gen/note"
	"note/kitex_gen/user"
)

type MGetNoteService struct {
	ctx context.Context
}

// NewMGetNoteService new MGetNoteService
func NewMGetNoteService(ctx context.Context) *MGetNoteService {
	return &MGetNoteService{ctx: ctx}
}

// MGetNote multiple get list of note info
func (s *MGetNoteService) MGetNote(req *note.MGetNoteRequest) ([]*note.Note, error) {
	// 从数据库查询一组Note ID，返回一组noteModel
	noteModels, err := db.MGetNotes(s.ctx, req.NoteIds)
	if err != nil {
		return nil, err
	}
	// 从所有notemodel中获取用户IDs
	uIds := pack.UserIds(noteModels)
	// 根据用户id 查询用户信息
	userMap, err := rpc.MGetUser(s.ctx, &user.MGetUserRequest{UserIds: uIds})
	if err != nil {
		return nil, err
	}
	notes := pack.Notes(noteModels)
	// 为每一个Note添加用户信息
	for i := 0; i < len(notes); i++ {
		if u, ok := userMap[notes[i].UserId]; ok {
			notes[i].Username = u.Username
			notes[i].UserAvatar = u.Avatar
		}
	}
	return notes, nil
}
