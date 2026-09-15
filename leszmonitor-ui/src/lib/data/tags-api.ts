import { SERVER_API_URL } from "@/lib/consts.ts";
import { authFetch } from "@/lib/data/utils.ts";
import type { Tag } from "@/lib/types.ts";

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

export const getAllTags = async (): Promise<Tag[]> => {
  const res = await authFetch(`${SERVER_API_URL}/tags`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const tags = (await res.json()) as Tag[];
  return tags.map(mapTag);
};

export const getTagById = async (tagId: string): Promise<Tag> => {
  const res = await authFetch(`${SERVER_API_URL}/tags/${tagId}`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const tag = (await res.json()) as Tag;
  return mapTag(tag);
};

export const createTag = async (tagData: TagPayload): Promise<Tag> => {
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

export const updateTag = async (
  tagId: string,
  tagData: TagPayload,
): Promise<Tag> => {
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

export const deleteTag = async (tagId: string): Promise<void> => {
  await authFetch(`${SERVER_API_URL}/tags/${tagId}`, {
    method: "DELETE",
  });
};
