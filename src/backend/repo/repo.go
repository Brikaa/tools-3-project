package repo

import (
	"database/sql"
	"fmt"

	"github.com/Brikaa/tools-3-project/src/backend/model"
)

func selectOneUser(db *sql.DB, condition string, arguments []any) (*model.User, error) {
	var user model.User
	err := db.QueryRow(
		"SELECT id, username, password, role FROM User WHERE "+condition, arguments...).Scan(
		&user.ID, &user.Username, &user.Password, &user.Role)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("%v, %v: %v", condition, arguments, err)
	}
	return &user, nil
}

func SelectUserByUsername(db *sql.DB, username string) (*model.User, error) {
	return selectOneUser(db, "username = ?", []any{username})
}

func SelectUserByUsernameAndPassword(db *sql.DB, username string, password string) (*model.User, error) {
	return selectOneUser(db, "username = ? AND password = ?", []any{username, password})
}

func InsertUser(db *sql.DB, user model.User) error {
	result, err := db.Exec(
		"INSERT INTO User (username, password, role) VALUES (?, ?, ?)", user.Username, user.Password, user.Role)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if _, err := result.LastInsertId(); err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}
