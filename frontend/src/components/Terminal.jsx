import { useEffect, useRef } from "react";
import { useStore } from "../store/useStore";
import { TerminalSquare } from "lucide-react";

export default function Terminal() {
  const terminalLines = useStore((state) => state.terminalLines);
  const bodyRef = useRef(null);

  useEffect(() => {
    if (bodyRef.current) {
      bodyRef.current.scrollTop = bodyRef.current.scrollHeight;
    }
  }, [terminalLines]);

  return (
    <section className="mx-auto max-w-6xl rounded-[32px] border border-emerald-500/15 bg-[#f9fcfa]/95 p-8 shadow-[0_24px_80px_-40px_rgba(5,46,22,0.45)] backdrop-blur-sm dark:border-emerald-400/10 dark:bg-[#03130c] dark:shadow-[0_24px_80px_-40px_rgba(16,185,129,0.28)]">
      <div className="mb-6 flex items-center justify-between gap-4">
        <div className="flex items-center gap-4">
          <span className="flex h-12 w-12 items-center justify-center rounded-2xl border border-emerald-500/12 bg-emerald-500/[0.06] text-emerald-700 dark:border-emerald-400/10 dark:bg-white/[0.03] dark:text-emerald-200">
            <TerminalSquare className="h-5 w-5" strokeWidth={1.8} />
          </span>
          <div>
            <p className="mb-1 text-[11px] font-semibold uppercase tracking-[0.26em] text-emerald-700/55 dark:text-emerald-200/55">
              Runtime
            </p>
            <h2 className="text-2xl font-semibold text-slate-950 dark:text-emerald-50">
              Terminal
            </h2>
          </div>
        </div>
        <span className="rounded-full border border-emerald-500/12 bg-emerald-500/[0.05] px-3 py-1 text-xs font-semibold uppercase tracking-[0.18em] text-emerald-700/70 dark:border-emerald-400/10 dark:bg-white/[0.03] dark:text-emerald-200/80">
          {terminalLines.length} entries
        </span>
      </div>

      <div
        ref={bodyRef}
        className="max-h-[64vh] overflow-auto rounded-[28px] border border-emerald-500/12 bg-[#0b1711] p-4 font-mono text-sm leading-7 text-emerald-50 shadow-inner shadow-black/20 dark:border-emerald-400/10 dark:bg-[#010804]"
      >
        {terminalLines.map((line, index) => (
          <p
            key={`${line}-${index}`}
            className="border-b border-emerald-400/10 py-1 text-emerald-100/90 last:border-b-0"
          >
            {line}
          </p>
        ))}
      </div>
    </section>
  );
}
