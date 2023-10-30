package model

import "time"

type User struct {
	ID       string
	Username string
	Password string
	Role     string
}

type Slot struct {
	ID        string
	Start     time.Time
	End       time.Time
	DoctorId  string
	PatientId string
}

type SlotXDoctorXPatient struct {
	Slot
	PatientName string
	DoctorName  string
}
