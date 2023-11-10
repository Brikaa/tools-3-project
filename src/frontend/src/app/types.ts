type Role = 'patient' | 'doctor';

export interface UserContext {
  id: string;
  username: string;
  role: Role;
  token: string;
}

export interface BadRequestResponse {
  message: string;
}

export interface User {
  id: string;
  username: string;
  role: Role;
}

export type Doctor = Omit<User, 'role'>;

export interface Slot {
  id: string;
  start: string;
  end: string;
}

interface Appointment {
  id: string;
  slotId: string;
  start: string;
  end: string;
}

export type PatientAppointment = Appointment & {
  doctorId: string;
  doctorUsername: string;
};

export type DoctorAppointment = Appointment & {
  patientId: string;
  patientUsername: string;
};
