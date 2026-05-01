import { useStore } from "../store/useStore";

export default function JobPanel() {
  const activeJobs = useStore((state) => state.activeJobs);

  return (
    <section className="panel">
      <div className="panel-header">
        <div>
          <p className="eyebrow">Queue</p>
          <h2>Active jobs</h2>
        </div>
      </div>

      {activeJobs.length === 0 ? (
        <p className="muted">No jobs yet. Start a transcript, conversion, or cloud sync.</p>
      ) : (
        <div className="job-list">
          {activeJobs.map((job) => (
            <article key={job.id} className="job-card">
              <div className="job-topline">
                <div>
                  <h3>{job.title}</h3>
                  <p className="muted">{job.message}</p>
                </div>
                <span className={job.status === "done" ? "badge success" : "badge"}>{job.status}</span>
              </div>

              <div className="progress-track" aria-hidden="true">
                <span style={{ width: `${job.progress}%` }} />
              </div>

              <div className="job-footer">
                <span>{job.option}</span>
                <span>{job.progress}%</span>
              </div>
            </article>
          ))}
        </div>
      )}
    </section>
  );
}
