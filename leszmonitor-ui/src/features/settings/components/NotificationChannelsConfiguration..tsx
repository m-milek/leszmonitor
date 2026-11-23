import { useQuery } from "@tanstack/react-query";
import {
  Card,
  CardAction,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import { NewNotificationChannelDialog } from "@/features/notifications/components/NewNotificationChannelDialog";
import { NotificationChannelsTable } from "@/features/notifications/components/NotificationChannelsTable";
import { NotificationChannelsApi } from "@/features/notifications/notification-channels-api";
import { QUERY_KEYS } from "@/lib/consts";

export const NotificationChannelsConfiguration = () => {
  const { data: channels } = useQuery({
    queryKey: [QUERY_KEYS.NOTIFICATION_CHANNELS],
    queryFn: () => NotificationChannelsApi.getAll(),
  });

  return (
    <Card>
      <CardHeader>
        <CardTitle>Notification Channels</CardTitle>
        <CardDescription>
          Where alerts are sent when a monitor changes status
        </CardDescription>
        <CardAction>
          <NewNotificationChannelDialog />
        </CardAction>
      </CardHeader>
      <CardContent className="px-0">
        {channels && <NotificationChannelsTable channels={channels} />}
      </CardContent>
    </Card>
  );
};
