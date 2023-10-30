package model

import "time"

type User struct {
	ID       string
	Username string
	Password string
	Role     string
}

func CreateUser(username, password, role string) *User {
	return &User{Username: username, Password: password, Role: role}
}

type Slot struct {
	ID        string
	Start     time.Time
	end       time.Time
	doctorId  string
	patientId string
}
