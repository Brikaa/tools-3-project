DROP DATABASE IF EXISTS app;
CREATE DATABASE IF NOT EXISTS app;
USE app;

CREATE TABLE User (
  `id` varchar(36) DEFAULT (UUID()) PRIMARY KEY,
  `username` varchar(255) UNIQUE NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` ENUM('patient', 'doctor') NOT NULL
);

CREATE TABLE Slot (
  `id` varchar(36) DEFAULT (UUID()) PRIMARY KEY,
  `start` TIMESTAMP NOT NULL,
  `end` TIMESTAMP NOT NULL,
  doctorId varchar(36) NOT NULL,
  CONSTRAINT FK_SLOT_DOCTOR FOREIGN KEY (doctorId) REFERENCES User(id) ON DELETE CASCADE,
);

CREATE TABLE Appointment (
  `id` varchar(36) DEFAULT (UUID()) PRIMARY KEY,
  slotId varchar(36) UNIQUE NOT NULL,
  patientId varchar(36) NOT NULL,
  CONSTRAINT FK_APPOINTMENT_PATIENT FOREIGN KEY (patientId) REFERENCES User(id) ON DELETE CASCADE,
  CONSTRAINT FK_APPOINTMENT_SLOT FOREIGN KEY (slotId) REFERENCES Slot(id) ON DELETE CASCADE,
)
