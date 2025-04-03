package database

import "database/sql"

type Models struct {
	Users     UserModel
	Friends   FriendModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Users:     UserModel{DB: db},
		Friends:   FriendModel{DB: db},
	}
}