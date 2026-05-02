import Header from "../components/Header";
import NotificationTray from "../components/NotificationTray";
import Sidebar from "../components/Sidebar";
import { useStore } from "../store/useStore";
import { useEffect, useState } from "react";

export default function Settings() {
  const storedSettings = useStore((state) => state.settings);
  const saveAppSettings = useStore((state) => state.saveAppSettings);
  const [form, setForm] = useState(storedSettings);
  const [saving, setSaving] = useState(false);

  useEffect(() => {
    setForm(storedSettings);
  }, [storedSettings]);

  const onChange = (key, value) => {
    setForm((state) => ({
      ...state,
      [key]: value,
    }));
  };

  const onSubmit = async (event) => {
    event.preventDefault();
    setSaving(true);
    try {
      await saveAppSettings(form);
    } finally {
      setSaving(false);
    }
  };

  return (
    <div className="app-shell">
      <Header />
      <NotificationTray />

      <div className="workspace">
        <Sidebar />

        <main className="relative flex-1 overflow-hidden">
          <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.15),transparent_28%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.12),transparent_26%),linear-gradient(180deg,rgba(255,255,255,0.8),rgba(240,253,250,0.55))] dark:bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.16),transparent_24%),radial-gradient(circle_at_top_right,rgba(6,182,212,0.14),transparent_24%),linear-gradient(180deg,rgba(4,17,11,1),rgba(6,24,17,0.95))]" />
          <section className="relative flex min-h-full items-center justify-center px-6 py-8 md:px-8 xl:px-10">
            <section className="w-full max-w-3xl rounded-[32px] border border-emerald-500/15 bg-white/85 p-8 shadow-[0_24px_80px_-40px_rgba(5,46,22,0.45)] backdrop-blur-sm dark:border-emerald-400/10 dark:bg-[#07150f]/90 dark:shadow-[0_24px_80px_-40px_rgba(16,185,129,0.3)]">
              <p className="mb-2 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
                Configuration
              </p>
              <h2 className="mb-3 text-3xl font-semibold text-slate-950 dark:text-emerald-50">
                Service endpoints
              </h2>
              <p className="mb-6 text-sm leading-6 text-slate-600 dark:text-emerald-100/65">
                These values are already wired through the app bridge, so
                swapping to real backend implementations later will be
                straightforward.
              </p>

              <form className="grid grid-cols-2 space-x-4" onSubmit={onSubmit}>
                <div>
                  <label
                    className="mb-2 block text-sm font-semibold text-slate-800 dark:text-emerald-100"
                    htmlFor="whisper-url"
                  >
                    Whisper API URL
                  </label>
                  <input
                    id="whisper-url"
                    className="w-full rounded-2xl border border-emerald-500/12 bg-white px-4 py-3 text-slate-900 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50"
                    value={form.whisperUrl}
                    onChange={(event) =>
                      onChange("whisperUrl", event.target.value)
                    }
                  />
                </div>

                <div>
                  <label
                    className="mb-2 block text-sm font-semibold text-slate-800 dark:text-emerald-100"
                    htmlFor="ollama-url"
                  >
                    Ollama URL
                  </label>
                  <input
                    id="ollama-url"
                    className="w-full rounded-2xl border border-emerald-500/12 bg-white px-4 py-3 text-slate-900 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50"
                    value={form.ollamaUrl}
                    onChange={(event) =>
                      onChange("ollamaUrl", event.target.value)
                    }
                  />
                </div>

                <div>
                  <label
                    className="mb-2 block text-sm font-semibold text-slate-800 dark:text-emerald-100"
                    htmlFor="cloud-url"
                  >
                    Cloud URL
                  </label>
                  <input
                    id="cloud-url"
                        className="w-full rounded-2xl border border-emerald-500/12 bg-white px-4 py-3 text-slate-900 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50"
                    value={form.cloudUrl}
                    onChange={(event) =>
                      onChange("cloudUrl", event.target.value)
                    }
                  />
                </div>

                <div>
                  <label
                    className="mb-2 block text-sm font-semibold text-slate-800 dark:text-emerald-100"
                    htmlFor="output-dir"
                  >
                    Output directory
                  </label>
                  <input
                    id="output-dir"
                    className="w-full rounded-2xl border border-emerald-500/12 bg-white px-4 py-3 text-slate-900 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50"
                    value={form.outputDir}
                    onChange={(event) =>
                      onChange("outputDir", event.target.value)
                    }
                  />
                </div>

                <div>
                  <label
                    className="mb-2 block text-sm font-semibold text-slate-800 dark:text-emerald-100"
                    htmlFor="theme"
                  >
                    Theme
                  </label>
                  <select
                    id="theme"
                    className="w-full rounded-2xl border border-emerald-500/12 bg-white px-4 py-3 text-slate-900 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50"
                    value={form.theme}
                    onChange={(event) => onChange("theme", event.target.value)}
                  >
                    <option value="midnight">Midnight</option>
                    <option value="paper">Paper</option>
                    <option value="graphite">Graphite</option>
                  </select>
                </div>

                <button
                  className="mt-6 inline-flex w-full items-center justify-center rounded-2xl bg-slate-950 px-4 py-3 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-emerald-400 dark:text-slate-950 dark:hover:bg-emerald-300"
                  type="submit"
                  disabled={saving}
                >
                  {saving ? "Saving..." : "Save settings"}
                </button>
              </form>
            </section>
          </section>
        </main>
      </div>
    </div>
  );
}
