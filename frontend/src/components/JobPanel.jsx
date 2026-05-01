import { useStore } from "../store/useStore";

export default function JobPanel() {
  const activeJobs = useStore((state) => state.activeJobs);

  return (

    
    <section className="bg-white dark:bg-black border-2 border-emerald-500 dark:border-emerald-700 rounded-xl p-8">
      <div className="mb-6">
        <p className="text-xs uppercase tracking-widest text-emerald-500/70 font-semibold mb-2">
          Queue
        </p>
        <h2 className="text-2xl font-bold text-emerald-600 dark:text-emerald-300">
          Active jobs
        </h2>
      </div>

      {activeJobs.length === 0 ? (
        <p className="text-sm text-emerald-500/70">
          No jobs yet. Start a transcript, conversion, or cloud sync.
        </p>
      ) : (
        <div className="space-y-3">
          {activeJobs.map((job) => (
            <article
              key={job.id}
              className="bg-emerald-50 dark:bg-emerald-950 border border-emerald-200 dark:border-emerald-800 rounded-md p-4 hover:border-emerald-400 dark:hover:border-emerald-600 transition-colors"
            >
              <div className="flex items-center justify-between gap-3 mb-3">
                <div className="flex-1">
                  <h3 className="font-semibold text-emerald-700 dark:text-emerald-300">
                    {job.title}
                  </h3>
                  <p className="text-xs text-emerald-500/70 dark:text-emerald-400/70">
                    {job.message}
                  </p>
                </div>
                <span
                  className={`text-xs font-medium px-2 py-1 rounded-full ${job.status === "done" ? "bg-emerald-200 dark:bg-emerald-800 text-emerald-700 dark:text-emerald-300" : "bg-amber-100 dark:bg-amber-900 text-amber-700 dark:text-amber-300"}`}
                >
                  {job.status}
                </span>
              </div>

              <div
                className="w-full h-2 bg-emerald-100 dark:bg-emerald-900 rounded-full overflow-hidden mb-2"
                aria-hidden="true"
              >
                <div
                  className="h-full bg-gradient-to-r from-emerald-500 to-emerald-600"
                  style={{
                    width: `${job.progress}%`,
                    transition: "width 300ms ease",
                  }}
                />
              </div>

              <div className="flex items-center justify-between text-xs">
                <span className="text-emerald-600 dark:text-emerald-400">
                  {job.option}
                </span>
                <span className="font-semibold text-emerald-700 dark:text-emerald-300">
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
