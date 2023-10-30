package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Brikaa/tools-3-project/src/backend/model"
)

func selectOne[T any](db *sql.DB, query string, arguments []any, entity *T, rows []any) (*T, error) {
	err := db.QueryRow(query, arguments...).Scan(rows...)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("%v, %v, %v: %v", query, rows, arguments, err)
	}
	return entity, nil
}

func selectOneUser(db *sql.DB, condition string, arguments []any) (*model.User, error) {
	var user model.User
	return selectOne(db,
		"SELECT id, username, password, role FROM User WHERE "+condition,
		arguments, &user, []any{&user.ID, &user.Username, &user.Password, &user.Role})
}

func SelectUserByUsername(db *sql.DB, username string) (*model.User, error) {
	return selectOneUser(db, "username = ?", []any{username})
}

func SelectUserByUsernameAndPassword(db *sql.DB, username string, password string) (*model.User, error) {
	return selectOneUser(db, "username = ? AND password = ?", []any{username, password})
}

func InsertUser(db *sql.DB, user *model.User) error {
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

func GetOverlappingSlot(db *sql.DB, doctorId string, start time.Time, end time.Time) (*model.Slot, error) {
	var slot model.Slot
	return selectOne(
		db,
		"SELECT id FROM Slot WHERE doctorId = ? AND ? >= start AND ? <= end",
		[]any{doctorId, start, end},
		&slot,
		[]any{&slot.ID},
	)
}
