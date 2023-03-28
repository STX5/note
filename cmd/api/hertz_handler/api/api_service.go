// Code generated by hertz generator.

package api

import (
	"context"

	"note/cmd/api/middleware"
	"note/cmd/api/rpc"
	api "note/hertz_gen/api"
	"note/kitex_gen/note"
	"note/kitex_gen/user"
	"note/pkg/consts"
	"note/pkg/errno"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
)

// CreateUser .
// @router /v1/user/register [POST]
func CreateUser(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.CreateUserRequest
	// 验证c里面的数据是否可以绑定到req中，并验证约束
	err = c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	err = rpc.CreateUser(context.Background(), &user.CreateUserRequest{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	SendResponse(c, errno.Success, nil)
}

// CheckUser .
// @router /v1/user/login [POST]
func CheckUser(ctx context.Context, c *app.RequestContext) {
	middleware.JwtMiddleware.LoginHandler(ctx, c)
}

// CreateNote .
// @router /v1/note [POST]
func CreateNote(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.CreateNoteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	v, _ := c.Get(consts.IdentityKey)
	err = rpc.CreateNote(context.Background(), &note.CreateNoteRequest{
		Title:   req.Title,
		Content: req.Content,
		UserId:  v.(*api.User).UserID,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	SendResponse(c, errno.Success, nil)
}

// QueryNote .
// @router /v1/note/query [GET]
func QueryNote(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.QueryNoteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	v, _ := c.Get(consts.IdentityKey)
	notes, total, err := rpc.QueryNotes(context.Background(), &note.QueryNoteRequest{
		UserId:    v.(*api.User).UserID,
		SearchKey: req.SearchKey,
		Offset:    req.Offset,
		Limit:     req.Limit,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	SendResponse(c, errno.Success, utils.H{
		consts.Total: total,
		consts.Notes: notes,
	})
}

// UpdateNote .
// @router /v1/note/:note_id [PUT]
func UpdateNote(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.UpdateNoteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	v, _ := c.Get(consts.IdentityKey)
	err = rpc.UpdateNote(context.Background(), &note.UpdateNoteRequest{
		NoteId:  req.NoteID,
		UserId:  v.(*api.User).UserID,
		Title:   req.Title,
		Content: req.Content,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	SendResponse(c, errno.Success, nil)
}

// DeleteNote .
// @router /v1/note/:note_id [DELETE]
func DeleteNote(ctx context.Context, c *app.RequestContext) {
	var err error
	var req api.DeleteNoteRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	v, _ := c.Get(consts.IdentityKey)
	err = rpc.DeleteNote(context.Background(), &note.DeleteNoteRequest{
		NoteId: req.NoteID,
		UserId: v.(*api.User).UserID,
	})
	if err != nil {
		SendResponse(c, errno.ConvertErr(err), nil)
		return
	}
	SendResponse(c, errno.Success, nil)
}
