import { useStore } from "../store/useStore";

export default function NotificationTray() {
  const notifications = useStore((state) => state.notifications);
  const dismissNotification = useStore((state) => state.dismissNotification);

  if (notifications.length === 0) {
    return null;
  }

  return (
    <div className="toast-stack" aria-live="polite">
      {notifications.map((notification) => (
        <section key={notification.id} className={`toast ${notification.tone}`}>
          <div>
            <strong>{notification.title}</strong>
            <p>{notification.message}</p>
          </div>
          <button className="toast-close" onClick={() => dismissNotification(notification.id)}>
            Close
          </button>
        </section>
      ))}
    </div>
  );
}
