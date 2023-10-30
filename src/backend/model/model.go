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
	DoctorId string
}

type SlotXReserved struct {
	Slot
	Reserved bool
}
