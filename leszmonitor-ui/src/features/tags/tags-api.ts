import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";
import type { Tag } from "@/features/tags/types";

export interface TagPayload {
  name: string;
  description: string;
  colorHex: string;
}

const mapTag = (tag: Tag): Tag => {
  return {
    ...tag,
    createdAt: new Date(tag.createdAt),
    updatedAt: new Date(tag.updatedAt),
  };
};

const getAll = async (): Promise<Tag[]> => {
  const res = await authFetch(`${SERVER_API_URL}/tags`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const tags = (await res.json()) as Tag[];
  return tags.map(mapTag);
};

const getById = async (tagId: string): Promise<Tag> => {
  const res = await authFetch(`${SERVER_API_URL}/tags/${tagId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const tag = (await res.json()) as Tag;
  return mapTag(tag);
};

const create = async (tagData: TagPayload): Promise<Tag> => {
  const res = await authFetch(`${SERVER_API_URL}/tags`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(tagData),
  });

  const tag = (await res.json()) as Tag;
  return mapTag(tag);
};

const update = async (tagId: string, tagData: TagPayload): Promise<Tag> => {
  const res = await authFetch(`${SERVER_API_URL}/tags/${tagId}`, {
    method: "PATCH",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(tagData),
  });

  const tag = (await res.json()) as Tag;
  return mapTag(tag);
};

const remove = async (tagId: string): Promise<void> => {
  await authFetch(`${SERVER_API_URL}/tags/${tagId}`, {
    method: "DELETE",
  });
};

export const TagsApi = {
  getAll,
  getById,
  create,
  update,
  remove,
};
