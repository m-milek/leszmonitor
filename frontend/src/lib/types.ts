export interface User {
  id: string;
  username: string;
  createdAt: Date;
  updatedAt: Date;
}

export interface LoginPayload {
  username: string;
  password: string;
}

export interface LoginResponse {
  jwt: string;
}

export interface JwtClaims {
  username: string;
}
