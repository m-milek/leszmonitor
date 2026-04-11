import { BACKEND_URL } from "@/lib/consts.ts";
import type { Project, Timestamps } from "@/lib/types.ts";
import { authFetch } from "@/lib/data/utils.ts";

export type ProjectInput = Omit<Project, "id" | keyof Timestamps>;

export const getProjects = async (orgId: string): Promise<Project[]> => {
  const res = await authFetch(`${BACKEND_URL}/orgs/${orgId}/projects`);

  if (!res.ok) {
    throw new Error("Failed to fetch projects");
  }

  const projects = (await res.json()) as Project[];
  return projects.map((project) => ({
    ...project,
    createdAt: new Date(project.createdAt),
    updatedAt: new Date(project.updatedAt),
  }));
};

export const getProject = async (
  orgId: string,
  projectId: string,
): Promise<Project> => {
  const res = await authFetch(
    `${BACKEND_URL}/orgs/${orgId}/projects/${projectId}`,
  );

  if (!res.ok) {
    throw new Error("Failed to fetch project");
  }

  const project = (await res.json()) as Project;
  return {
    ...project,
    createdAt: new Date(project.createdAt),
    updatedAt: new Date(project.updatedAt),
  };
};

export const addProject = async (
  orgId: string,
  project: ProjectInput,
): Promise<Project> => {
  const res = await authFetch(`${BACKEND_URL}/orgs/${orgId}/projects`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(project),
  });

  if (!res.ok) {
    throw new Error("Failed to add project");
  }

  const newProject = (await res.json()) as Project;
  return {
    ...newProject,
    createdAt: new Date(newProject.createdAt),
    updatedAt: new Date(newProject.updatedAt),
  };
};

export const deleteProject = async (
  orgId: string,
  projectId: string,
): Promise<void> => {
  const res = await authFetch(
    `${BACKEND_URL}/orgs/${orgId}/projects/${projectId}`,
    {
      method: "DELETE",
    },
  );

  if (!res.ok) {
    throw new Error("Failed to delete project");
  }
};

