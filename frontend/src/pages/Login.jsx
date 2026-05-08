import { useState } from "react";
import { ArrowRight, KeyRound, ShieldCheck, UserPlus } from "lucide-react";
import { useNavigate } from "react-router-dom";
import { useStore } from "../store/useStore";

export default function Login() {
  const navigate = useNavigate();
  const loginUser = useStore((state) => state.loginUser);
  const signupUser = useStore((state) => state.signupUser);
  const [mode, setMode] = useState("login");
  const [pending, setPending] = useState(false);
  const [loginForm, setLoginForm] = useState({
    identifier: "",
    password: "",
  });
  const [signupForm, setSignupForm] = useState({
    username: "",
    email: "",
    password: "",
  });

  const handleLogin = async (event) => {
    event.preventDefault();
    setPending(true);
    try {
      await loginUser(loginForm);
      navigate("/");
    } finally {
      setPending(false);
    }
  };

  const handleSignup = async (event) => {
    event.preventDefault();
    setPending(true);
    try {
      await signupUser(signupForm);
      navigate("/");
    } finally {
      setPending(false);
    }
  };

  return (
    <main className="relative min-h-screen overflow-hidden bg-[#eef8f2] px-4 py-10 dark:bg-[#03110a]">
      <div className="pointer-events-none absolute inset-0 bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.18),transparent_24%),radial-gradient(circle_at_bottom_right,rgba(6,182,212,0.14),transparent_22%),linear-gradient(180deg,rgba(255,255,255,0.88),rgba(236,253,245,0.72))] dark:bg-[radial-gradient(circle_at_top_left,rgba(16,185,129,0.16),transparent_22%),radial-gradient(circle_at_bottom_right,rgba(6,182,212,0.12),transparent_20%),linear-gradient(180deg,rgba(3,17,10,1),rgba(5,25,17,0.96))]" />
      <section className="relative mx-auto grid min-h-[calc(100vh-5rem)] w-full max-w-6xl items-center gap-8 xl:grid-cols-[1.05fr_0.95fr]">
        <div className="rounded-[36px] border border-emerald-500/12 bg-slate-950 p-8 text-white shadow-[0_30px_90px_-45px_rgba(2,6,23,0.65)] dark:border-emerald-400/10 dark:bg-[#0b1f15] md:p-10">
          <p className="mb-3 text-[11px] font-semibold uppercase tracking-[0.28em] text-emerald-200/60">
            traer access
          </p>
          <h1 className="max-w-xl text-4xl font-semibold tracking-tight text-white dark:text-emerald-50 md:text-5xl">
            Sign in to your media workspace with the backend session flow.
          </h1>
          <p className="mt-5 max-w-xl text-base leading-7 text-emerald-50/75">
            Authentication is now backed by the Go API. Successful login stores the JWT in an HttpOnly cookie, and protected job requests reuse that session automatically.
          </p>

          <div className="mt-8 space-y-4">
            {[
              {
                icon: ShieldCheck,
                title: "Cookie-backed session",
                description: "The backend issues a JWT cookie and validates it on every protected route.",
              },
              {
                icon: KeyRound,
                title: "Login by username or email",
                description: "The form accepts either identifier plus password, matching the backend contract.",
              },
              {
                icon: UserPlus,
                title: "Create an account first",
                description: "Signup posts to `/users`, then signs you in so the protected UI is ready immediately.",
              },
            ].map(({ icon: Icon, title, description }) => (
              <div
                key={title}
                className="rounded-2xl border border-white/10 bg-white/5 p-4"
              >
                <div className="mb-3 flex h-10 w-10 items-center justify-center rounded-2xl border border-white/10 bg-white/10 text-emerald-100">
                  <Icon className="h-4 w-4" strokeWidth={1.8} />
                </div>
                <h2 className="text-base font-semibold text-white dark:text-emerald-50">
                  {title}
                </h2>
                <p className="mt-1 text-sm leading-6 text-emerald-50/70">
                  {description}
                </p>
              </div>
            ))}
          </div>
        </div>

        <section className="rounded-[36px] border border-emerald-500/15 bg-white/88 p-8 shadow-[0_24px_80px_-40px_rgba(5,46,22,0.45)] backdrop-blur-sm dark:border-emerald-400/10 dark:bg-[#07150f]/92 dark:shadow-[0_24px_80px_-40px_rgba(16,185,129,0.28)] md:p-10">
          <div className="mb-8 inline-flex rounded-full border border-emerald-500/12 bg-emerald-500/[0.04] p-1 dark:border-emerald-400/10 dark:bg-white/[0.03]">
            <button
              type="button"
              onClick={() => setMode("login")}
              className={`${mode === "login" ? "bg-slate-950 text-white dark:bg-emerald-400 dark:text-slate-950" : "text-slate-600 dark:text-emerald-100/70"} rounded-full px-4 py-2 text-sm font-semibold transition`}
            >
              Login
            </button>
            <button
              type="button"
              onClick={() => setMode("signup")}
              className={`${mode === "signup" ? "bg-slate-950 text-white dark:bg-emerald-400 dark:text-slate-950" : "text-slate-600 dark:text-emerald-100/70"} rounded-full px-4 py-2 text-sm font-semibold transition`}
            >
              Sign up
            </button>
          </div>

          {mode === "login" ? (
            <form className="space-y-5" onSubmit={handleLogin}>
              <div>
                <label
                  className="mb-2 block text-sm font-semibold text-slate-800 dark:text-emerald-100"
                  htmlFor="identifier"
                >
                  Username or email
                </label>
                <input
                  id="identifier"
                  className="w-full rounded-2xl border border-emerald-500/12 bg-white px-4 py-3 text-slate-900 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50"
                  value={loginForm.identifier}
                  onChange={(event) =>
                    setLoginForm((state) => ({
                      ...state,
                      identifier: event.target.value,
                    }))
                  }
                  placeholder="operator or operator@example.com"
                />
              </div>

              <div>
                <label
                  className="mb-2 block text-sm font-semibold text-slate-800 dark:text-emerald-100"
                  htmlFor="login-password"
                >
                  Password
                </label>
                <input
                  id="login-password"
                  type="password"
                  className="w-full rounded-2xl border border-emerald-500/12 bg-white px-4 py-3 text-slate-900 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50"
                  value={loginForm.password}
                  onChange={(event) =>
                    setLoginForm((state) => ({
                      ...state,
                      password: event.target.value,
                    }))
                  }
                  placeholder="Enter your password"
                />
              </div>

              <button
                className="inline-flex w-full items-center justify-center gap-2 rounded-2xl bg-slate-950 px-4 py-3 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-emerald-400 dark:text-slate-950 dark:hover:bg-emerald-300"
                type="submit"
                disabled={pending}
              >
                {pending ? "Signing in..." : "Sign in"}
                <ArrowRight className="h-4 w-4" strokeWidth={1.8} />
              </button>
            </form>
          ) : (
            <form className="space-y-5" onSubmit={handleSignup}>
              <div>
                <label
                  className="mb-2 block text-sm font-semibold text-slate-800 dark:text-emerald-100"
                  htmlFor="signup-username"
                >
                  Username
                </label>
                <input
                  id="signup-username"
                  className="w-full rounded-2xl border border-emerald-500/12 bg-white px-4 py-3 text-slate-900 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50"
                  value={signupForm.username}
                  onChange={(event) =>
                    setSignupForm((state) => ({
                      ...state,
                      username: event.target.value,
                    }))
                  }
                  placeholder="Choose a username"
                />
              </div>

              <div>
                <label
                  className="mb-2 block text-sm font-semibold text-slate-800 dark:text-emerald-100"
                  htmlFor="signup-email"
                >
                  Email
                </label>
                <input
                  id="signup-email"
                  type="email"
                  className="w-full rounded-2xl border border-emerald-500/12 bg-white px-4 py-3 text-slate-900 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50"
                  value={signupForm.email}
                  onChange={(event) =>
                    setSignupForm((state) => ({
                      ...state,
                      email: event.target.value,
                    }))
                  }
                  placeholder="name@example.com"
                />
              </div>

              <div>
                <label
                  className="mb-2 block text-sm font-semibold text-slate-800 dark:text-emerald-100"
                  htmlFor="signup-password"
                >
                  Password
                </label>
                <input
                  id="signup-password"
                  type="password"
                  className="w-full rounded-2xl border border-emerald-500/12 bg-white px-4 py-3 text-slate-900 outline-none transition focus:border-emerald-500/40 focus:ring-4 focus:ring-emerald-500/10 dark:border-emerald-400/10 dark:bg-[#0a1912] dark:text-emerald-50"
                  value={signupForm.password}
                  onChange={(event) =>
                    setSignupForm((state) => ({
                      ...state,
                      password: event.target.value,
                    }))
                  }
                  placeholder="Create a secure password"
                />
              </div>

              <button
                className="inline-flex w-full items-center justify-center gap-2 rounded-2xl bg-slate-950 px-4 py-3 text-sm font-semibold text-white transition hover:bg-emerald-600 disabled:cursor-not-allowed disabled:opacity-50 dark:bg-emerald-400 dark:text-slate-950 dark:hover:bg-emerald-300"
                type="submit"
                disabled={pending}
              >
                {pending ? "Creating account..." : "Create account"}
                <ArrowRight className="h-4 w-4" strokeWidth={1.8} />
              </button>
            </form>
          )}
        </section>
      </section>
    </main>
  );
}
