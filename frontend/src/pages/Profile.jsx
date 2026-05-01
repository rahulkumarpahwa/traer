import { useEffect } from "react";
import Header from "../components/Header";
import NotificationTray from "../components/NotificationTray";
import Sidebar from "../components/Sidebar";
import Terminal from "../components/Terminal";
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

        <main className="flex-1 p-8 overflow-auto">
          <section className="bg-white dark:bg-black border-2 border-emerald-500 dark:border-emerald-700 rounded-xl p-8 max-w-2xl">
            <p className="text-xs uppercase tracking-widest text-emerald-500/70 font-semibold mb-2">
              Operator profile
            </p>
            <h2 className="text-2xl font-bold text-emerald-600 dark:text-emerald-300 mb-2">
              {profile?.name ?? "Loading profile..."}
            </h2>
            <p className="text-sm text-emerald-500/70 dark:text-emerald-400/70 mb-8">
              {profile?.bio}
            </p>

            <div className="grid grid-cols-2 gap-4">
              <article className="bg-emerald-50 dark:bg-emerald-950 border border-emerald-200 dark:border-emerald-800 rounded-md p-4 hover:border-emerald-400 dark:hover:border-emerald-600 transition-colors">
                <span className="text-xs uppercase tracking-widest text-emerald-500/70 font-semibold">
                  Email
                </span>
                <strong className="block text-emerald-700 dark:text-emerald-300 mt-2">
                  {user?.email ?? profile?.email ?? "operator@traer.local"}
                </strong>
              </article>
              <article className="bg-emerald-50 dark:bg-emerald-950 border border-emerald-200 dark:border-emerald-800 rounded-md p-4 hover:border-emerald-400 dark:hover:border-emerald-600 transition-colors">
                <span className="text-xs uppercase tracking-widest text-emerald-500/70 font-semibold">
                  Plan
                </span>
                <strong className="block text-emerald-700 dark:text-emerald-300 mt-2">
                  {profile?.plan ?? "Studio"}
                </strong>
              </article>
              <article className="bg-emerald-50 dark:bg-emerald-950 border border-emerald-200 dark:border-emerald-800 rounded-md p-4 hover:border-emerald-400 dark:hover:border-emerald-600 transition-colors">
                <span className="text-xs uppercase tracking-widest text-emerald-500/70 font-semibold">
                  Role
                </span>
                <strong className="block text-emerald-700 dark:text-emerald-300 mt-2">
                  {user?.role ?? "Admin"}
                </strong>
              </article>
              <article className="bg-emerald-50 dark:bg-emerald-950 border border-emerald-200 dark:border-emerald-800 rounded-md p-4 hover:border-emerald-400 dark:hover:border-emerald-600 transition-colors">
                <span className="text-xs uppercase tracking-widest text-emerald-500/70 font-semibold">
                  Last active
                </span>
                <strong className="block text-emerald-700 dark:text-emerald-300 mt-2">
                  {profile?.lastActive
                    ? new Date(profile.lastActive).toLocaleString()
                    : "Updating..."}
                </strong>
              </article>
            </div>
          </section>
        </main>
      </div>

      <Terminal />
    </div>
  );
}
