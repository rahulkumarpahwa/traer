import { useStore } from "../store/useStore";

const NotesIcon = () => (
  <svg viewBox="0 0 24 24" style={{width: "24px", height: "24px"}}><path d="M14 2H6a2 2 0 0 0-2 2v16a2 2 0 0 0 2 2h12a2 2 0 0 0 2-2V8z"/><polyline points="14 2 14 8 20 8"/><line x1="12" y1="13" x2="16" y2="13"/><line x1="12" y1="17" x2="16" y2="17"/><line x1="12" y1="9" x2="16" y2="9"/></svg>
);

const AudioIcon = () => (
  <svg viewBox="0 0 24 24" style={{width: "24px", height: "24px"}}><polygon points="11 5 6 9 2 9 2 15 6 15 11 19 11 5"/><path d="M15.54 5.54a7 7 0 0 1 0 9.9M19.07 4a11 11 0 0 1 0 15.6"/></svg>
);

const VideoIcon = () => (
  <svg viewBox="0 0 24 24" style={{width: "24px", height: "24px"}}><polygon points="23 7 16 12 23 17 23 7"/><rect x="1" y="5" width="15" height="14" rx="2" ry="2"/></svg>
);

const CloudIcon = () => (
  <svg viewBox="0 0 24 24" style={{width: "24px", height: "24px"}}><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"/></svg>
);

const iconMap = {
  transcribe: NotesIcon,
  audio: AudioIcon,
  video: VideoIcon,
  cloud: CloudIcon,
};

export default function ActionCard({ kind, title, description, accent, cta, disabled, url }) {
  const actionOptions = useStore((state) => state.actionOptions);
  const selectedOption = useStore((state) => state.selectedOptions[kind]);
  const setOption = useStore((state) => state.setOption);
  const runJob = useStore((state) => state.runJob);

  const IconComponent = iconMap[kind];

  return (
    <article className="action-card" data-accent={accent}>
      <div className="card-header">
        <div>
          <p className="eyebrow">{kind}</p>
          <h3 style={{ display: "flex", alignItems: "center", gap: "0.8rem" }}>
            {IconComponent && <IconComponent />}
            {title}
          </h3>
        </div>
        <span className={`pill ${accent}`}>{selectedOption}</span>
      </div>
      <p className="card-copy">{description}</p>
      <select
        className="select-field"
        value={selectedOption}
        onChange={(event) => setOption(kind, event.target.value)}
      >
        {actionOptions[kind].map((option) => (
          <option key={option} value={option}>
            {option}
          </option>
        ))}
      </select>
      <button className="primary-button" onClick={() => runJob({ kind, url })} disabled={disabled}>
        {cta}
      </button>
    </article>
  );
}
