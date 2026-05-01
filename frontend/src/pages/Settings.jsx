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

        <main className="content">
          <section className="panel">
            <p className="eyebrow">Configuration</p>
            <h2>Service endpoints</h2>
            <p className="muted">
              These values are already wired through the app bridge, so swapping
              to real backend implementations later will be straightforward.
            </p>

            <form className="settings-form" onSubmit={onSubmit}>
              <label className="input-label" htmlFor="whisper-url">
                Whisper API URL
              </label>
              <input
                id="whisper-url"
                className="input-field"
                value={form.whisperUrl}
                onChange={(event) => onChange("whisperUrl", event.target.value)}
              />

              <label className="input-label" htmlFor="ollama-url">
                Ollama URL
              </label>
              <input
                id="ollama-url"
                className="input-field"
                value={form.ollamaUrl}
                onChange={(event) => onChange("ollamaUrl", event.target.value)}
              />

              <label className="input-label" htmlFor="cloud-url">
                Cloud URL
              </label>
              <input
                id="cloud-url"
                className="input-field"
                value={form.cloudUrl}
                onChange={(event) => onChange("cloudUrl", event.target.value)}
              />

              <label className="input-label" htmlFor="output-dir">
                Output directory
              </label>
              <input
                id="output-dir"
                className="input-field"
                value={form.outputDir}
                onChange={(event) => onChange("outputDir", event.target.value)}
              />

              <label className="input-label" htmlFor="theme">
                Theme
              </label>
              <select
                id="theme"
                className="select-field"
                value={form.theme}
                onChange={(event) => onChange("theme", event.target.value)}
              >
                <option value="midnight">Midnight</option>
                <option value="paper">Paper</option>
                <option value="graphite">Graphite</option>
              </select>

              <button
                className="primary-button"
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
