import { Link, useLocation } from "react-router-dom";
import { useStore } from "../store/useStore";
import { useState, useEffect } from "react";
import Icon from "./Icon";
import StatusPill from "./StatusPill";

const MenuIcon = () => (
  <Icon>
    <line x1="3" y1="6" x2="21" y2="6" />
    <line x1="3" y1="12" x2="21" y2="12" />
    <line x1="3" y1="18" x2="21" y2="18" />
  </Icon>
);

const SunIcon = () => (
  <Icon>
    <circle cx="12" cy="12" r="5" />
    <line x1="12" y1="1" x2="12" y2="3" />
    <line x1="12" y1="21" x2="12" y2="23" />
    <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
    <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
    <line x1="1" y1="12" x2="3" y2="12" />
    <line x1="21" y1="12" x2="23" y2="12" />
    <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
    <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
  </Icon>
);

const MoonIcon = () => (
  <Icon>
    <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
  </Icon>
);

const TerminalIcon = () => (
  <Icon>
    <polyline points="4 17 10 11 4 5" />
    <line x1="12" y1="19" x2="20" y2="19" />
  </Icon>
);

export default function Header() {
  const location = useLocation();
  const [theme, setTheme] = useState(localStorage.getItem("theme") || "light");
  const appName = useStore((state) => state.appName);
  const toggleSidebar = useStore((state) => state.toggleSidebar);
  const toggleTerminal = useStore((state) => state.toggleTerminal);
  const activeJobs = useStore((state) => state.activeJobs);
  const capabilities = useStore((state) => state.capabilities);
  const terminalOpen = useStore((state) => state.terminalOpen);

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
              <MenuIcon />
            </button>
            <div className="min-w-0">
              <p className="text-[11px] font-semibold uppercase tracking-[0.28em] text-emerald-700/55 dark:text-emerald-200/55">
                Capture workspace
              </p>
              <h1 className="truncate text-lg font-semibold text-slate-950 dark:text-emerald-50">
                {appName}
              </h1>
            </div>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <nav className="flex flex-wrap items-center gap-2">
              <Link
                to="/"
                className={`${location.pathname === "/" ? "border-emerald-500/20 bg-slate-950 text-white dark:bg-emerald-400 dark:text-slate-950" : "border-emerald-500/10 bg-white/70 text-slate-700 dark:bg-white/[0.03] dark:text-emerald-100/75"} inline-flex items-center rounded-2xl border px-4 py-2 text-sm font-semibold transition hover:border-emerald-500/20 hover:text-emerald-700 dark:hover:text-emerald-50`}
              >
                Dashboard
                {activeJobs && activeJobs.length > 0 ? (
                  <span className="tab-badge">{activeJobs.length}</span>
                ) : null}
              </Link>

              <Link
                to="/settings"
                className={`${location.pathname === "/settings" ? "border-emerald-500/20 bg-slate-950 text-white dark:bg-emerald-400 dark:text-slate-950" : "border-emerald-500/10 bg-white/70 text-slate-700 dark:bg-white/[0.03] dark:text-emerald-100/75"} inline-flex items-center rounded-2xl border px-4 py-2 text-sm font-semibold transition hover:border-emerald-500/20 hover:text-emerald-700 dark:hover:text-emerald-50`}
              >
                Endpoints
              </Link>

              <Link
                to="/jobs"
                className={`${location.pathname === "/jobs" ? "border-emerald-500/20 bg-slate-950 text-white dark:bg-emerald-400 dark:text-slate-950" : "border-emerald-500/10 bg-white/70 text-slate-700 dark:bg-white/[0.03] dark:text-emerald-100/75"} inline-flex items-center rounded-2xl border px-4 py-2 text-sm font-semibold transition hover:border-emerald-500/20 hover:text-emerald-700 dark:hover:text-emerald-50`}
              >
                Active Jobs
              </Link>

              <Link
                to="/profile"
                className={`${location.pathname === "/profile" ? "border-emerald-500/20 bg-slate-950 text-white dark:bg-emerald-400 dark:text-slate-950" : "border-emerald-500/10 bg-white/70 text-slate-700 dark:bg-white/[0.03] dark:text-emerald-100/75"} inline-flex items-center rounded-2xl border px-4 py-2 text-sm font-semibold transition hover:border-emerald-500/20 hover:text-emerald-700 dark:hover:text-emerald-50`}
              >
                Profile
              </Link>
            </nav>
          </div>

          <div className="flex flex-wrap items-center gap-2">
            <StatusPill tone={capabilities.backendBound ? "strong" : "warm"}>
              {capabilities.backendBound ? "Backend live" : "Preview mode"}
            </StatusPill>
            <StatusPill tone="neutral">
              {terminalOpen ? "Terminal open" : "Terminal hidden"}
            </StatusPill>
            <button
              onClick={toggleTerminal}
              title="Toggle terminal"
              aria-label="Toggle terminal"
              className="inline-flex h-11 w-11 items-center justify-center rounded-2xl border border-emerald-500/10 bg-white/80 text-slate-700 shadow-sm shadow-emerald-950/5 transition hover:border-emerald-500/20 hover:text-emerald-700 dark:border-emerald-400/10 dark:bg-white/[0.03] dark:text-emerald-100 dark:hover:text-emerald-50"
            >
              <TerminalIcon />
            </button>

            <button
              onClick={toggleTheme}
              aria-label="Toggle theme"
              title="Toggle dark mode"
              className="inline-flex h-11 w-11 items-center justify-center rounded-2xl border border-emerald-500/10 bg-white/80 text-slate-700 shadow-sm shadow-emerald-950/5 transition hover:border-emerald-500/20 hover:text-emerald-700 dark:border-emerald-400/10 dark:bg-white/[0.03] dark:text-emerald-100 dark:hover:text-emerald-50"
            >
              {theme === "light" ? <MoonIcon /> : <SunIcon />}
            </button>
          </div>
        </div>
      </div>
    </header>
  );
}
