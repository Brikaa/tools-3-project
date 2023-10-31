package model

import "time"

type User struct {
	ID       string
	Username string
	Password string
	Role     string
}

type Doctor struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type Slot struct {
	ID       string    `json:"id"`
	Start    time.Time `json:"start"`
	End      time.Time `json:"end"`
	DoctorID string    `json:"doctorId"`
}

type Appointment struct {
	ID     string `json:"id"`
	SlotID string `json:"slotId"`
}

type AppointmentXSlot struct {
	Appointment
	SlotStart time.Time `json:"start"`
	SlotEnd   time.Time `json:"end"`
}

type AppointmentXSlotXPatient struct {
	AppointmentXSlot
	PatientID       string `json:"patientId"`
	PatientUsername string `json:"patientUsername"`
}

type AppointmentXSlotXDoctor struct {
	AppointmentXSlot
	DoctorID       string
	DoctorUsername string
}
