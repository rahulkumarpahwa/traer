import { useEffect, useRef } from "react";
import { useStore } from "../store/useStore";

export default function Terminal() {
  const terminalOpen = useStore((state) => state.terminalOpen);
  const terminalLines = useStore((state) => state.terminalLines);
  const bodyRef = useRef(null);

  useEffect(() => {
    if (bodyRef.current) {
      bodyRef.current.scrollTop = bodyRef.current.scrollHeight;
    }
  }, [terminalLines]);

  if (!terminalOpen) {
    return null;
  }

  return (
    <section className="terminal">
      <div className="terminal-header">
        <span>Runtime terminal</span>
        <span className="muted">{terminalLines.length} entries</span>
      </div>
      <div className="terminal-body" ref={bodyRef}>
        {terminalLines.map((line, index) => (
          <p key={`${line}-${index}`}>{line}</p>
        ))}
      </div>
    </section>
  );
}
