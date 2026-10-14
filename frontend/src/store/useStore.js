import { useEffect } from "react";
import { create } from "zustand";
import {
  createMusicJobs,
  downloadJobFile,
  getActiveJobs,
  getInitialState,
  getJobStatus,
  getProfile,
  getUserById,
  login,
  retryJob,
  saveSettings,
  signup,
  startJob,
} from "../lib/backend";

const USER_ID_STORAGE_KEY = "traer.userId";

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

function persistUserId(userId) {
  if (typeof window === "undefined") return;

  if (userId) {
    window.localStorage.setItem(USER_ID_STORAGE_KEY, String(userId));
  } else {
    window.localStorage.removeItem(USER_ID_STORAGE_KEY);
  }
}

function readPersistedUserId() {
  if (typeof window === "undefined") return 0;
  const raw = window.localStorage.getItem(USER_ID_STORAGE_KEY);
  return raw ? Number(raw) : 0;
}

export const useStore = create((set, get) => ({
  appName: "traer",
  hydrated: false,
  authReady: false,
  authenticated: false,
  currentUserId: 0,
  sidebarOpen: true,
  capabilities: {
    backendBound: false,
    mockMode: false,
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
    const persistedUserId = readPersistedUserId();

    set({
      hydrated: true,
      appName: data.appName,
      recentLinks: data.recentLinks ?? [],
      settings: data.settings ?? initialSettings,
      profile: data.profile ?? null,
      suggestions: data.suggestions ?? [],
      capabilities: data.capabilities ?? {
        backendBound: false,
        mockMode: false,
      },
      terminalLines: [
        "[system] traer ready",
        data.capabilities?.backendBound
          ? "[bridge] wails backend detected"
          : "[bridge] browser client active",
      ],
      currentUserId: persistedUserId,
    });

    if (persistedUserId) {
      try {
        const [user, profile] = await Promise.all([
          getUserById(persistedUserId),
          getProfile(persistedUserId),
        ]);

        set({
          authenticated: true,
          authReady: true,
          currentUserId: persistedUserId,
          user,
          profile,
        });

        await get().refreshActiveJobs();
        return;
      } catch {
        persistUserId(0);
      }
    }

    set({
      authReady: true,
      authenticated: false,
      currentUserId: 0,
      user: null,
      profile: null,
    });
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
    const sessionUser = await login(payload);
    const resolvedUser = sessionUser?.username && sessionUser?.email
      ? sessionUser
      : await getUserById(sessionUser.id);
    const profile = await getProfile(resolvedUser.id);

    persistUserId(resolvedUser.id);
    set({
      authenticated: true,
      authReady: true,
      currentUserId: resolvedUser.id,
      user: resolvedUser,
      profile,
    });

    get().addNotification({
      tone: "success",
      title: `Welcome back, ${resolvedUser.username || resolvedUser.email}`,
      message: "Your authenticated workspace is ready.",
    });
    get().addTerminalLine(`[auth] signed in as ${resolvedUser.email}`);
    await get().refreshActiveJobs();
    return resolvedUser;
  },
  signupUser: async (payload) => {
    await signup(payload);
    get().addNotification({
      tone: "success",
      title: "Account created",
      message: "Your account has been created. Signing you in now.",
    });
    return get().loginUser({
      identifier: payload.email || payload.username,
      password: payload.password,
    });
  },
  logoutUser: () => {
    persistUserId(0);
    set({
      authenticated: false,
      currentUserId: 0,
      user: null,
      profile: null,
      activeJobs: [],
    });
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
    const userId = get().currentUserId;
    if (!userId) {
      throw new Error("No authenticated user");
    }

    const [user, profile] = await Promise.all([getUserById(userId), getProfile(userId)]);
    set({ user, profile });
    return profile;
  },
  refreshActiveJobs: async () => {
    if (!get().authenticated) return;

    try {
      const remoteActiveJobs = await getActiveJobs();
      const trackedJobs = get().activeJobs.filter((job) => job.source === "backend");
      const trackedStatuses = await Promise.all(
        trackedJobs.map(async (job) => {
          try {
            return {
              ...(await getJobStatus(job.id)),
              option: job.option,
              createdAt: job.createdAt,
              source: "backend",
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
    } catch (error) {
      if (error?.status === 401) {
        get().logoutUser();
      }
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
  saveDownloadedJob: async (job) => {
    const savedTo = await downloadJobFile(job.id, job.fileName || job.title);
    return savedTo;
  },
  runMusicJobs: async (input) => {
    const urls = Array.isArray(input)
      ? input
      : String(input ?? "")
          .split(/[\s,]+/)
          .map((value) => value.trim())
          .filter(Boolean);

    // Collapse duplicates within the submission and skip links already tracked
    // as a queued/running/uploading/done backend job to avoid re-downloading.
    const existing = new Set(
      get()
        .activeJobs.filter(
          (job) =>
            job.kind === "music" &&
            job.source === "backend" &&
            ["queued", "running", "uploading", "done"].includes(job.status),
        )
        .map((job) => job.sourceUrl),
    );

    const unique = [];
    const seen = new Set();
    for (const url of urls) {
      if (seen.has(url) || existing.has(url)) continue;
      seen.add(url);
      unique.push(url);
    }

    if (unique.length === 0) {
      get().addNotification({
        tone: "warning",
        title: "Nothing to queue",
        message: existing.size
          ? "Those links are already queued or completed."
          : "Add one or more YouTube or YouTube Music links first.",
      });
      return [];
    }

    unique.forEach((url) => {
      get().addLink(url);
      get().addTerminalLine(`[music] queued for Nextcloud upload: ${url}`);
    });

    const created = await createMusicJobs(unique);

    set((state) => ({
      activeJobs: mergeJobs(
        state.activeJobs,
        created.map((job) => ({
          ...job,
          progress: Number(job.progress ?? 5),
          status: job.status ?? "queued",
        })),
      ).slice(0, 24),
    }));

    get().addNotification({
      tone: "info",
      title: `${created.length} music ${created.length === 1 ? "track" : "tracks"} queued`,
      message: "Audio will be extracted to MP3 and uploaded to your Nextcloud Music folder.",
    });

    await get().refreshActiveJobs();
    return created;
  },
  retryJob: async (job) => {
    try {
      const updated = await retryJob(job.id);
      set((state) => ({
        activeJobs: mergeJobs(state.activeJobs, [
          {
            ...job,
            ...updated,
            status: updated.status ?? "queued",
            progress: 0,
            error: "",
            uploadError: "",
            uploaded: false,
            source: "backend",
          },
        ]),
      }));
      get().addTerminalLine(`[retry] re-queued ${job.id}`);
      get().addNotification({
        tone: "info",
        title: "Retrying job",
        message: updated.message ?? "The job has been re-queued.",
      });
      await get().refreshActiveJobs();
    } catch (error) {
      get().addNotification({
        tone: "warning",
        title: "Retry failed",
        message: error?.message ?? "Unable to retry this job.",
      });
    }
  },
  actionOptions,
}));

export function useHydrateApp() {
  const hydrated = useStore((state) => state.hydrated);
  const authenticated = useStore((state) => state.authenticated);
  const initialize = useStore((state) => state.initialize);
  const refreshActiveJobs = useStore((state) => state.refreshActiveJobs);

  useEffect(() => {
    if (!hydrated) {
      initialize();
    }
  }, [hydrated, initialize]);

  useEffect(() => {
    if (!hydrated || !authenticated) return undefined;

    const intervalId = window.setInterval(() => {
      refreshActiveJobs();
    }, 3000);

    return () => window.clearInterval(intervalId);
  }, [hydrated, authenticated, refreshActiveJobs]);
}
