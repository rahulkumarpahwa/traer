import { useStore } from "../store/useStore";
import StatusPill from "./StatusPill";

export default function JobPanel() {
  const activeJobs = useStore((state) => state.activeJobs);

  return (
    <section className="mx-auto max-w-5xl rounded-[32px] border border-emerald-500/15 bg-white/85 p-8 shadow-[0_24px_80px_-40px_rgba(5,46,22,0.45)] backdrop-blur-sm dark:border-emerald-400/10 dark:bg-[#07150f]/90 dark:shadow-[0_24px_80px_-40px_rgba(16,185,129,0.3)]">
      <div className="mb-6">
        <p className="mb-2 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
          Queue
        </p>
        <h2 className="text-2xl font-semibold text-slate-950 dark:text-emerald-50">
          Active jobs
        </h2>
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
              <div className="flex items-center justify-between gap-3 mb-3">
                <div className="flex-1">
                  <h3 className="font-semibold text-slate-900 dark:text-emerald-50">
                    {job.title}
                  </h3>
                  <p className="text-xs text-slate-500 dark:text-emerald-100/55">
                    {job.message}
                  </p>
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
