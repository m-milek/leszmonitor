import {
  Card,
  CardAction,
  CardContent,
  CardHeader,
  CardTitle,
} from "@/components/ui/card.tsx";
import { NewNotificationChannelDialog } from "@/features/notifications/components/NewNotificationChannelDialog";

export const NotificationChannelsConfiguration = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Notification Channels</CardTitle>
        <CardAction>
          <NewNotificationChannelDialog />
        </CardAction>
      </CardHeader>
      <CardContent>content</CardContent>
    </Card>
  );
};
