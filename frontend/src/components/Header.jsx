import { Link, useLocation } from "react-router-dom";
import { useStore } from "../store/useStore";
import { useState, useEffect } from "react";

const navItems = [
  { label: "Workspace", to: "/" },
  { label: "Settings", to: "/settings" },
  { label: "Profile", to: "/profile" },
];

const MenuIcon = () => (
  <svg viewBox="0 0 24 24">
    <line x1="3" y1="6" x2="21" y2="6" />
    <line x1="3" y1="12" x2="21" y2="12" />
    <line x1="3" y1="18" x2="21" y2="18" />
  </svg>
);

const SunIcon = () => (
  <svg viewBox="0 0 24 24">
    <circle cx="12" cy="12" r="5" />
    <line x1="12" y1="1" x2="12" y2="3" />
    <line x1="12" y1="21" x2="12" y2="23" />
    <line x1="4.22" y1="4.22" x2="5.64" y2="5.64" />
    <line x1="18.36" y1="18.36" x2="19.78" y2="19.78" />
    <line x1="1" y1="12" x2="3" y2="12" />
    <line x1="21" y1="12" x2="23" y2="12" />
    <line x1="4.22" y1="19.78" x2="5.64" y2="18.36" />
    <line x1="18.36" y1="5.64" x2="19.78" y2="4.22" />
  </svg>
);

const MoonIcon = () => (
  <svg viewBox="0 0 24 24">
    <path d="M21 12.79A9 9 0 1 1 11.21 3 7 7 0 0 0 21 12.79z" />
  </svg>
);

const TerminalIcon = () => (
  <svg viewBox="0 0 24 24">
    <polyline points="4 17 10 11 4 5" />
    <line x1="12" y1="19" x2="20" y2="19" />
  </svg>
);

const CheckIcon = () => (
  <svg viewBox="0 0 24 24">
    <polyline points="20 6 9 17 4 12" />
  </svg>
);

const CircleIcon = () => (
  <svg viewBox="0 0 24 24">
    <circle cx="12" cy="12" r="10" />
  </svg>
);

export default function Header() {
  const location = useLocation();
  const [theme, setTheme] = useState(localStorage.getItem("theme") || "light");
  const appName = useStore((state) => state.appName);
  const toggleSidebar = useStore((state) => state.toggleSidebar);
  const toggleTerminal = useStore((state) => state.toggleTerminal);
  const activeJobs = useStore((state) => state.activeJobs);
  const capabilities = useStore((state) => state.capabilities);

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
    <header className="bg-white dark:bg-black border-b border-emerald-100 dark:border-emerald-700 sticky top-0 z-20">
      <div className="flex flex-col gap-1 px-4 py-3">
        <div className="flex items-center justify-between">
          <div className="flex items-center">
            <button
              onClick={toggleSidebar}
              aria-label="Toggle sidebar"
              title="Toggle sidebar"
              className="p-2 rounded-md text-emerald-600 dark:text-emerald-400 hover:bg-emerald-50 dark:hover:bg-emerald-900"
            >
              <MenuIcon />
            </button>
          </div>

          <div className="pt-2">
            <nav className="flex items-center gap-2">
              <Link
                to="/"
                className={`${location.pathname === "/" ? "bg-emerald-100 dark:bg-emerald-900 text-emerald-700 dark:text-emerald-300" : "text-emerald-500/70"} px-3 py-1 rounded-md font-semibold flex items-center`}
              >
                Dashboard
                {activeJobs && activeJobs.length > 0 ? (
                  <span className="tab-badge">{activeJobs.length}</span>
                ) : null}
              </Link>

              <Link
                to="/settings"
                className={`${location.pathname === "/settings" ? "bg-emerald-100 dark:bg-emerald-900 text-emerald-700 dark:text-emerald-300" : "text-emerald-500/70"} px-3 py-1 rounded-md font-semibold`}
              >
                Endpoints
              </Link>

              <Link
                to="/profile"
                className={`${location.pathname === "/profile" ? "bg-emerald-100 dark:bg-emerald-900 text-emerald-700 dark:text-emerald-300" : "text-emerald-500/70"} px-3 py-1 rounded-md font-semibold`}
              >
                Profile
              </Link>
            </nav>
          </div>

          <div className="flex items-center gap-2">
            <button
              onClick={toggleTerminal}
              title="Toggle terminal"
              aria-label="Toggle terminal"
              className="p-2 rounded-md text-emerald-600 dark:text-emerald-400 hover:bg-emerald-50 dark:hover:bg-emerald-900"
            >
              <TerminalIcon />
            </button>

            <button
              onClick={toggleTheme}
              aria-label="Toggle theme"
              title="Toggle dark mode"
              className="p-2 rounded-md text-emerald-600 dark:text-emerald-400 hover:bg-emerald-50 dark:hover:bg-emerald-900"
            >
              {theme === "light" ? <MoonIcon /> : <SunIcon />}
            </button>
          </div>
        </div>
      </div>
    </header>
  );
}
