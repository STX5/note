package api_service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	api "note/hertz_gen/api"
	"testing"
)

var (
	defaultHeader = http.Header{}
	token         string
)

func TestCreatUser(t *testing.T) {
	req := api.CreateUserRequest{
		Username: "lorain",
		Password: "123456",
	}
	resp, rawResp, err := CreateUser(context.Background(), &req)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(resp)
	fmt.Println(string(rawResp.Body()))
}

// Login
func TestCheckUser(t *testing.T) {
	req := api.CheckUserRequest{
		Username: "lorain",
		Password: "123456",
	}
	_, rawResp, err := CheckUser(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	var res map[string]interface{}
	err = json.Unmarshal(rawResp.Body(), &res)
	if err != nil {
		t.Error(err)
	}
	// JWT Token
	token = res["token"].(string)
	defaultHeader.Add("Authorization", "Bearer "+token)
	fmt.Printf("token:%s\n", token)
}

func TestCreateNote(t *testing.T) {
	TestCheckUser(t)
	authorizationClient, _ := NewApiServiceClient("http://127.0.0.1:8080", WithHeader(defaultHeader))
	req := api.CreateNoteRequest{
		Title:   "test title",
		Content: "test content",
	}
	resp, rawResp, err := authorizationClient.CreateNote(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	fmt.Println(string(rawResp.Body()))
}

func TestQueryNote(t *testing.T) {
	TestCheckUser(t)
	authorizationClient, _ := NewApiServiceClient("http://127.0.0.1:8080", WithHeader(defaultHeader))
	key := "test"
	req := api.QueryNoteRequest{
		Offset:    0,
		Limit:     20,
		SearchKey: &key,
	}
	_, rawResp, err := authorizationClient.QueryNote(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rawResp.Body()))
}

func TestUpdateNote(t *testing.T) {
	TestCheckUser(t)
	authorizationClient, _ := NewApiServiceClient("http://127.0.0.1:8080", WithHeader(defaultHeader))
	title := "test"
	content := "test"
	req := api.UpdateNoteRequest{
		Title:   &title,
		Content: &content,
		NoteID:  1,
	}
	resp, rawResp, err := authorizationClient.UpdateNote(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(resp)
	fmt.Println(string(rawResp.Body()))
}

func TestDeleteNote(t *testing.T) {
	TestCheckUser(t)
	authorizationClient, _ := NewApiServiceClient("http://127.0.0.1:8080", WithHeader(defaultHeader))
	req := api.DeleteNoteRequest{
		NoteID: 1,
	}
	_, rawResp, err := authorizationClient.DeleteNote(context.Background(), &req)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(rawResp.Body()))
}
