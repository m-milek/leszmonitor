import z from "zod";
import type {
  NotificationChannel,
  NotificationChannelKind,
} from "@/features/notifications/types";

const httpUrl = (field: string) =>
  z.url({
    protocol: /^https?$/,
    error: `${field} must be an http or https URL`,
  });

const baseFields = {
  name: z.string().trim().min(1, "Name cannot be empty"),
  enabled: z.boolean(),
};

const webhookChannelSchema = z.object({
  ...baseFields,
  kind: z.literal("webhook"),
  config: z.object({
    url: httpUrl("URL"),
    headers: z.record(z.string(), z.string()),
  }),
});

const buildNtfyChannelSchema = (isEdit: boolean) =>
  z.object({
    ...baseFields,
    kind: z.literal("ntfy"),
    config: z
      .object({
        serverUrl: httpUrl("Server URL"),
        topic: z.string().trim().min(1, "Topic cannot be empty"),
        token: z.string(),
        username: z.string(),
        password: z.string(),
        priority: z
          .number()
          .int()
          .min(0, "Priority must be between 1 and 5")
          .max(5, "Priority must be between 1 and 5"),
      })
      .refine((c) => !(c.token && c.username), {
        path: ["token"],
        message: "Use either a token or username/password",
      })
      .refine((c) => !c.password || c.username, {
        path: ["username"],
        message: "Username is required when a password is set",
      })
      .refine((c) => isEdit || !c.username || c.password, {
        path: ["password"],
        message: "Password is required when a username is set",
      }),
  });

const buildGotifyChannelSchema = (isEdit: boolean) =>
  z.object({
    ...baseFields,
    kind: z.literal("gotify"),
    config: z.object({
      serverUrl: httpUrl("Server URL"),
      appToken: isEdit
        ? z.string()
        : z.string().min(1, "App token cannot be empty"),
      priority: z
        .number()
        .int()
        .min(0, "Priority must be between 0 and 10")
        .max(10, "Priority must be between 0 and 10"),
    }),
  });

export const buildNotificationChannelSchema = (isEdit: boolean) =>
  z.discriminatedUnion("kind", [
    webhookChannelSchema,
    buildNtfyChannelSchema(isEdit),
    buildGotifyChannelSchema(isEdit),
  ]);

export type NotificationChannelFormValues = z.infer<
  ReturnType<typeof buildNotificationChannelSchema>
>;

export const defaultChannelConfigs: {
  [K in NotificationChannelKind]: Extract<
    NotificationChannelFormValues,
    { kind: K }
  >["config"];
} = {
  webhook: { url: "", headers: {} },
  ntfy: {
    serverUrl: "",
    topic: "",
    token: "",
    username: "",
    password: "",
    priority: 0,
  },
  gotify: { serverUrl: "", appToken: "", priority: 0 },
};

export const newNotificationChannelDefaultValues: NotificationChannelFormValues =
  {
    name: "",
    enabled: true,
    kind: "webhook",
    config: defaultChannelConfigs.webhook,
  };

export const channelToFormValues = (
  channel: NotificationChannel,
): NotificationChannelFormValues => {
  const base = { name: channel.name, enabled: channel.enabled };
  switch (channel.kind) {
    case "webhook":
      return {
        ...base,
        kind: "webhook",
        config: { ...defaultChannelConfigs.webhook, ...channel.config },
      };
    case "ntfy":
      return {
        ...base,
        kind: "ntfy",
        config: { ...defaultChannelConfigs.ntfy, ...channel.config },
      };
    case "gotify":
      return {
        ...base,
        kind: "gotify",
        config: { ...defaultChannelConfigs.gotify, ...channel.config },
      };
  }
};
