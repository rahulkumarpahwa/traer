import { useEffect } from "react";
import Header from "../components/Header";
import NotificationTray from "../components/NotificationTray";
import Sidebar from "../components/Sidebar";
import { useStore } from "../store/useStore";

export default function Profile() {
  const profile = useStore((state) => state.profile);
  const user = useStore((state) => state.user);
  const refreshProfile = useStore((state) => state.refreshProfile);

  useEffect(() => {
    refreshProfile();
  }, [refreshProfile]);

  return (
    <div className="app-shell">
      <Header />
      <NotificationTray />

      <div className="workspace">
        <Sidebar />

        <main className="relative flex-1 overflow-hidden">
          <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.15),transparent_28%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.12),transparent_26%),linear-gradient(180deg,rgba(255,255,255,0.8),rgba(240,253,250,0.55))] dark:bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.16),transparent_24%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.14),transparent_24%),linear-gradient(180deg,rgba(4,17,11,1),rgba(6,24,17,0.95))]" />
          <section className="relative flex min-h-full items-center justify-center px-6 py-8 md:px-8 xl:px-10">
            <section className="w-full max-w-3xl rounded-[32px] border border-emerald-500/15 bg-white/85 p-8 shadow-[0_24px_80px_-40px_rgba(5,46,22,0.45)] backdrop-blur-sm dark:border-emerald-400/10 dark:bg-[#07150f]/90 dark:shadow-[0_24px_80px_-40px_rgba(16,185,129,0.3)]">
            <p className="mb-2 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
              Operator profile
            </p>
            <h2 className="mb-2 text-3xl font-semibold text-slate-950 dark:text-emerald-50">
              {profile?.name ?? "Loading profile..."}
            </h2>
            <p className="mb-8 text-sm leading-6 text-slate-600 dark:text-emerald-100/65">
              {profile?.bio}
            </p>

            <div className="grid grid-cols-2 gap-4">
              <article className="rounded-2xl border border-emerald-500/12 bg-emerald-500/[0.05] p-4 transition hover:border-emerald-500/20 dark:border-emerald-400/10 dark:bg-white/[0.03]">
                <span className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700/55 dark:text-emerald-200/55">
                  Email
                </span>
                <strong className="mt-2 block text-slate-900 dark:text-emerald-50">
                  {user?.email ?? profile?.email ?? "operator@traer.local"}
                </strong>
              </article>
              <article className="rounded-2xl border border-emerald-500/12 bg-emerald-500/[0.05] p-4 transition hover:border-emerald-500/20 dark:border-emerald-400/10 dark:bg-white/[0.03]">
                <span className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700/55 dark:text-emerald-200/55">
                  Plan
                </span>
                <strong className="mt-2 block text-slate-900 dark:text-emerald-50">
                  {profile?.plan ?? "Studio"}
                </strong>
              </article>
              <article className="rounded-2xl border border-emerald-500/12 bg-emerald-500/[0.05] p-4 transition hover:border-emerald-500/20 dark:border-emerald-400/10 dark:bg-white/[0.03]">
                <span className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700/55 dark:text-emerald-200/55">
                  Role
                </span>
                <strong className="mt-2 block text-slate-900 dark:text-emerald-50">
                  {user?.role ?? "Admin"}
                </strong>
              </article>
              <article className="rounded-2xl border border-emerald-500/12 bg-emerald-500/[0.05] p-4 transition hover:border-emerald-500/20 dark:border-emerald-400/10 dark:bg-white/[0.03]">
                <span className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700/55 dark:text-emerald-200/55">
                  Last active
                </span>
                <strong className="mt-2 block text-slate-900 dark:text-emerald-50">
                  {profile?.lastActive
                    ? new Date(profile.lastActive).toLocaleString()
                    : "Updating..."}
                </strong>
              </article>
            </div>
            </section>
          </section>
        </main>
      </div>
    </div>
  );
}
