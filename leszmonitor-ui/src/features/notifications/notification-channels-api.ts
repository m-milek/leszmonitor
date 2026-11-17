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

export const NotificationChannelsApi = {
  create,
};
