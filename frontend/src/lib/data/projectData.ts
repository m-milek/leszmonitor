import { BACKEND_URL } from "@/lib/consts.ts";
import type {
  Project,
  ProjectMember,
  ProjectRole,
  Timestamps,
} from "@/lib/types.ts";
import { authFetch } from "@/lib/data/utils.ts";

export type ProjectInput = Omit<Project, "id" | "members" | keyof Timestamps>;

const parseProject = (project: Project): Project => ({
  ...project,
  createdAt: new Date(project.createdAt),
  updatedAt: new Date(project.updatedAt),
  members:
    project.members?.map((m) => ({
      ...m,
      createdAt: new Date(m.createdAt),
      updatedAt: new Date(m.updatedAt),
    })) ?? [],
});

export const getProjects = async (): Promise<Project[]> => {
  const res = await authFetch(`${BACKEND_URL}/projects`);

  const projects = (await res.json()) as Project[];
  return projects.map(parseProject);
};

export const getProject = async (projectId: string): Promise<Project> => {
  const res = await authFetch(`${BACKEND_URL}/projects/${projectId}`);

  return parseProject((await res.json()) as Project);
};

export const addProject = async (project: ProjectInput): Promise<Project> => {
  const res = await authFetch(`${BACKEND_URL}/projects`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(project),
  });

  return parseProject((await res.json()) as Project);
};

export const deleteProject = async (projectId: string): Promise<void> => {
  const res = await authFetch(`${BACKEND_URL}/projects/${projectId}`, {
    method: "DELETE",
  });

  if (!res.ok) {
    throw new Error("Failed to delete project");
  }
};

export interface AddProjectMemberPayload {
  username: string;
  role: ProjectRole;
}

export const addMemberToProject = async (
  projectId: string,
  payload: AddProjectMemberPayload,
): Promise<void> => {
  const res = await authFetch(`${BACKEND_URL}/projects/${projectId}/members`, {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    throw new Error("Failed to add member to project");
  }
};

export interface RemoveProjectMemberPayload {
  username: string;
}

export const removeMemberFromProject = async (
  projectId: string,
  payload: RemoveProjectMemberPayload,
): Promise<void> => {
  const res = await authFetch(`${BACKEND_URL}/projects/${projectId}/members`, {
    method: "DELETE",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(payload),
  });

  if (!res.ok) {
    throw new Error("Failed to remove member from project");
  }
};

export type { ProjectMember };
