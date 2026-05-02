import { Link, useLocation } from "react-router-dom";
import { useStore } from "../store/useStore";
import { useEffect, useState } from "react";
import {
  LayoutDashboard,
  Menu,
  Moon,
  PlaySquare,
  SunMedium,
  TerminalSquare,
  UserRound,
} from "lucide-react";
import StatusPill from "./StatusPill";

const tabs = [
  { label: "Dashboard", to: "/", icon: LayoutDashboard },
  { label: "Active Jobs", to: "/jobs", icon: PlaySquare },
  { label: "Terminal", to: "/terminal", icon: TerminalSquare },
];

export default function Header() {
  const location = useLocation();
  const [theme, setTheme] = useState(localStorage.getItem("theme") || "light");
  const toggleSidebar = useStore((state) => state.toggleSidebar);
  const activeJobs = useStore((state) => state.activeJobs);

  useEffect(() => {
    const root = document.documentElement;
    root.setAttribute("data-theme", theme);
    if (theme === "dark") root.classList.add("dark");
    else root.classList.remove("dark");
    localStorage.setItem("theme", theme);
  }, [theme]);

  const toggleTheme = () => {
    setTheme(theme === "light" ? "dark" : "light");
  };

  return (
    <header className="sticky top-0 z-20 border-b border-emerald-500/10 bg-white/80 backdrop-blur-xl dark:border-emerald-400/10 dark:bg-[#04110b]/80">
      <div className="px-4 py-4 md:px-6">
        <div className="flex flex-col gap-4 xl:flex-row xl:items-center xl:justify-between">
          <div className="flex items-center justify-between gap-3">
            <button
              onClick={toggleSidebar}
              aria-label="Toggle sidebar"
              title="Toggle sidebar"
              className="inline-flex h-11 w-11 items-center justify-center rounded-2xl border border-emerald-500/10 bg-white/80 text-slate-700 shadow-sm shadow-emerald-950/5 transition hover:border-emerald-500/20 hover:text-emerald-700 dark:border-emerald-400/10 dark:bg-white/[0.03] dark:text-emerald-100 dark:hover:text-emerald-50"
            >
              <Menu className="h-5 w-5" strokeWidth={1.8} />
            </button>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <nav className="flex flex-wrap items-center gap-2">
              {tabs.map(({ label, to, icon: TabIcon }) => (
                <Link
                  key={to}
                  to={to}
                  className={`${location.pathname === to ? "border-emerald-500/20 bg-slate-950 text-white dark:bg-emerald-400 dark:text-slate-950" : "border-emerald-500/10 bg-white/70 text-slate-700 dark:bg-white/[0.03] dark:text-emerald-100/75"} inline-flex items-center gap-2 rounded-2xl border px-4 py-2 text-sm font-semibold transition hover:border-emerald-500/20 hover:text-emerald-700 dark:hover:text-emerald-50`}
                >
                  <TabIcon className="h-4 w-4" strokeWidth={1.8} />
                  {label}
                  {to === "/jobs" && activeJobs.length > 0 ? (
                    <span className="tab-badge">{activeJobs.length}</span>
                  ) : null}
                </Link>
              ))}
            </nav>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <StatusPill tone="neutral">{activeJobs.length} active jobs</StatusPill>
            <Link
              to="/profile"
              aria-label="Open profile"
              title="Open profile"
              className={`${location.pathname === "/profile" ? "border-emerald-500/20 bg-slate-950 text-white dark:bg-emerald-400 dark:text-slate-950" : "border-emerald-500/10 bg-white/80 text-slate-700 dark:bg-white/[0.03] dark:text-emerald-100"} inline-flex h-11 w-11 items-center justify-center rounded-2xl border shadow-sm shadow-emerald-950/5 transition hover:border-emerald-500/20 hover:text-emerald-700 dark:hover:text-emerald-50`}
            >
              <UserRound className="h-5 w-5" strokeWidth={1.8} />
            </Link>
            <button
              onClick={toggleTheme}
              aria-label="Toggle theme"
              title="Toggle dark mode"
              className="inline-flex h-11 w-11 items-center justify-center rounded-2xl border border-emerald-500/10 bg-white/80 text-slate-700 shadow-sm shadow-emerald-950/5 transition hover:border-emerald-500/20 hover:text-emerald-700 dark:border-emerald-400/10 dark:bg-white/[0.03] dark:text-emerald-100 dark:hover:text-emerald-50"
            >
              {theme === "light" ? (
                <Moon className="h-5 w-5" strokeWidth={1.8} />
              ) : (
                <SunMedium className="h-5 w-5" strokeWidth={1.8} />
              )}
            </button>
          </div>
        </div>
      </div>
    </header>
  );
}
