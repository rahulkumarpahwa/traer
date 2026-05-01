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
    <section className="fixed bottom-0 left-0 right-0 max-h-56 bg-white dark:bg-black border-t border-emerald-100 dark:border-emerald-700 overflow-auto text-sm font-mono text-emerald-600">
      <div
        className="flex items-center justify-between p-2 border-b border-emerald-50 dark:border-emerald-800"
        onClick={toggleTerminal}
        title="Click to minimize"
      >
        <span className="text-emerald-700 dark:text-emerald-300">
          Runtime terminal
        </span>
        <span className="text-emerald-500/70">
          {terminalLines.length} entries
        </span>
      </div>
      <div className="p-3 space-y-1" ref={bodyRef}>
        {terminalLines.map((line, index) => (
          <p
            key={`${line}-${index}`}
            className="text-emerald-600 dark:text-emerald-400"
          >
            {line}
          </p>
        ))}
      </div>
    </section>
  );
}
