import { useEffect } from "react";
import { create } from "zustand";
import { getInitialState, getProfile, login, saveSettings, startJob } from "../lib/backend";

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

export const useStore = create((set, get) => ({
  appName: "traer",
  hydrated: false,
  sidebarOpen: true,
  terminalOpen: true,
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
        data.capabilities?.backendBound ? "[bridge] wails backend detected" : "[bridge] running in mock preview mode",
      ],
    });
  },
  toggleSidebar: () => set((state) => ({ sidebarOpen: !state.sidebarOpen })),
  toggleTerminal: () => set((state) => ({ terminalOpen: !state.terminalOpen })),
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

    response.logs.forEach((line) => get().addTerminalLine(line));

    const job = {
      id: response.id,
      kind: response.kind,
      option: response.option,
      title: response.title,
      message: response.message,
      progress: 8,
      status: "running",
      createdAt: response.createdAt,
    };

    set((state) => ({
      activeJobs: [job, ...state.activeJobs].slice(0, 6),
    }));

    get().addNotification({
      tone: "info",
      title: response.title,
      message: response.message,
    });

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

  useEffect(() => {
    if (!hydrated) {
      initialize();
    }
  }, [hydrated, initialize]);
}
