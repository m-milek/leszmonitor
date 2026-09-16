import type { Timestamps } from "@/lib/types";

export enum UserRole {
  Owner = "owner",
  Admin = "admin",
  Member = "member",
  Viewer = "viewer",
}

export const mapUserRoleToDisplayName: Record<UserRole, string> = {
  [UserRole.Owner]: "Owner",
  [UserRole.Admin]: "Admin",
  [UserRole.Member]: "Member",
  [UserRole.Viewer]: "Viewer",
};

export interface User extends Timestamps {
  id: string;
  username: string;
  role: UserRole;
}
