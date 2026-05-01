import { useStore } from "../store/useStore";

export default function NotificationTray() {
  const notifications = useStore((state) => state.notifications);
  const dismissNotification = useStore((state) => state.dismissNotification);

  if (notifications.length === 0) {
    return null;
  }

  return (
    <div
      className="fixed top-24 right-4 w-80 space-y-3 z-30"
      aria-live="polite"
    >
      {notifications.map((notification) => {
        const bgColor =
          notification.tone === "success"
            ? "bg-emerald-50 dark:bg-emerald-950 border-emerald-200 dark:border-emerald-800"
            : "bg-amber-50 dark:bg-amber-950 border-amber-200 dark:border-amber-800";
        const textColor =
          notification.tone === "success"
            ? "text-emerald-700 dark:text-emerald-300"
            : "text-amber-700 dark:text-amber-300";

        return (
          <section
            key={notification.id}
            className={`${bgColor} border-2 rounded-md p-4 flex justify-between items-start gap-3 animate-slideIn`}
          >
            <div>
              <strong className={textColor}>{notification.title}</strong>
              <p className={`text-xs ${textColor}/70 mt-1`}>
                {notification.message}
              </p>
            </div>
            <button
              onClick={() => dismissNotification(notification.id)}
              className={`text-xs font-semibold ${textColor} hover:opacity-70 transition-opacity whitespace-nowrap`}
            >
              Close
            </button>
          </section>
        );
      })}
    </div>
  );
}
