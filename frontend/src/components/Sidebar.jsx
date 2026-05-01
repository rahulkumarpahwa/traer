import { Link, useLocation } from "react-router-dom";
import { useStore } from "../store/useStore";
import Icon from "./Icon";
import StatusPill from "./StatusPill";

const DashboardIcon = () => (
  <Icon className="h-4 w-4">
    <rect x="3" y="3" width="7" height="7" />
    <rect x="14" y="3" width="7" height="7" />
    <rect x="14" y="14" width="7" height="7" />
    <rect x="3" y="14" width="7" height="7" />
  </Icon>
);

const SettingsIcon = () => (
  <Icon className="h-4 w-4">
    <circle cx="12" cy="12" r="3" />
    <path d="M12 1v6m0 6v6M4.22 4.22l4.24 4.24m4.24 4.24l4.24 4.24M1 12h6m6 0h6m-15.78 7.78l4.24-4.24m4.24-4.24l4.24-4.24" />
  </Icon>
);

const UserIcon = () => (
  <Icon className="h-4 w-4">
    <path d="M20 21v-2a4 4 0 0 0-4-4H8a4 4 0 0 0-4 4v2" />
    <circle cx="12" cy="7" r="4" />
  </Icon>
);

const LinkIcon = () => (
  <Icon className="h-4 w-4">
    <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
    <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
  </Icon>
);

const LightbulbIcon = () => (
  <Icon className="h-4 w-4">
    <path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z" />
  </Icon>
);

export default function Sidebar() {
  const location = useLocation();
  const sidebarOpen = useStore((state) => state.sidebarOpen);
  const recentLinks = useStore((state) => state.recentLinks);
  const suggestions = useStore((state) => state.suggestions);

  const quickLinks = [
    { label: "Dashboard", to: "/", icon: DashboardIcon },
    { label: "Endpoints", to: "/settings", icon: SettingsIcon },
    { label: "Profile", to: "/profile", icon: UserIcon },
  ];

  return (
    <aside
      className={`${!sidebarOpen ? "hidden" : "block"} w-80 shrink-0 border-r border-emerald-500/10 bg-white/70 p-6 backdrop-blur-xl dark:border-emerald-400/10 dark:bg-[#04110b]/70`}
    >
      <div className="mb-8 rounded-[28px] border border-emerald-500/10 bg-gradient-to-br from-emerald-500/[0.08] via-white/80 to-cyan-500/[0.08] p-5 dark:border-emerald-400/10 dark:from-emerald-400/[0.08] dark:via-white/[0.03] dark:to-cyan-400/[0.04]">
        <p className="mb-2 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
          Workspace guide
        </p>
        <h2 className="mb-3 text-xl font-semibold text-slate-950 dark:text-emerald-50">
          Keep source links, jobs, and tips in one place.
        </h2>
        <p className="mb-4 text-sm leading-6 text-slate-600 dark:text-emerald-100/70">
          Start from a media link, choose a preset, and watch each task move through the queue.
        </p>
        <div className="flex flex-wrap gap-2">
          <StatusPill tone="strong">{recentLinks.length} recent links</StatusPill>
          <StatusPill tone="neutral">{suggestions.length} suggestions</StatusPill>
        </div>
      </div>

      <div className="mb-8">
        <p className="mb-3 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
          Quick access
        </p>
        <nav className="space-y-2">
          {quickLinks.map(({ label, to, icon: QuickIcon }) => (
            <Link
              key={to}
              to={to}
              className={`${location.pathname === to ? "border-emerald-500/20 bg-slate-950 text-white dark:bg-emerald-400 dark:text-slate-950" : "border-emerald-500/10 bg-white/65 text-slate-700 dark:bg-white/[0.03] dark:text-emerald-100/80"} flex items-center gap-3 rounded-2xl border px-4 py-3 text-sm font-medium transition hover:border-emerald-500/20 hover:text-emerald-700 dark:hover:text-emerald-50`}
            >
              <span className="flex h-9 w-9 items-center justify-center rounded-xl border border-current/10">
                <QuickIcon />
              </span>
              {label}
            </Link>
          ))}
        </nav>
      </div>

      <div className="mb-6">
        <p className="mb-3 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
          Recent sources
        </p>
        {recentLinks.length === 0 ? (
          <p className="rounded-2xl border border-dashed border-emerald-500/15 bg-white/50 p-4 text-sm leading-6 text-slate-500 dark:border-emerald-400/10 dark:bg-white/[0.02] dark:text-emerald-100/55">
            Your latest media links will appear here.
          </p>
        ) : (
          recentLinks.map((link) => (
            <div
              key={link}
              className="mb-2 flex items-start gap-3 rounded-2xl border border-emerald-500/10 bg-white/65 p-3 transition hover:border-emerald-500/20 hover:bg-white dark:border-emerald-400/10 dark:bg-white/[0.03] dark:hover:bg-white/[0.05]"
              title={link}
            >
              <span className="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-xl border border-emerald-500/10 text-emerald-600 dark:text-emerald-200">
                <LinkIcon />
              </span>
              <span className="text-sm text-slate-700 dark:text-emerald-100/80 break-all">
                {link}
              </span>
            </div>
          ))
        )}
      </div>

      <div>
        <p className="mb-3 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
          Suggestions
        </p>
        {suggestions.length === 0 ? (
          <div className="rounded-2xl border border-dashed border-emerald-500/15 bg-white/50 p-4 text-sm leading-6 text-slate-500 dark:border-emerald-400/10 dark:bg-white/[0.02] dark:text-emerald-100/55">
            Suggestions from the backend will show up here when available.
          </div>
        ) : (
          <div className="space-y-2">
            {suggestions.map((suggestion, index) => (
              <div
                key={`${suggestion}-${index}`}
                className="flex items-start gap-3 rounded-2xl border border-emerald-500/10 bg-white/65 p-3 dark:border-emerald-400/10 dark:bg-white/[0.03]"
              >
                <span className="mt-0.5 flex h-8 w-8 shrink-0 items-center justify-center rounded-xl border border-emerald-500/10 text-emerald-600 dark:text-emerald-200">
                  <LightbulbIcon />
                </span>
                <p className="text-sm leading-6 text-slate-700 dark:text-emerald-100/80">
                  {suggestion}
                </p>
              </div>
            ))}
          </div>
        )}
      </div>
    </aside>
  );
}
