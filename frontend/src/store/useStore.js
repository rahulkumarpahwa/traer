import { useEffect } from "react";
import { create } from "zustand";
import {
  getActiveJobs,
  getDownloadUrl,
  getInitialState,
  getJobStatus,
  getProfile,
  login,
  saveSettings,
  startJob,
} from "../lib/backend";

const actionOptions = {
  transcribe: ["Markdown", "MDX", "TXT", "JSON"],
  audio: ["MP3", "WAV", "AAC", "FLAC"],
  video: ["360p", "720p", "1080p", "4K"],
  cloud: ["Audio", "Video", "Notes", "Bundle"],
};

const initialSettings = {
  whisperUrl: "",
  ollamaUrl: "",
  cloudUrl: "",
  outputDir: "",
  theme: "midnight",
};

function mergeJobs(existingJobs, incomingJobs) {
  const byId = new Map(existingJobs.map((job) => [job.id, job]));

  incomingJobs.forEach((job) => {
    const previous = byId.get(job.id) ?? {};
    byId.set(job.id, {
      ...previous,
      ...job,
      option: job.option || previous.option || String(job.kind || "").toUpperCase(),
      title: job.title || previous.title,
      message: job.message || previous.message,
      createdAt: previous.createdAt || job.createdAt || new Date().toISOString(),
    });
  });

  return Array.from(byId.values()).sort(
    (left, right) => new Date(right.createdAt).getTime() - new Date(left.createdAt).getTime(),
  );
}

