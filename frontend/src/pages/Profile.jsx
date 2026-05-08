import { useEffect } from "react";
import { CalendarDays, Fingerprint, Mail, UserRound } from "lucide-react";
import Header from "../components/Header";
import NotificationTray from "../components/NotificationTray";
import Sidebar from "../components/Sidebar";
import { useStore } from "../store/useStore";

export default function Profile() {
  const profile = useStore((state) => state.profile);
  const user = useStore((state) => state.user);
  const refreshProfile = useStore((state) => state.refreshProfile);

  useEffect(() => {
    refreshProfile().catch(() => undefined);
  }, [refreshProfile]);

  const cards = [
    {
      label: "Username",
      value: user?.username || profile?.username || "Loading...",
      icon: UserRound,
    },
    {
      label: "Email",
      value: user?.email || profile?.email || "Loading...",
      icon: Mail,
    },
    {
      label: "User ID",
      value: user?.id ? `#${user.id}` : "Loading...",
      icon: Fingerprint,
    },
    {
      label: "Created",
      value: profile?.createdAt
        ? new Date(profile.createdAt).toLocaleString()
        : "Not available",
      icon: CalendarDays,
    },
  ];

  return (
    <div className="app-shell">
      <Header />
      <NotificationTray />

      <div className="workspace">
        <Sidebar />

        <main className="relative flex-1 overflow-hidden">
          <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.15),transparent_28%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.12),transparent_26%),linear-gradient(180deg,rgba(255,255,255,0.8),rgba(240,253,250,0.55))] dark:bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.16),transparent_24%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.14),transparent_24%),linear-gradient(180deg,rgba(4,17,11,1),rgba(6,24,17,0.95))]" />
          <section className="relative flex min-h-full items-center justify-center px-6 py-8 md:px-8 xl:px-10">
            <section className="w-full max-w-4xl rounded-[32px] border border-emerald-500/15 bg-white/85 p-8 shadow-[0_24px_80px_-40px_rgba(5,46,22,0.45)] backdrop-blur-sm dark:border-emerald-400/10 dark:bg-[#07150f]/90 dark:shadow-[0_24px_80px_-40px_rgba(16,185,129,0.3)]">
              <p className="mb-2 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
                Authenticated profile
              </p>
              <h2 className="mb-2 text-3xl font-semibold text-slate-950 dark:text-emerald-50">
                {profile?.name ?? user?.username ?? "Loading profile..."}
              </h2>
              <p className="mb-8 text-sm leading-6 text-slate-600 dark:text-emerald-100/65">
                {profile?.bio ?? "Your backend-authenticated account details are shown below."}
              </p>

              <div className="grid gap-4 md:grid-cols-2">
                {cards.map(({ label, value, icon: Icon }) => (
                  <article
                    key={label}
                    className="rounded-2xl border border-emerald-500/12 bg-emerald-500/[0.05] p-5 transition hover:border-emerald-500/20 dark:border-emerald-400/10 dark:bg-white/[0.03]"
                  >
                    <div className="mb-4 flex h-11 w-11 items-center justify-center rounded-2xl border border-emerald-500/12 bg-white text-emerald-700 dark:border-emerald-400/10 dark:bg-[#102118] dark:text-emerald-200">
                      <Icon className="h-5 w-5" strokeWidth={1.8} />
                    </div>
                    <span className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700/55 dark:text-emerald-200/55">
                      {label}
                    </span>
                    <strong className="mt-2 block break-all text-slate-900 dark:text-emerald-50">
                      {value}
                    </strong>
                  </article>
                ))}
              </div>

              <div className="mt-8 rounded-2xl border border-emerald-500/12 bg-white/70 p-5 dark:border-emerald-400/10 dark:bg-white/[0.03]">
                <p className="text-xs font-semibold uppercase tracking-[0.2em] text-emerald-700/55 dark:text-emerald-200/55">
                  Session details
                </p>
                <div className="mt-4 grid gap-4 md:grid-cols-2">
                  <div>
                    <span className="text-sm text-slate-500 dark:text-emerald-100/55">
                      Plan
                    </span>
                    <p className="mt-1 font-semibold text-slate-900 dark:text-emerald-50">
                      {profile?.plan ?? "Studio"}
                    </p>
                  </div>
                  <div>
                    <span className="text-sm text-slate-500 dark:text-emerald-100/55">
                      Updated
                    </span>
                    <p className="mt-1 font-semibold text-slate-900 dark:text-emerald-50">
                      {profile?.updatedAt
                        ? new Date(profile.updatedAt).toLocaleString()
                        : "Not available"}
                    </p>
                  </div>
                </div>
              </div>
            </section>
          </section>
        </main>
      </div>
    </div>
  );
}
