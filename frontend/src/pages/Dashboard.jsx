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

  return (
    <div className="app-shell">
      <Header />
      <NotificationTray />

      <div className="workspace">
        <Sidebar />

        <main className="relative flex-1 overflow-hidden">
          <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.15),transparent_28%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.12),transparent_26%),linear-gradient(180deg,rgba(255,255,255,0.8),rgba(240,253,250,0.55))] dark:bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.16),transparent_24%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.14),transparent_24%),linear-gradient(180deg,rgba(4,17,11,1),rgba(6,24,17,0.95))]" />
          <section className="relative mx-auto flex w-full max-w-7xl flex-col gap-10 px-6 py-8 md:px-8 xl:px-10">
            <div className="grid gap-6 xl:grid-cols-[1.2fr_0.8fr]">
              <div className="rounded-[32px] border border-emerald-500/12 bg-white/82 p-8 shadow-[0_30px_90px_-45px_rgba(5,46,22,0.45)] backdrop-blur-sm dark:border-emerald-400/10 dark:bg-[#07150f]/88 dark:shadow-[0_30px_90px_-45px_rgba(16,185,129,0.28)]">
                <p className="mb-3 text-[11px] font-semibold uppercase tracking-[0.28em] text-emerald-700/55 dark:text-emerald-200/55">
                  Media intake
                </p>
                <h2 className="max-w-2xl text-4xl font-semibold tracking-tight text-slate-950 dark:text-emerald-50">
                  Turn any link into notes, exports, and synced deliverables.
                </h2>
                <p className="mt-4 max-w-2xl text-base leading-7 text-slate-600 dark:text-emerald-100/70">
                  Drop in a recording once, then route it through transcription, conversion, rendering, or cloud delivery with consistent presets.
                </p>

                <div className="mt-8 rounded-[28px] border border-emerald-500/12 bg-emerald-500/[0.05] p-4 dark:border-emerald-400/10 dark:bg-white/[0.03]">
                  <label
                    htmlFor="source-url"
                    className="mb-3 block text-[11px] font-semibold uppercase tracking-[0.22em] text-emerald-700/60 dark:text-emerald-200/60"
                  >
                    Source URL
                  </label>
                  <input
                    id="source-url"
                    className="w-full rounded-2xl border border-emerald-500/12 bg-white px-5 py-4 text-slate-900 outline-none transition placeholder:text-slate-400 focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50 dark:placeholder:text-emerald-100/35"
                    placeholder="Paste a YouTube, Loom, Vimeo, or podcast URL"
                    value={url}
                    onChange={(event) => setUrl(event.target.value)}
                  />
                  <p className="mt-3 text-sm text-slate-500 dark:text-emerald-100/55">
                    Best results come from public links or sources reachable by your configured backend.
                  </p>
                </div>
              </div>

              <div className="rounded-[32px] border border-emerald-500/12 bg-slate-950 p-8 text-white shadow-[0_30px_90px_-45px_rgba(2,6,23,0.65)] dark:border-emerald-400/10 dark:bg-[#0d2017]">
                <p className="mb-3 text-[11px] font-semibold uppercase tracking-[0.28em] text-emerald-200/60">
                  Workflow
                </p>
                <h3 className="text-2xl font-semibold text-white dark:text-emerald-50">
                  A cleaner handoff from raw media to finished assets.
                </h3>
                <div className="mt-6 space-y-4">
                  {[
                    "Paste the link once and reuse it across every workflow card.",
                    "Choose an output preset tailored for notes, audio, video, or bundles.",
                    "Track progress in the active jobs queue and runtime terminal.",
                  ].map((item) => (
                    <div
                      key={item}
                      className="rounded-2xl border border-white/10 bg-white/5 p-4 text-sm leading-6 text-emerald-50/80"
                    >
                      {item}
                    </div>
                  ))}
                </div>
              </div>
            </div>

            <div className="grid grid-cols-1 gap-6 md:grid-cols-2">
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
