import { useState } from "react";
import ActionCard from "../components/ActionCard";
import Header from "../components/Header";
import JobPanel from "../components/JobPanel";
import NotificationTray from "../components/NotificationTray";
import Sidebar from "../components/Sidebar";
import Terminal from "../components/Terminal";
import { useStore } from "../store/useStore";

const cards = [
  {
    kind: "transcribe",
    title: "Transcribe to notes",
    description: "Generate structured notes or raw text for meetings, podcasts, and lectures.",
    accent: "blue",
    cta: "Generate notes",
  },
  {
    kind: "audio",
    title: "Convert audio",
    description: "Package the source into delivery-ready audio formats with sensible presets.",
    accent: "green",
    cta: "Convert audio",
  },
  {
    kind: "video",
    title: "Export video",
    description: "Create clean video exports for archive, upload, or lightweight review.",
    accent: "amber",
    cta: "Render video",
  },
  {
    kind: "cloud",
    title: "Send to cloud",
    description: "Push artifacts to your storage endpoint as a neat audio, video, notes, or bundle package.",
    accent: "rose",
    cta: "Sync output",
  },
];

export default function Dashboard() {
  const [url, setUrl] = useState("");
  const hydrated = useStore((state) => state.hydrated);
  const capabilities = useStore((state) => state.capabilities);

  return (
    <div className="app-shell">
      <Header />
      <NotificationTray />

      <div className="workspace">
        <Sidebar />

        <main className="content">
          <section className="hero panel">
            <div>
              <p className="eyebrow">Workflow launcher</p>
              <h2 style={{ marginBottom: "0.5rem" }}>Move from raw link to usable assets</h2>
              <p className="hero-copy">
                traer keeps the interface minimal while still surfacing the controls and runtime visibility
                you need for real work.
              </p>
            </div>
            <div className="hero-status">
              <span className="metric">{hydrated ? "Ready" : "Loading"}</span>
              <span className="muted">
                {capabilities.backendBound ? "Connected to Wails backend" : "Preview with mock bridge"}
              </span>
            </div>
          </section>

          <section className="panel">
            <div className="panel-header">
              <div>
                <p className="eyebrow">Source</p>
                <h2>Media link</h2>
              </div>
            </div>

            <label className="input-label" htmlFor="source-url">
              Paste a YouTube, Loom, Vimeo, or podcast URL
            </label>
            <input
              id="source-url"
              className="input-field"
              placeholder="https://"
              value={url}
              onChange={(event) => setUrl(event.target.value)}
            />
          </section>

          <section className="card-grid">
            {cards.map((card) => (
              <ActionCard key={card.kind} url={url} disabled={!hydrated} {...card} />
            ))}
          </section>

          <JobPanel />
        </main>
      </div>

      <Terminal />
    </div>
  );
}
