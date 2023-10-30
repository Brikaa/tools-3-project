package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Brikaa/tools-3-project/src/backend/model"
)

func selectOne(db *sql.DB, query string, arguments []any, rows []any) error {
	err := db.QueryRow(query, arguments...).Scan(rows...)
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return fmt.Errorf("%v, %v, %v: %v", query, rows, arguments, err)
	}
	return nil
}

func insert(db *sql.DB, query string, arguments []any) error {
	result, err := db.Exec(query, arguments...)
	if err != nil {
		return fmt.Errorf("%v", err)
	}
	if _, err := result.LastInsertId(); err != nil {
		return fmt.Errorf("%v", err)
	}
	return nil
}

func update(db *sql.DB, query string, arguments []any) (bool, error) {
	result, err := db.Exec(query, arguments...)
	if err != nil {
		return false, fmt.Errorf("%v", err)
	}
	rowsAffected, rErr := result.RowsAffected()
	if rErr != nil {
		return false, fmt.Errorf("%v", err)
	}
	return rowsAffected >= 1, nil
}

func selectOneUser(db *sql.DB, condition string, arguments []any) (*model.User, error) {
	var user model.User
	return &user, selectOne(db,
		"SELECT id, username, password, role FROM User WHERE "+condition,
		arguments, []any{&user.ID, &user.Username, &user.Password, &user.Role})
}

func SelectUserByUsername(db *sql.DB, username string) (*model.User, error) {
	return selectOneUser(db, "username = ?", []any{username})
}

func SelectUserByUsernameAndPassword(db *sql.DB, username string, password string) (*model.User, error) {
	return selectOneUser(db, "username = ? AND password = ?", []any{username, password})
}

func InsertUser(db *sql.DB, username, password, role string) error {
	return insert(
		db,
		"INSERT INTO User (username, password, role) VALUES (?, ?, ?)",
		[]any{username, password, role},
	)
}

func GetOverlappingSlotId(db *sql.DB, doctorId string, start time.Time, end time.Time) (*int64, error) {
	var slotId int64
	return &slotId, selectOne(
		db,
		"SELECT id FROM Slot WHERE doctorId = ? AND ? >= start AND ? <= end",
		[]any{doctorId, start, end},
		[]any{&slotId},
	)
}

func InsertSlot(db *sql.DB, start time.Time, end time.Time, doctorId string) error {
	return insert(
		db,
		"INSERT INTO Slot (start, end, doctorId) VALUES (?, ?, ?)",
		[]any{start, end, doctorId},
	)
}

func DeleteSlotByIdAndDoctorId(db *sql.DB, slotId string, doctorId string) (bool, error) {
	return update(db, "DELETE FROM Slot WHERE id = ? AND doctorId = ?", []any{slotId, doctorId})
}

func GetSlotsByDoctorId(db *sql.DB, doctorId string) ([]*model.SlotXReserved, error) {
	var slots []*model.SlotXReserved

	rows, err := db.Query(
		`SELECT Slot.id, Slot.start, Slot.end, (Appointment.id IS NOT NULL) AS reserved FROM Slot WHERE doctorId = ?
LEFT JOIN Appointment ON Appointment.slotId = Slot.id ORDER BY Slot.start`,
		doctorId,
	)
	if err != nil {
		return nil, fmt.Errorf("%q: %v", doctorId, err)
	}
	defer rows.Close()

	for rows.Next() {
		var slot model.SlotXReserved
		if err := rows.Scan(&slot.ID, &slot.Start, &slot.End, &slot.Reserved); err != nil {
			return nil, fmt.Errorf("%q: %v", doctorId, err)
		}
		slots = append(slots, &slot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%q: %v", doctorId, err)
	}
	return slots, nil
}
