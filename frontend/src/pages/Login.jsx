import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { useStore } from "../store/useStore";

export default function Login() {
  const navigate = useNavigate();
  const loginUser = useStore((state) => state.loginUser);
  const [form, setForm] = useState({
    username: "",
    password: "",
  });
  const [pending, setPending] = useState(false);

  const onSubmit = async (event) => {
    event.preventDefault();
    setPending(true);
    try {
      await loginUser(form);
      navigate("/");
    } finally {
      setPending(false);
    }
  };

  return (
    <main className="auth-shell">
      <section className="auth-card">
        <p className="eyebrow">traer access</p>
        <h1>Sign in to your operator workspace</h1>
        <p className="muted">
          This starter uses a lightweight mock auth flow so you can wire real identity later without reshaping the UI.
        </p>

        <form className="auth-form" onSubmit={onSubmit}>
          <label className="input-label" htmlFor="username">
            Username
          </label>
          <input
            id="username"
            className="input-field"
            value={form.username}
            onChange={(event) => setForm((state) => ({ ...state, username: event.target.value }))}
            placeholder="operator"
          />

          <label className="input-label" htmlFor="password">
            Password
          </label>
          <input
            id="password"
            type="password"
            className="input-field"
            value={form.password}
            onChange={(event) => setForm((state) => ({ ...state, password: event.target.value }))}
            placeholder="••••••••"
          />

          <button className="primary-button" type="submit" disabled={pending}>
            {pending ? "Signing in..." : "Sign in"}
          </button>
        </form>
      </section>
    </main>
  );
}
