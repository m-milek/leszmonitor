import type { ResourceSource, Timestamps } from "@/lib/types";

export const notificationChannelKinds = ["webhook", "ntfy", "gotify"] as const;
export type NotificationChannelKind = (typeof notificationChannelKinds)[number];

export interface WebhookConfig {
  url: string;
  headers?: Record<string, string>;
}

export interface NtfyConfig {
  serverUrl: string;
  topic: string;
  token?: string;
  username?: string;
  password?: string;
  priority?: number;
}

export interface GotifyConfig {
  serverUrl: string;
  appToken?: string;
  priority: number;
}

export type NotificationChannelConfig =
  | { kind: "webhook"; config: WebhookConfig }
  | { kind: "ntfy"; config: NtfyConfig }
  | { kind: "gotify"; config: GotifyConfig };

export type NotificationChannel = Timestamps &
  NotificationChannelConfig & {
    id: string;
    name: string;
    enabled: boolean;
    source: ResourceSource;
  };
