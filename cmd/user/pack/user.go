package pack

import (
	"note/cmd/user/dal/db"
	"note/kitex_gen/user"
)

// User pack user info
func User(u *db.User) *user.User {
	if u == nil {
		return nil
	}

	return &user.User{UserId: int64(u.ID), Username: u.Username, Avatar: "test"}
}

// Users pack list of user info
func Users(us []*db.User) []*user.User {
	users := make([]*user.User, 0)
	for _, u := range us {
		if temp := User(u); temp != nil {
			users = append(users, temp)
		}
	}
	return users
}
