import { useStore } from "../store/useStore";
import { CheckCircle2, Clock3, Sparkles } from "lucide-react";
import Header from "./Header";
import NotificationTray from "./NotificationTray";
import Sidebar from "./Sidebar";
import StatusPill from "./StatusPill";

function JobPanelContent() {
  const activeJobs = useStore((state) => state.activeJobs);

  return (
    <section className="mx-auto max-w-6xl rounded-[32px] border border-emerald-500/15 bg-white/85 p-8 shadow-[0_24px_80px_-40px_rgba(5,46,22,0.45)] backdrop-blur-sm dark:border-emerald-400/10 dark:bg-[#07150f]/90 dark:shadow-[0_24px_80px_-40px_rgba(16,185,129,0.3)]">
      <div className="mb-8 flex flex-col gap-6 lg:flex-row lg:items-end lg:justify-between">
        <div>
          <p className="mb-2 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
            Queue
          </p>
          <h2 className="text-3xl font-semibold text-slate-950 dark:text-emerald-50">
            Active jobs
          </h2>
          <p className="mt-3 max-w-2xl text-sm leading-6 text-slate-600 dark:text-emerald-100/65">
            Watch every workflow move from queue to completion, with progress, preset, and live status all in one view.
          </p>
        </div>
        <div className="grid grid-cols-1 gap-3 sm:grid-cols-3">
          <div className="rounded-2xl border border-emerald-500/12 bg-emerald-500/[0.05] p-4 dark:border-emerald-400/10 dark:bg-white/[0.03]">
            <Sparkles className="mb-3 h-5 w-5 text-emerald-600 dark:text-emerald-200" strokeWidth={1.8} />
            <p className="text-xs uppercase tracking-[0.2em] text-emerald-700/55 dark:text-emerald-200/55">Total</p>
            <strong className="mt-2 block text-2xl text-slate-950 dark:text-emerald-50">{activeJobs.length}</strong>
          </div>
          <div className="rounded-2xl border border-emerald-500/12 bg-emerald-500/[0.05] p-4 dark:border-emerald-400/10 dark:bg-white/[0.03]">
            <Clock3 className="mb-3 h-5 w-5 text-amber-600 dark:text-amber-200" strokeWidth={1.8} />
            <p className="text-xs uppercase tracking-[0.2em] text-emerald-700/55 dark:text-emerald-200/55">Running</p>
            <strong className="mt-2 block text-2xl text-slate-950 dark:text-emerald-50">
              {activeJobs.filter((job) => job.status !== "done").length}
            </strong>
          </div>
          <div className="rounded-2xl border border-emerald-500/12 bg-emerald-500/[0.05] p-4 dark:border-emerald-400/10 dark:bg-white/[0.03]">
            <CheckCircle2 className="mb-3 h-5 w-5 text-emerald-600 dark:text-emerald-200" strokeWidth={1.8} />
            <p className="text-xs uppercase tracking-[0.2em] text-emerald-700/55 dark:text-emerald-200/55">Completed</p>
            <strong className="mt-2 block text-2xl text-slate-950 dark:text-emerald-50">
              {activeJobs.filter((job) => job.status === "done").length}
            </strong>
          </div>
        </div>
      </div>

      {activeJobs.length === 0 ? (
        <p className="rounded-2xl border border-dashed border-emerald-500/15 bg-white/50 p-5 text-sm leading-6 text-slate-500 dark:border-emerald-400/10 dark:bg-white/[0.02] dark:text-emerald-100/55">
          No jobs yet. Start a transcript, conversion, or cloud sync.
        </p>
      ) : (
        <div className="space-y-3">
          {activeJobs.map((job) => (
            <article
              key={job.id}
              className="rounded-3xl border border-emerald-500/12 bg-emerald-500/[0.05] p-5 transition hover:border-emerald-500/20 dark:border-emerald-400/10 dark:bg-white/[0.03]"
            >
              <div className="mb-3 flex items-center justify-between gap-3">
                <div className="flex flex-1 items-start gap-4">
                  <span className="flex h-11 w-11 shrink-0 items-center justify-center rounded-2xl border border-emerald-500/12 bg-white/75 text-emerald-700 dark:border-emerald-400/10 dark:bg-white/[0.04] dark:text-emerald-200">
                    {job.status === "done" ? (
                      <CheckCircle2 className="h-5 w-5" strokeWidth={1.8} />
                    ) : (
                      <Clock3 className="h-5 w-5" strokeWidth={1.8} />
                    )}
                  </span>
                  <div className="flex-1">
                    <h3 className="font-semibold text-slate-900 dark:text-emerald-50">
                      {job.title}
                    </h3>
                    <p className="text-xs text-slate-500 dark:text-emerald-100/55">
                      {job.message}
                    </p>
                  </div>
                </div>
                <StatusPill tone={job.status === "done" ? "strong" : "warm"}>
                  {job.status}
                </StatusPill>
              </div>

              <div
                className="mb-2 h-2 w-full overflow-hidden rounded-full bg-emerald-500/10 dark:bg-emerald-400/10"
                aria-hidden="true"
              >
                <div
                  className="h-full rounded-full bg-gradient-to-r from-emerald-500 via-teal-500 to-cyan-500"
                  style={{
                    width: `${job.progress}%`,
                    transition: "width 300ms ease",
                  }}
                />
              </div>

              <div className="flex items-center justify-between text-xs">
                <span className="text-slate-600 dark:text-emerald-100/70">
                  {job.option}
                </span>
                <span className="font-semibold text-slate-900 dark:text-emerald-50">
                  {job.progress}%
                </span>
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}

export default function JobPanel() {
  return (
    <div className="app-shell">
      <Header />
      <NotificationTray />

      <div className="workspace">
        <Sidebar />

        <main className="relative flex-1 overflow-hidden">
          <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.15),transparent_28%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.12),transparent_26%),linear-gradient(180deg,rgba(255,255,255,0.8),rgba(240,253,250,0.55))] dark:bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.16),transparent_24%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.14),transparent_24%),linear-gradient(180deg,rgba(4,17,11,1),rgba(6,24,17,0.95))]" />
          <section className="relative mx-auto flex w-full max-w-7xl flex-col gap-10 px-6 py-8 md:px-8 xl:px-10">
            <JobPanelContent />
          </section>
        </main>
      </div>
    </div>
  );
}
