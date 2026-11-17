import { SERVER_API_URL } from "@/lib/consts";
import { authFetch } from "@/lib/api-client";
import type {
  NotificationChannel,
  NotificationChannelConfig,
} from "@/features/notifications/types";

export type NotificationChannelPayload = NotificationChannelConfig & {
  name: string;
  enabled: boolean;
};

const mapChannel = (channel: NotificationChannel): NotificationChannel => {
  return {
    ...channel,
    createdAt: new Date(channel.createdAt),
    updatedAt: new Date(channel.updatedAt),
  };
};

const getAll = async (): Promise<NotificationChannel[]> => {
  const res = await authFetch(`${SERVER_API_URL}/notification-channels`, {
    method: "GET",
    headers: {
      "Content-Type": "application/json",
    },
  });

  const channels = (await res.json()) as NotificationChannel[];
  return channels.map(mapChannel);
};

const create = async (
  payload: NotificationChannelPayload,
): Promise<NotificationChannel> => {
  const res = await authFetch(`${SERVER_API_URL}/notification-channels`, {
    method: "POST",
    headers: {
      "Content-Type": "application/json",
    },
    body: JSON.stringify(payload),
  });

  const channel = (await res.json()) as NotificationChannel;
  return mapChannel(channel);
};

const update = async (
  channelId: string,
  payload: NotificationChannelPayload,
): Promise<NotificationChannel> => {
  const res = await authFetch(
    `${SERVER_API_URL}/notification-channels/${channelId}`,
    {
      method: "PATCH",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify(payload),
    },
  );

  const channel = (await res.json()) as NotificationChannel;
  return mapChannel(channel);
};

const remove = async (channelId: string): Promise<void> => {
  await authFetch(`${SERVER_API_URL}/notification-channels/${channelId}`, {
    method: "DELETE",
  });
};

const exportYaml = async (channelId: string): Promise<string> => {
  const res = await authFetch(
    `${SERVER_API_URL}/notification-channels/${channelId}/export`,
    {
      method: "GET",
    },
  );

  return res.text();
};

export const NotificationChannelsApi = {
  getAll,
  create,
  update,
  remove,
  exportYaml,
};
