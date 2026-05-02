import { useState } from "react";
import {
  ArrowUpRight,
  CloudUpload,
  FileAudio2,
  FileText,
  Search,
  Video,
} from "lucide-react";
import Header from "../components/Header";
import NotificationTray from "../components/NotificationTray";
import Sidebar from "../components/Sidebar";
import { useStore } from "../store/useStore";

const actions = [
  {
    kind: "transcribe",
    label: "Transcribe",
    icon: FileText,
  },
  {
    kind: "audio",
    label: "Audio",
    icon: FileAudio2,
  },
  {
    kind: "video",
    label: "Video",
    icon: Video,
  },
  {
    kind: "cloud",
    label: "Upload",
    icon: CloudUpload,
  },
];

export default function Dashboard() {
  const [url, setUrl] = useState("");
  const hydrated = useStore((state) => state.hydrated);
  const actionOptions = useStore((state) => state.actionOptions);
  const selectedOptions = useStore((state) => state.selectedOptions);
  const setOption = useStore((state) => state.setOption);
  const runJob = useStore((state) => state.runJob);

  return (
    <div className="app-shell">
      <Header />
      <NotificationTray />

      <div className="workspace">
        <Sidebar />

        <main className="relative flex-1 overflow-hidden">
          <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.14),transparent_28%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.10),transparent_24%),linear-gradient(180deg,rgba(255,255,255,0.75),rgba(240,253,250,0.5))] dark:bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.14),transparent_24%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.12),transparent_22%),linear-gradient(180deg,rgba(4,17,11,1),rgba(6,24,17,0.95))]" />
          <section className="relative flex min-h-full items-center justify-center px-6 py-8 md:px-8 xl:px-10">
            <div className="w-full max-w-5xl rounded-[36px] border border-emerald-500/12 bg-white/82 p-8 shadow-[0_30px_90px_-45px_rgba(5,46,22,0.45)] backdrop-blur-sm dark:border-emerald-400/10 dark:bg-[#07150f]/90 dark:shadow-[0_30px_90px_-45px_rgba(16,185,129,0.28)] md:p-10">

              <div className="mx-auto  max-w-3xl">
                <div className="flex items-center gap-4 rounded-[28px] border border-emerald-500/12 bg-white px-5 py-4 shadow-sm shadow-emerald-950/5 dark:border-emerald-400/10 dark:bg-[#0a1912]">
                  <Search
                    className="h-5 w-5 shrink-0 text-slate-400 dark:text-emerald-100/35"
                    strokeWidth={1.8}
                  />
                  <input
                    id="source-url"
                    className="w-full bg-transparent text-base text-slate-900 outline-none placeholder:text-slate-400 dark:text-emerald-50 dark:placeholder:text-emerald-100/35"
                    placeholder="Youtube URL"
                    value={url}
                    onChange={(event) => setUrl(event.target.value)}
                  />
                </div>
              </div>

              <div className="mx-auto mt-8 grid max-w-4xl gap-4 md:grid-cols-2 xl:grid-cols-4">
                {actions.map(({ kind, label, icon: ActionIcon }) => (
                  <button
                    key={kind}
                    className="group flex items-center justify-between gap-4 rounded-[28px] border border-emerald-500/12 bg-emerald-500/[0.05] px-5 py-5 text-left transition hover:-translate-y-1 hover:border-emerald-500/25 hover:bg-white dark:border-emerald-400/10 dark:bg-white/[0.03] dark:hover:bg-white/[0.05] disabled:cursor-not-allowed disabled:opacity-45"
                    onClick={() => runJob({ kind, url })}
                    disabled={!hydrated}
                  >
                    <div className="flex items-center justify-center gap-2  p-2">
                      <div className="mb-3 flex h-11 w-11 items-center justify-center rounded-2xl border border-emerald-500  text-emerald-700 shadow-inner shadow-white/50 dark:border-emerald-400 dark:bg-[#102118] dark:text-emerald-200 dark:shadow-none">
                        <ActionIcon className="h-3 w-3" strokeWidth={1.4} />
                      </div>
                      <h2 className="text-base font-semibold text-slate-900 dark:text-emerald-50">
                      {label}
                      </h2>
                    </div>
                    <ArrowUpRight
                      className="h-5 w-5 shrink-0 text-slate-400 transition group-hover:text-emerald-600 dark:text-emerald-100/35 dark:group-hover:text-emerald-200"
                      strokeWidth={1.8}
                    />
                  </button>
                ))}
              </div>

              <div className="mx-auto mt-10 max-w-4xl">
                <p className="mb-4 text-center text-[11px] font-semibold uppercase tracking-[0.24em] text-emerald-700/55 dark:text-emerald-200/55">
                  Choose output type
                </p>
                <div className="grid gap-4 md:grid-cols-2 xl:grid-cols-4">
                  {actions.map(({ kind, label }) => (
                    <div
                      key={`${kind}-option`}
                      className="rounded-[24px] border border-emerald-500/10 bg-white/70 p-4 dark:border-emerald-400/10 dark:bg-white/[0.03]"
                    >
                      <label
                        htmlFor={`${kind}-option`}
                        className="mb-2 block text-xs font-semibold uppercase tracking-[0.18em] text-emerald-700/60 dark:text-emerald-200/60"
                      >
                        {label}
                      </label>
                      <select
                        id={`${kind}-option`}
                        className="w-full rounded-2xl border border-emerald-500/10 bg-white px-3 py-3 text-sm font-medium text-slate-800 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-100"
                        value={selectedOptions[kind]}
                        onChange={(event) =>
                          setOption(kind, event.target.value)
                        }
                      >
                        {actionOptions[kind].map((option) => (
                          <option key={option} value={option}>
                            {option}
                          </option>
                        ))}
                      </select>
                    </div>
                  ))}
                </div>
              </div>
            </div>
          </section>
        </main>
      </div>
    </div>
  );
}
