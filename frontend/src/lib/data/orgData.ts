import type { Org } from "@/lib/types.ts";
import { BACKEND_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";

export async function getOrg(orgName: string): Promise<Org | null> {
  const res = await authFetch(`${BACKEND_URL}/orgs/${orgName}`);

  if (!res.ok) {
    throw new Error("Failed to fetch org");
  }

  const org = (await res.json()) as Org;

  org.updatedAt = new Date(org.updatedAt);
  org.createdAt = new Date(org.createdAt);
  org.members = org.members.map((member) => ({
    ...member,
    createdAt: new Date(member.createdAt),
    updatedAt: new Date(member.updatedAt),
  }));

  return org;
}

export async function fetchOrgs(): Promise<Org[]> {
  const res = await authFetch(`${BACKEND_URL}/orgs`);

  if (!res.ok) {
    throw new Error("Failed to fetch orgs");
  }

  return await res.json();
}

