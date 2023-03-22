package dal

import "note/cmd/note/dal/db"

// Init init dal
func Init() {
	db.Init() // mysql init
}
