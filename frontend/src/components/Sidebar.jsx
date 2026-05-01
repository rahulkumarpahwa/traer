import { Link, useLocation } from "react-router-dom";
import { useStore } from "../store/useStore";

const DashboardIcon = () => (
  <svg viewBox="0 0 24 24">
    <rect x="3" y="3" width="7" height="7" />
    <rect x="14" y="3" width="7" height="7" />
    <rect x="14" y="14" width="7" height="7" />
    <rect x="3" y="14" width="7" height="7" />
  </svg>
);

const SettingsIcon = () => (
  <svg viewBox="0 0 24 24">
    <circle cx="12" cy="12" r="3" />
    <path d="M12 1v6m0 6v6M4.22 4.22l4.24 4.24m4.24 4.24l4.24 4.24M1 12h6m6 0h6m-15.78 7.78l4.24-4.24m4.24-4.24l4.24-4.24" />
  </svg>
);

const UserIcon = () => (
  <svg viewBox="0 0 24 24">
    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
    <circle cx="12" cy="7" r="4" />
  </svg>
);

const LinkIcon = () => (
  <svg viewBox="0 0 24 24">
    <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
    <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
  </svg>
);

const LightbulbIcon = () => (
  <svg viewBox="0 0 24 24">
    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
  </svg>
);

export default function Sidebar() {
  const location = useLocation();
  const sidebarOpen = useStore((state) => state.sidebarOpen);
  const recentLinks = useStore((state) => state.recentLinks);
  const suggestions = useStore((state) => state.suggestions);
  return (
    <aside
      className={`${!sidebarOpen ? "hidden" : "block"} bg-white dark:bg-black border-r border-emerald-100 dark:border-emerald-700 p-6 w-72 h-screen overflow-auto`}
    >
      <div className="mb-6">
        <p className="text-xs uppercase tracking-wider text-emerald-600 dark:text-emerald-400 font-semibold mb-3">
          Recent sources
        </p>
        {recentLinks.length === 0 ? (
          <p className="text-sm text-emerald-500/70">
            Your latest media links will appear here.
          </p>
        ) : (
          recentLinks.map((link) => (
            <div
              key={link}
              className="flex items-start gap-3 p-2 rounded-md hover:bg-emerald-50 dark:hover:bg-emerald-900"
              title={link}
            >
              <span className="w-4 h-4 text-emerald-500">
                <LinkIcon />
              </span>
              <span className="text-sm text-emerald-600 dark:text-emerald-300 break-all">
                {link}
              </span>
            </div>
          ))
        )}
      </div>
    </aside>
  );
}
