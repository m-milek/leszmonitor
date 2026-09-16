export interface LoginPayload {
  username: string;
  password: string;
}

export interface LoginResponse {
  jwt: string;
}
