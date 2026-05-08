import { HashRouter, Navigate, Route, Routes } from "react-router-dom";
import Dashboard from "./pages/Dashboard";
import Login from "./pages/Login";
import Profile from "./pages/Profile";
import Settings from "./pages/Settings";
import TerminalPage from "./pages/TerminalPage";
import { useHydrateApp, useStore } from "./store/useStore";
import JobPanel from "./components/JobPanel";

function ProtectedRoute({ children }) {
  const authReady = useStore((state) => state.authReady);
  const authenticated = useStore((state) => state.authenticated);

  if (!authReady) {
    return null;
  }

  if (!authenticated) {
    return <Navigate to="/login" replace />;
  }

  return children;
}

function LoginRoute() {
  const authReady = useStore((state) => state.authReady);
  const authenticated = useStore((state) => state.authenticated);

  if (!authReady) {
    return null;
  }

  return authenticated ? <Navigate to="/" replace /> : <Login />;
}

function AppRoutes() {
  useHydrateApp();

  return (
    <Routes>
      <Route path="/login" element={<LoginRoute />} />
      <Route
        path="/"
        element={
          <ProtectedRoute>
            <Dashboard />
          </ProtectedRoute>
        }
      />
      <Route
        path="/settings"
        element={
          <ProtectedRoute>
            <Settings />
          </ProtectedRoute>
        }
      />
      <Route
        path="/profile"
        element={
          <ProtectedRoute>
            <Profile />
          </ProtectedRoute>
        }
      />
      <Route
        path="/jobs"
        element={
          <ProtectedRoute>
            <JobPanel />
          </ProtectedRoute>
        }
      />
      <Route
        path="/terminal"
        element={
          <ProtectedRoute>
            <TerminalPage />
          </ProtectedRoute>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  );
}

export default function App() {
  return (
    <HashRouter>
      <AppRoutes />
    </HashRouter>
  );
}
