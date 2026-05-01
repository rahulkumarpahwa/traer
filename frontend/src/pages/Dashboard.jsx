import { useState } from "react";
import ActionCard from "../components/ActionCard";
import Header from "../components/Header";
// JobPanel moved to tab/indicator area; remove inline panel from Dashboard
import NotificationTray from "../components/NotificationTray";
import Sidebar from "../components/Sidebar";
import Terminal from "../components/Terminal";
import { useStore } from "../store/useStore";

const cards = [
  {
    kind: "transcribe",
    title: "Transcribe to notes",
    description:
      "Generate structured notes or raw text for meetings, podcasts, and lectures.",
    accent: "blue",
    cta: "Generate notes",
  },
  {
    kind: "audio",
    title: "Convert audio",
    description:
      "Package the source into delivery-ready audio formats with sensible presets.",
    accent: "green",
    cta: "Convert audio",
  },
  {
    kind: "video",
    title: "Export video",
    description:
      "Create clean video exports for archive, upload, or lightweight review.",
    accent: "amber",
    cta: "Render video",
  },
  {
    kind: "cloud",
    title: "Send to cloud",
    description:
      "Push artifacts to your storage endpoint as a neat audio, video, notes, or bundle package.",
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

        <main className="flex-1 flex items-center justify-center p-16 md:p-8">
          <section className="w-full max-w-2xl flex flex-col gap-8">
            <label htmlFor="source-url" className="sr-only">
              Media URL
            </label>
            <input
              id="source-url"
              className="w-full px-6 py-4 rounded-xl border-2 border-emerald-500 text-emerald-600 bg-white dark:bg-black dark:border-emerald-700 focus:outline-none focus:ring-2 focus:ring-emerald-400 dark:focus:ring-emerald-600 transition-colors"
              placeholder="Paste a YouTube, Loom, Vimeo, or podcast URL"
              value={url}
              onChange={(event) => setUrl(event.target.value)}
            />

            <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
              {cards.map((card) => (
                <ActionCard
                  key={card.kind}
                  url={url}
                  disabled={!hydrated}
                  {...card}
                />
              ))}
            </div>
          </section>
        </main>
      </div>

      <Terminal />
    </div>
  );
}
