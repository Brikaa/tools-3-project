package model

import "time"

type User struct {
	ID       string
	Username string
	Password string
	Role     string
}

type Slot struct {
	ID       string
	Start    time.Time
	End      time.Time
	DoctorID string
}

type Appointment struct {
	ID        string
	SlotID    string
	PatientID string
}

type AppointmentXSlotXPatient struct {
	Appointment
	SlotStart       time.Time
	SlotEnd         time.Time
	PatientUsername string
}
