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
