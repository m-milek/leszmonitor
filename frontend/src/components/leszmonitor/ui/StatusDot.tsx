export const StatusDot = ({
  status,
}: {
  status: "success" | "failure" | "pending";
}) => {
  const color = {
    success: "bg-green-500",
    failure: "bg-red-700",
    pending: "bg-yellow-500",
  }[status];

  return <span className={`inline-block w-3 h-3 rounded-full ${color}`} />;
};
