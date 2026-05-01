import Header from "../components/Header";
import NotificationTray from "../components/NotificationTray";
import Sidebar from "../components/Sidebar";
import Terminal from "../components/Terminal";
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

        <main className="flex-1 p-8 overflow-auto">
          <section className="bg-white dark:bg-black border-2 border-emerald-500 dark:border-emerald-700 rounded-xl p-8 max-w-2xl">
            <p className="text-xs uppercase tracking-widest text-emerald-500/70 font-semibold mb-2">
              Configuration
            </p>
            <h2 className="text-2xl font-bold text-emerald-600 dark:text-emerald-300 mb-3">
              Service endpoints
            </h2>
            <p className="text-sm text-emerald-500/70 dark:text-emerald-400/70 mb-6">
              These values are already wired through the app bridge, so swapping
              to real backend implementations later will be straightforward.
            </p>

            <form className="space-y-4" onSubmit={onSubmit}>
              <div>
                <label
                  className="block text-sm font-semibold text-emerald-600 dark:text-emerald-300 mb-2"
                  htmlFor="whisper-url"
                >
                  Whisper API URL
                </label>
                <input
                  id="whisper-url"
                  className="w-full px-4 py-2 border-2 border-emerald-100 dark:border-emerald-800 bg-white dark:bg-black text-emerald-600 dark:text-emerald-300 rounded-md focus:border-emerald-500 dark:focus:border-emerald-400 outline-none transition-colors"
                  value={form.whisperUrl}
                  onChange={(event) =>
                    onChange("whisperUrl", event.target.value)
                  }
                />
              </div>

              <div>
                <label
                  className="block text-sm font-semibold text-emerald-600 dark:text-emerald-300 mb-2"
                  htmlFor="ollama-url"
                >
                  Ollama URL
                </label>
                <input
                  id="ollama-url"
                  className="w-full px-4 py-2 border-2 border-emerald-100 dark:border-emerald-800 bg-white dark:bg-black text-emerald-600 dark:text-emerald-300 rounded-md focus:border-emerald-500 dark:focus:border-emerald-400 outline-none transition-colors"
                  value={form.ollamaUrl}
                  onChange={(event) =>
                    onChange("ollamaUrl", event.target.value)
                  }
                />
              </div>

              <div>
                <label
                  className="block text-sm font-semibold text-emerald-600 dark:text-emerald-300 mb-2"
                  htmlFor="cloud-url"
                >
                  Cloud URL
                </label>
                <input
                  id="cloud-url"
                  className="w-full px-4 py-2 border-2 border-emerald-100 dark:border-emerald-800 bg-white dark:bg-black text-emerald-600 dark:text-emerald-300 rounded-md focus:border-emerald-500 dark:focus:border-emerald-400 outline-none transition-colors"
                  value={form.cloudUrl}
                  onChange={(event) => onChange("cloudUrl", event.target.value)}
                />
              </div>

              <div>
                <label
                  className="block text-sm font-semibold text-emerald-600 dark:text-emerald-300 mb-2"
                  htmlFor="output-dir"
                >
                  Output directory
                </label>
                <input
                  id="output-dir"
                  className="w-full px-4 py-2 border-2 border-emerald-100 dark:border-emerald-800 bg-white dark:bg-black text-emerald-600 dark:text-emerald-300 rounded-md focus:border-emerald-500 dark:focus:border-emerald-400 outline-none transition-colors"
                  value={form.outputDir}
                  onChange={(event) =>
                    onChange("outputDir", event.target.value)
                  }
                />
              </div>

              <div>
                <label
                  className="block text-sm font-semibold text-emerald-600 dark:text-emerald-300 mb-2"
                  htmlFor="theme"
                >
                  Theme
                </label>
                <select
                  id="theme"
                  className="w-full px-4 py-2 border-2 border-emerald-100 dark:border-emerald-800 bg-white dark:bg-black text-emerald-600 dark:text-emerald-300 rounded-md focus:border-emerald-500 dark:focus:border-emerald-400 outline-none transition-colors"
                  value={form.theme}
                  onChange={(event) => onChange("theme", event.target.value)}
                >
                  <option value="midnight">Midnight</option>
                  <option value="paper">Paper</option>
                  <option value="graphite">Graphite</option>
                </select>
              </div>

              <button
                className="w-full mt-6 px-4 py-2 bg-gradient-to-r from-emerald-600 to-emerald-500 text-white font-semibold rounded-md hover:from-emerald-700 hover:to-emerald-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
                type="submit"
                disabled={saving}
              >
                {saving ? "Saving..." : "Save settings"}
              </button>
            </form>
          </section>
        </main>
      </div>

      <Terminal />
    </div>
  );
}
