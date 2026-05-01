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
            ? "border-emerald-500/15 bg-white/90 dark:border-emerald-400/10 dark:bg-[#07150f]/95"
            : "border-amber-500/20 bg-white/90 dark:border-amber-400/15 dark:bg-[#1a1206]/95";
        const textColor =
          notification.tone === "success"
            ? "text-slate-900 dark:text-emerald-50"
            : "text-slate-900 dark:text-amber-100";

        return (
          <section
            key={notification.id}
            className={`${bgColor} animate-slideIn flex items-start justify-between gap-3 rounded-2xl border p-4 shadow-[0_18px_50px_-30px_rgba(2,6,23,0.45)] backdrop-blur-xl`}
          >
            <div>
              <strong className={textColor}>{notification.title}</strong>
              <p className={`mt-1 text-xs ${notification.tone === "success" ? "text-slate-600 dark:text-emerald-100/60" : "text-slate-600 dark:text-amber-100/60"}`}>
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
