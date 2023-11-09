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

export interface PatientAppointment {
  id: string;
  slotStart: Date;
  slotEnd: Date;
  doctorId: string;
  doctorUsername: string;
}
