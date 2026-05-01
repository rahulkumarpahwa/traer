import { useStore } from "../store/useStore";
import { Link2 } from "lucide-react";

export default function Sidebar() {
  const sidebarOpen = useStore((state) => state.sidebarOpen);
  const recentLinks = useStore((state) => state.recentLinks);

  return (
    <aside
      className={`${!sidebarOpen ? "hidden" : "block"} w-80 shrink-0 border-r border-emerald-500/10 bg-white/70 p-6 backdrop-blur-xl dark:border-emerald-400/10 dark:bg-[#04110b]/70`}
    >
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
                <Link2 className="h-4 w-4" strokeWidth={1.8} />
              </span>
              <span className="text-sm text-slate-700 dark:text-emerald-100/80 break-all">
                {link}
              </span>
            </div>
          ))
        )}
      </div>
    </aside>
  );
}
