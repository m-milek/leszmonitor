export interface Timestamps {
  createdAt: Date;
  updatedAt: Date;
}

export interface User extends Timestamps {
  id: string;
  username: string;
}

export enum TeamRole {
  Owner = "owner",
  Admin = "admin",
  Member = "member",
  Viewer = "viewer",
}

export interface TeamMember extends Timestamps {
  id: string;
  role: TeamRole;
}

export interface Team extends Timestamps {
  id: string;
  displayId: string;
  name: string;
  description: string;
  members: TeamMember;
}

export interface Group extends Timestamps {
  id: string;
  name: string;
  displayId: string;
  description: string;
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
