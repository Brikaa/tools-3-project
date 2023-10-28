DROP TABLE IF EXISTS User
DROP TABLE IF EXISTS Slot

CREATE TABLE User (
  `id` BINARY(16) DEFAULT (UUID_TO_BIN(UUID())) PRIMARY KEY,
  `name` varchar(255) UNIQUE NOT NULL,
  `password` varchar(255) NOT NULL,
  `role` ENUM('patient', 'doctor') NOT NULL
)

CREATE TABLE Slot (
  `id` BINARY(16) DEFAULT (UUID_TO_BIN(UUID())) PRIMARY KEY,
  `start` TIMESTAMP NOT NULL,
  `end` TIMESTAMP NOT NULL,
  doctorId BINARY(16),
  patientId BINARY(16),
  CONSTRAINT FK_SLOT_DOCTOR FOREIGN KEY (doctorId) REFERENCES User(id) ON DELETE CASCADE,
  CONSTRAINT FK_SLOT_DOCTOR FOREIGN KEY (patientId) REFERENCES User(id)
)
