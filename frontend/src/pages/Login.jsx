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
    <main className="min-h-screen flex items-center justify-center p-4 bg-white dark:bg-black">
      <section className="w-full max-w-md bg-white dark:bg-black border-2 border-emerald-500 dark:border-emerald-700 rounded-xl p-8">
        <p className="text-xs uppercase tracking-widest text-emerald-500/70 font-semibold mb-2">
          traer access
        </p>
        <h1 className="text-2xl font-bold text-emerald-600 dark:text-emerald-300 mb-3">
          Sign in to your operator workspace
        </h1>
        <p className="text-sm text-emerald-500/70 dark:text-emerald-400/70 mb-6">
          This starter uses a lightweight mock auth flow so you can wire real
          identity later without reshaping the UI.
        </p>

        <form className="space-y-4" onSubmit={onSubmit}>
          <div>
            <label
              className="block text-sm font-semibold text-emerald-600 dark:text-emerald-300 mb-2"
              htmlFor="username"
            >
              Username
            </label>
            <input
              id="username"
              className="w-full px-4 py-2 border-2 border-emerald-100 dark:border-emerald-800 bg-white dark:bg-black text-emerald-600 dark:text-emerald-300 rounded-md focus:border-emerald-500 dark:focus:border-emerald-400 outline-none transition-colors"
              value={form.username}
              onChange={(event) =>
                setForm((state) => ({ ...state, username: event.target.value }))
              }
              placeholder="operator"
            />
          </div>

          <div>
            <label
              className="block text-sm font-semibold text-emerald-600 dark:text-emerald-300 mb-2"
              htmlFor="password"
            >
              Password
            </label>
            <input
              id="password"
              type="password"
              className="w-full px-4 py-2 border-2 border-emerald-100 dark:border-emerald-800 bg-white dark:bg-black text-emerald-600 dark:text-emerald-300 rounded-md focus:border-emerald-500 dark:focus:border-emerald-400 outline-none transition-colors"
              value={form.password}
              onChange={(event) =>
                setForm((state) => ({ ...state, password: event.target.value }))
              }
              placeholder="••••••••"
            />
          </div>

          <button
            className="w-full mt-6 px-4 py-2 bg-gradient-to-r from-emerald-600 to-emerald-500 text-white font-semibold rounded-md hover:from-emerald-700 hover:to-emerald-600 disabled:opacity-50 disabled:cursor-not-allowed transition-all duration-200"
            type="submit"
            disabled={pending}
          >
            {pending ? "Signing in..." : "Sign in"}
          </button>
        </form>
      </section>
    </main>
  );
}
