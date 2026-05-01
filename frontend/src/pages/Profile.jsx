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

        <main className="content">
          <section className="panel profile-panel">
            <p className="eyebrow">Operator profile</p>
            <h2>{profile?.name ?? "Loading profile..."}</h2>
            <p className="muted">{profile?.bio}</p>

            <div className="stats-grid">
              <article className="stat-card">
                <span className="stat-label">Email</span>
                <strong>
                  {user?.email ?? profile?.email ?? "operator@traer.local"}
                </strong>
              </article>
              <article className="stat-card">
                <span className="stat-label">Plan</span>
                <strong>{profile?.plan ?? "Studio"}</strong>
              </article>
              <article className="stat-card">
                <span className="stat-label">Role</span>
                <strong>{user?.role ?? "Admin"}</strong>
              </article>
              <article className="stat-card">
                <span className="stat-label">Last active</span>
                <strong>
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
