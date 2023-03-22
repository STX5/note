package dal

import "note/cmd/user/dal/db"

// Init init dal
func Init() {
	db.Init() // mysql init
}
