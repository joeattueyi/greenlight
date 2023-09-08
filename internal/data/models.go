package data

import (
	"database/sql"
	"errors"
)

var (
	ErrRecordNotFound = errors.New("record not found")
	ErrEditConflict   = errors.New("edit conflict")
)

type Models struct {
	Movies      MovieModel
	Users       UserModel
	Tokens      TokenModel
	Permissions PermissionModel
}

func NewModels(db *sql.DB) Models {
	return Models{
		Movies:      DBMovieModel{DB: db},
		Users:       DBUserModel{DB: db},
		Tokens:      DBTokenModel{DB: db},
		Permissions: DBPermissionModel{DB: db},
	}

}
