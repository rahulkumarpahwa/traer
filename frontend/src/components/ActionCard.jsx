import { useStore } from "../store/useStore";
import { CloudUpload, FileAudio2, FileText, Video } from "lucide-react";
import StatusPill from "./StatusPill";

const iconMap = {
  transcribe: FileText,
  audio: FileAudio2,
  video: Video,
  cloud: CloudUpload,
};

export default function ActionCard({
  kind,
  title,
  description,
  cta,
  disabled,
  url,
}) {
  const actionOptions = useStore((state) => state.actionOptions);
  const selectedOption = useStore((state) => state.selectedOptions[kind]);
  const setOption = useStore((state) => state.setOption);
  const runJob = useStore((state) => state.runJob);

  const IconComponent = iconMap[kind] ?? FileText;
  const optionCount = actionOptions[kind]?.length ?? 0;

  return (
    <article className="group relative overflow-hidden rounded-[28px] border border-emerald-500/15 bg-white/85 p-6 shadow-[0_24px_70px_-36px_rgba(5,46,22,0.45)] backdrop-blur-sm transition duration-300 hover:-translate-y-1 hover:border-emerald-500/25 hover:shadow-[0_32px_90px_-40px_rgba(5,46,22,0.55)] dark:border-emerald-400/10 dark:bg-[#07150f]/90 dark:shadow-[0_24px_70px_-36px_rgba(16,185,129,0.25)]">
      <div className="pointer-events-none absolute inset-x-0 top-0 h-px bg-gradient-to-r from-transparent via-emerald-500/40 to-transparent" />
      <div className="mb-5 flex items-start justify-between gap-4">
        <div>
          <p className="mb-2 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
            {kind}
          </p>
          <h3 className="flex items-center gap-3 text-lg font-semibold text-slate-900 dark:text-emerald-50">
            <span className="flex h-11 w-11 items-center justify-center rounded-2xl border border-emerald-500/15 bg-gradient-to-br from-emerald-500/15 via-teal-500/10 to-cyan-500/15 text-emerald-700 shadow-inner shadow-white/50 dark:border-emerald-400/10 dark:text-emerald-200 dark:shadow-none">
              <IconComponent className="h-5 w-5" strokeWidth={1.8} />
            </span>
            <span>{title}</span>
          </h3>
        </div>
        <StatusPill tone="neutral">{optionCount} presets</StatusPill>
      </div>
      <p className="mb-5 text-sm leading-6 text-slate-600 dark:text-emerald-100/70">
        {description}
      </p>
      <div className="mb-4 rounded-2xl border border-emerald-500/10 bg-emerald-500/[0.04] p-3 dark:border-emerald-400/10 dark:bg-white/[0.03]">
        <label
          htmlFor={`${kind}-option`}
          className="mb-2 block text-[11px] font-semibold uppercase tracking-[0.2em] text-emerald-700/60 dark:text-emerald-200/60"
        >
          Output preset
        </label>
        <select
          id={`${kind}-option`}
          className="w-full rounded-xl border border-emerald-500/10 bg-white px-3 py-3 text-sm font-medium text-slate-800 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-100"
          value={selectedOption}
          onChange={(event) => setOption(kind, event.target.value)}
        >
          {actionOptions[kind].map((option) => (
            <option key={option} value={option}>
              {option}
            </option>
          ))}
        </select>
      </div>
      <button
        className="mt-auto inline-flex w-full items-center justify-center rounded-2xl bg-slate-950 px-4 py-3 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-45 dark:bg-emerald-400 dark:text-slate-950 dark:hover:bg-emerald-300"
        onClick={() => runJob({ kind, url })}
        disabled={disabled}
      >
        {cta}
      </button>
    </article>
  );
}
