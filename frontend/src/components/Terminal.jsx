import { useEffect, useRef } from "react";
import { useStore } from "../store/useStore";

export default function Terminal() {
  const terminalOpen = useStore((state) => state.terminalOpen);
  const toggleTerminal = useStore((state) => state.toggleTerminal);
  const terminalLines = useStore((state) => state.terminalLines);
  const bodyRef = useRef(null);

  useEffect(() => {
    if (bodyRef.current) {
      bodyRef.current.scrollTop = bodyRef.current.scrollHeight;
    }
  }, [terminalLines]);

  if (!terminalOpen) return null;

  return (
    <section className="fixed bottom-0 left-0 right-0 max-h-64 overflow-auto border-t border-emerald-500/10 bg-[#f7faf8]/90 text-sm font-mono text-slate-700 backdrop-blur-xl dark:border-emerald-400/10 dark:bg-[#031009]/92 dark:text-emerald-100/80">
      <div
        className="flex items-center justify-between border-b border-emerald-500/10 px-3 py-3 dark:border-emerald-400/10"
        onClick={toggleTerminal}
        title="Click to minimize"
      >
        <span className="font-medium text-slate-900 dark:text-emerald-50">
          Runtime terminal
        </span>
        <span className="text-slate-500 dark:text-emerald-100/55">
          {terminalLines.length} entries
        </span>
      </div>
      <div className="p-3 space-y-1" ref={bodyRef}>
        {terminalLines.map((line, index) => (
          <p
            key={`${line}-${index}`}
            className="text-slate-700 dark:text-emerald-100/75"
          >
            {line}
          </p>
        ))}
      </div>
    </section>
  );
}
