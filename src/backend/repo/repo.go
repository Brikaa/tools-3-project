package repo

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Brikaa/tools-3-project/src/backend/model"
)

func selectOne[T any](db *sql.DB, query string, arguments []any, flatten func(*T) []any) (*T, error) {
	var entity T
	if err := db.QueryRow(query, arguments...).Scan(flatten(&entity)...); err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, fmt.Errorf("%v, %v: %v", query, arguments, err)
	}
	return &entity, nil
}

func insert(db *sql.DB, query string, arguments []any) error {
	if _, err := db.Exec(query, arguments...); err != nil {
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

func selectAll[T any](
	db *sql.DB,
	query string,
	arguments []any,
	fn func(*sql.Rows, *T) error,
) ([]*T, error) {
	var entities []*T

	rows, err := db.Query(query, arguments...)
	if err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var entity T
		if err := fn(rows, &entity); err != nil {
			return nil, fmt.Errorf("%v", err)
		}
		entities = append(entities, &entity)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%v", err)
	}
	return entities, nil
}

func selectOneUser(db *sql.DB, condition string, arguments []any) (*model.User, error) {
	return selectOne(db,
		"SELECT id, username, password, role FROM User WHERE "+condition,
		arguments, func(user *model.User) []any { return []any{&user.ID, &user.Username, &user.Password, &user.Role} })
}

func GetUserByUsername(db *sql.DB, username string) (*model.User, error) {
	return selectOneUser(db, "username = ?", []any{username})
}

func GetUserByUsernameAndPassword(db *sql.DB, username string, password string) (*model.User, error) {
	return selectOneUser(db, "username = ? AND password = ?", []any{username, password})
}

func GetUserByIdAndPassword(db *sql.DB, id string, password string) (*model.User, error) {
	return selectOneUser(db, "id = ? AND password = ?", []any{id, password})
}

func InsertUser(db *sql.DB, username, password, role string) error {
	return insert(
		db,
		"INSERT INTO User (username, password, role) VALUES (?, ?, ?)",
		[]any{username, password, role},
	)
}

func GetOverlappingSlotId(db *sql.DB, doctorId string, start time.Time, end time.Time) (*string, error) {
	return selectOne(
		db,
		"SELECT id FROM Slot WHERE doctorId = ? AND ((? >= start AND ? <= end) OR (? >= start AND ? <= end))",
		[]any{doctorId, start, start, end, end},
		func(slotId *string) []any { return []any{slotId} },
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

func GetSlotsByDoctorId(db *sql.DB, doctorId string) ([]*model.Slot, error) {
	return selectAll(
		db,
		"SELECT Slot.id, Slot.start, Slot.end, Slot.doctorId FROM Slot WHERE doctorId = ?",
		[]any{doctorId},
		func(rows *sql.Rows, slot *model.Slot) error {
			return rows.Scan(&slot.ID, &slot.Start, &slot.End, &slot.DoctorID)
		},
	)
}

func GetAppointmentsByDoctorId(db *sql.DB, doctorId string) ([]*model.AppointmentXSlotXPatient, error) {
	return selectAll(
		db,
		`SELECT Appointment.id, Slot.id, Slot.start, Slot.end, Patient.id, Patient.username
FROM Appointment WHERE Doctor.id = ?
LEFT JOIN Slot ON Appointment.slotId = Slot.id
LEFT JOIN User AS Patient ON Appointment.patientId = Patient.id
LEFT JOIN User AS Doctor ON Slot.doctorId = Doctor.id`,
		[]any{doctorId},
		func(rows *sql.Rows, appointment *model.AppointmentXSlotXPatient) error {
			return rows.Scan(
				&appointment.ID,
				&appointment.SlotID,
				&appointment.SlotStart,
				&appointment.SlotEnd,
				&appointment.PatientID,
				&appointment.PatientUsername,
			)
		},
	)
}

func GetAppointmentsByPatientId(db *sql.DB, patientId string) ([]*model.AppointmentXSlotXDoctor, error) {
	return selectAll(
		db,
		`SELECT Appointment.id, Slot.id, Slot.start, Slot.end, Doctor.id, Doctor.username
FROM Appointment WHERE Patient.id = ?
LEFT JOIN Slot ON Appointment.slotID = Slot.id
LEFT JOIN User AS Patient ON Appointment.patientId = Patient.id
LEFT JOIN USER AS Doctor ON Slot.doctorId = Doctor.id`,
		[]any{patientId},
		func(rows *sql.Rows, appointment *model.AppointmentXSlotXDoctor) error {
			return rows.Scan(
				&appointment.ID,
				&appointment.SlotID,
				&appointment.SlotStart,
				&appointment.SlotEnd,
				&appointment.DoctorID,
				&appointment.DoctorUsername,
			)
		},
	)
}

func GetAppointmentIdBySlotId(db *sql.DB, slotId string) (*string, error) {
	return selectOne(
		db,
		"SELECT Appointment.id FROM Appointment WHERE Appointment.slotId = ?",
		[]any{slotId},
		func(appointmentId *string) []any { return []any{appointmentId} },
	)
}

func GetSlotIdBySlotId(db *sql.DB, slotId string) (*string, error) {
	return selectOne(
		db,
		"SELECT Slot.id FROM Slot WHERE slotId = ?",
		[]any{slotId},
		func(slotId *string) []any { return []any{slotId} },
	)
}

func InsertAppointment(db *sql.DB, slotId string, patientId string) error {
	return insert(
		db,
		"INSERT INTO Appointment (slotId, patientId) VALUES (?, ?)",
		[]any{slotId, patientId},
	)
}

func DeleteAppointmentByIdAndPatientId(db *sql.DB, appointmentId string, patientId string) (bool, error) {
	return update(db, "DELETE FROM Appointment WHERE id = ? AND patientId = ?", []any{appointmentId, patientId})
}

func GetDoctors(db *sql.DB) ([]*model.Doctor, error) {
	return selectAll(
		db,
		`SELECT User.id, User.username FROM User WHERE User.role = "doctor"`,
		[]any{},
		func(rows *sql.Rows, doctor *model.Doctor) error {
			return rows.Scan(&doctor.ID, &doctor.Username)
		},
	)
}

func GetAvailableSlotsByDoctorId(db *sql.DB, doctorId string) ([]*model.Slot, error) {
	return selectAll(
		db,
		`SELECT Slot.id, Slot.start, Slot.end, Slot.doctorId FROM Slot WHERE doctorId = ? AND Appointment.id IS NULL
LEFT JOIN Appointment ON Appointment.slotId = Slot.id`,
		[]any{doctorId},
		func(rows *sql.Rows, slot *model.Slot) error {
			return rows.Scan(&slot.ID, &slot.Start, &slot.End)
		},
	)
}

func UpdateSlotByIdAndDoctorId(
	db *sql.DB, slotId string, doctorId string, start time.Time, end time.Time,
) (bool, error) {
	return update(
		db,
		"UPDATE Slot SET start = ?, end = ? WHERE id = ? AND doctorId = ?",
		[]any{slotId, doctorId, start, end},
	)
}

func UpdateAppointmentByIdAndPatientId(
	db *sql.DB, appointmentId string, patientId string, slotId string,
) (bool, error) {
	return update(
		db,
		"UPDATE Appointment SET slotId = ? WHERE id = ? AND patientId = ?",
		[]any{appointmentId, patientId, slotId},
	)
}
