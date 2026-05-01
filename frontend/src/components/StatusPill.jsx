export default function StatusPill({ children, tone = "neutral" }) {
  const toneMap = {
    neutral:
      "border border-emerald-500/15 bg-white/70 text-emerald-700 shadow-sm shadow-emerald-950/5 dark:border-emerald-400/15 dark:bg-white/5 dark:text-emerald-200",
    strong:
      "border border-emerald-500/20 bg-emerald-500/10 text-emerald-700 shadow-sm shadow-emerald-950/5 dark:border-emerald-400/15 dark:bg-emerald-400/10 dark:text-emerald-200",
    warm:
      "border border-amber-500/20 bg-amber-500/10 text-amber-700 shadow-sm shadow-amber-950/5 dark:border-amber-400/15 dark:bg-amber-400/10 dark:text-amber-200",
  };

  return (
    <span
      className={`inline-flex items-center rounded-full px-3 py-1 text-xs font-semibold tracking-[0.16em] uppercase ${toneMap[tone]}`}
    >
      {children}
    </span>
  );
}