export const useStore = create((set, get) => ({
  appName: "traer",
  hydrated: false,
  sidebarOpen: true,
  capabilities: {
    backendBound: false,
    mockMode: true,
  },
  user: null,
  profile: null,
  settings: initialSettings,
  recentLinks: [],
  suggestions: [],
  notifications: [],
  terminalLines: ["[system] traer ready"],
  activeJobs: [],
  selectedOptions: {
    transcribe: "Markdown",
    audio: "MP3",
    video: "1080p",
    cloud: "Bundle",
  },
  initialize: async () => {
    const data = await getInitialState();
    set({
      hydrated: true,
      appName: data.appName,
      recentLinks: data.recentLinks ?? [],
      settings: data.settings ?? initialSettings,
      profile: data.profile ?? null,
      suggestions: data.suggestions ?? [],
      capabilities: data.capabilities ?? {
        backendBound: false,
        mockMode: true,
      },
      terminalLines: [
        "[system] traer ready",
        data.capabilities?.backendBound
          ? "[bridge] wails backend detected"
          : "[bridge] running in mock preview mode",
      ],
    });

    await get().refreshActiveJobs();
  },
  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
  addNotification: (notification) =>
    set((state) => ({
      notifications: [
        { id: crypto.randomUUID(), tone: "info", ...notification },
        ...state.notifications,
      ].slice(0, 4),
    })),
  dismissNotification: (id) =>
    set((state) => ({
      notifications: state.notifications.filter((notification) => notification.id !== id),
    })),
  addTerminalLine: (line) =>
    set((state) => ({
      terminalLines: [...state.terminalLines, line].slice(-120),
    })),
  addLink: (link) => {
    const trimmed = link.trim();
    if (!trimmed) return;

    set((state) => ({
      recentLinks: [trimmed, ...state.recentLinks.filter((item) => item !== trimmed)].slice(0, 8),
    }));
  },
  setOption: (kind, option) =>
    set((state) => ({
      selectedOptions: {
        ...state.selectedOptions,
        [kind]: option,
      },
    })),
  loginUser: async (payload) => {
    const user = await login(payload);
    set({ user });
    get().addNotification({
      tone: "success",
      title: `Welcome back, ${user.name}`,
      message: "Your workspace is ready.",
    });
    get().addTerminalLine(`[auth] signed in as ${user.email}`);
    return user;
  },
  saveAppSettings: async (payload) => {
    const settings = await saveSettings(payload);
    set({ settings });
    get().addNotification({
      tone: "success",
      title: "Settings saved",
      message: "traer will use the updated endpoints for the next run.",
    });
    get().addTerminalLine(`[config] endpoints updated for theme=${settings.theme}`);
    return settings;
  },
  refreshProfile: async () => {
    const profile = await getProfile();
    set({ profile });
    return profile;
  },
  refreshActiveJobs: async () => {
    try {
      const remoteActiveJobs = await getActiveJobs();
      const trackedJobs = get().activeJobs.filter((job) => job.source === "backend");
      const trackedStatuses = await Promise.all(
        trackedJobs.map(async (job) => {
          try {
            const status = await getJobStatus(job.id);
            return {
              ...status,
              option: job.option,
              createdAt: job.createdAt,
              source: "backend",
              downloadUrl:
                status.downloadUrl || (status.status === "done" ? await getDownloadUrl(job.id) : ""),
            };
          } catch {
            return job;
          }
        }),
      );

      const normalizedRemote = remoteActiveJobs.map((job) => ({
        ...job,
        source: "backend",
      }));

      set((state) => ({
        activeJobs: mergeJobs(state.activeJobs, [...normalizedRemote, ...trackedStatuses]).slice(0, 12),
      }));
    } catch {
      // Ignore polling errors so the UI can keep rendering existing jobs.
    }
  },
  runJob: async ({ kind, url }) => {
    const option = get().selectedOptions[kind];
    const cleanUrl = url.trim();

    if (!cleanUrl) {
      get().addNotification({
        tone: "warning",
        title: "Source required",
        message: "Paste a media link before starting a workflow.",
      });
      return;
    }

    get().addLink(cleanUrl);
    get().addTerminalLine(`[queue] ${kind} requested for ${cleanUrl}`);

    const response = await startJob({
      kind,
      option,
      url: cleanUrl,
    });

    response.logs?.forEach((line) => get().addTerminalLine(line));

    const job = {
      id: response.id,
      kind: response.kind,
      option: response.option,
      title: response.title,
      message: response.message,
      progress: Number(response.progress ?? 8),
      status: response.status ?? "running",
      output: response.output ?? "",
      fileName: response.fileName ?? "",
      downloadUrl: response.downloadUrl ?? "",
      sourceUrl: response.sourceUrl ?? cleanUrl,
      createdAt: response.createdAt,
      source: response.kind === "audio" || response.kind === "video" ? "backend" : "mock",
    };

    set((state) => ({
      activeJobs: mergeJobs(state.activeJobs, [job]).slice(0, 12),
    }));

    get().addNotification({
      tone: "info",
      title: response.title,
      message: response.message,
    });

    if (job.source === "backend") {
      await get().refreshActiveJobs();
      return;
    }

    const progressMarks = [24, 47, 71, 92, 100];
    progressMarks.forEach((value, index) => {
      setTimeout(() => {
        set((state) => ({
          activeJobs: state.activeJobs.map((activeJob) =>
            activeJob.id === job.id
              ? {
                  ...activeJob,
                  progress: value,
                  status: value === 100 ? "done" : "running",
                }
              : activeJob,
          ),
        }));

        get().addTerminalLine(
          value === 100
            ? `[done] ${job.kind} completed with ${job.option}`
            : `[work] ${job.kind} progress ${value}%`,
        );

        if (value === 100) {
          get().addNotification({
            tone: "success",
            title: `${job.title.replace("queued", "complete")}`,
            message: `${job.kind} finished successfully with preset ${job.option}.`,
          });
        }
      }, 700 * (index + 1));
    });
  },
  actionOptions,
}));

export function useHydrateApp() {
  const hydrated = useStore((state) => state.hydrated);
  const initialize = useStore((state) => state.initialize);
  const refreshActiveJobs = useStore((state) => state.refreshActiveJobs);

  useEffect(() => {
    if (!hydrated) {
      initialize();
    }
  }, [hydrated, initialize]);

  useEffect(() => {
    if (!hydrated) return undefined;

    const intervalId = window.setInterval(() => {
      refreshActiveJobs();
    }, 3000);

    return () => window.clearInterval(intervalId);
  }, [hydrated, refreshActiveJobs]);
}
