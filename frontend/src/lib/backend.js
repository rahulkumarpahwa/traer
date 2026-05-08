const wait = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

const BACKEND_BASE_URL = import.meta.env.DEV ? "/api" : "http://localhost:8080";

function getWailsApp() {
  return globalThis?.go?.main?.App ?? null;
}

function labelForKind(kind) {
  const labelMap = {
    transcribe: "Transcript",
    audio: "Audio package",
    video: "Video export",
    cloud: "Cloud sync",
  };

  return labelMap[kind] ?? "Workflow";
}

function fileNameFromOutput(output) {
  if (!output) return "";
  return output.split(/[\\/]/).pop() ?? "";
}

function buildDownloadUrl(id) {
  if (!id) return "";
  const encodedId = encodeURIComponent(id);
  return `${BACKEND_BASE_URL}/jobs/download?id=${encodedId}`;
}

function normalizeLoginPayload(payload) {
  const identifier = payload.identifier?.trim() ?? "";
  return {
    username: identifier && !identifier.includes("@") ? identifier : "",
    email: identifier.includes("@") ? identifier : "",
    password: payload.password,
  };
}

function normalizeSignupPayload(payload) {
  return {
    username: payload.username?.trim() ?? "",
    email: payload.email?.trim() ?? "",
    password: payload.password,
  };
}

function mapBackendUser(user) {
  return {
    id: Number(user.id ?? 0),
    username: user.username ?? "",
    email: user.email ?? "",
    role: user.role ?? "Operator",
    createdAt: user.createdAt ?? user.created_at ?? "",
    updatedAt: user.updatedAt ?? user.updated_at ?? "",
  };
}

function mapUserToProfile(user) {
  return {
    name: user.username || "Operator",
    username: user.username || "",
    email: user.email || "",
    plan: "Studio",
    bio: "Authenticated workspace operator.",
    lastActive: new Date().toISOString(),
    createdAt: user.createdAt || "",
    updatedAt: user.updatedAt || "",
  };
}

async function backendFetch(path, options = {}) {
  const response = await fetch(`${BACKEND_BASE_URL}${path}`, {
    credentials: "include",
    headers: {
      "Content-Type": "application/json",
      ...(options.headers ?? {}),
    },
    ...options,
  });

  if (!response.ok) {
    let message = `Request failed with status ${response.status}`;
    try {
      const payload = await response.json();
      message = payload.error ?? payload.message ?? message;
    } catch {
      // Ignore JSON parse errors and use the fallback message.
    }
    const error = new Error(message);
    error.status = response.status;
    throw error;
  }

  return response.json();
}

function mapBackendJob(job, fallbackOption = "") {
  const kind = job.type ?? job.kind ?? "";
  const fileName = fileNameFromOutput(job.output);
  const title = fileName || labelForKind(kind);
  const message = job.error
    ? job.error
    : fileName
      ? `Output file ready: ${fileName}`
      : `${labelForKind(kind)} task is ${job.status}.`;

  return {
    id: job.id,
    kind,
    option: fallbackOption || String(kind || "").toUpperCase(),
    title,
    message,
    status: job.status,
    progress: Number(job.progress ?? 0),
    output: job.output ?? "",
    fileName,
    downloadUrl: buildDownloadUrl(job.id),
    sourceUrl: job.url ?? "",
    createdAt: job.createdAt ?? new Date().toISOString(),
    source: "backend",
  };
}

function makeMockJob(payload) {
  const now = new Date();
  return {
    id: `${payload.kind}-${now.getTime()}`,
    kind: payload.kind,
    option: payload.option,
    title: `${labelForKind(payload.kind)} queued`,
    message: `${labelForKind(payload.kind)} for ${payload.url} started with ${payload.option}.`,
    createdAt: now.toISOString(),
    logs: [
      `[${now.toLocaleTimeString()}] traer accepted ${payload.kind} request`,
      `[${new Date(now.getTime() + 800).toLocaleTimeString()}] source resolved: ${payload.url}`,
      `[${new Date(now.getTime() + 1600).toLocaleTimeString()}] preset selected: ${payload.option}`,
    ],
    sourceUrl: payload.url,
    source: "mock",
  };
}

export async function getInitialState() {
  const api = getWailsApp();
  if (api?.GetInitialState) {
    return api.GetInitialState();
  }

  await wait(120);
  return {
    appName: "traer",
    recentLinks: [],
    settings: {
      whisperUrl: "http://localhost:9000/transcribe",
      ollamaUrl: "http://localhost:11434",
      cloudUrl: "https://cloud.traer.app/upload",
      outputDir: "outputs",
      theme: "midnight",
    },
    profile: {
      name: "Traer Operator",
      username: "",
      email: "",
      plan: "Studio",
      bio: "Authenticated workspace operator.",
      lastActive: new Date().toISOString(),
    },
    suggestions: [
      "Sign in to create audio and video jobs.",
      "Your session is backed by a JWT cookie from the backend.",
      "Use the profile page to confirm which account is active.",
    ],
    capabilities: {
      backendBound: false,
      mockMode: false,
    },
  };
}

export async function login(payload) {
  const api = getWailsApp();
  const normalized = normalizeLoginPayload(payload);

  if (api?.Login) {
    const user = await api.Login(normalized);
    return mapBackendUser(user);
  }

  const response = await backendFetch("/login", {
    method: "POST",
    body: JSON.stringify(normalized),
  });

  return {
    id: Number(response.user_id),
    username: normalized.username || "",
    email: normalized.email || "",
  };
}

export async function signup(payload) {
  const api = getWailsApp();
  const normalized = normalizeSignupPayload(payload);

  if (api?.Signup) {
    return api.Signup(normalized);
  }

  return backendFetch("/users", {
    method: "POST",
    body: JSON.stringify(normalized),
  });
}

export async function saveSettings(payload) {
  const api = getWailsApp();
  if (api?.SaveSettings) {
    return api.SaveSettings(payload);
  }

  await wait(200);
  return payload;
}

export async function getUserById(userId) {
  const api = getWailsApp();

  if (api?.GetUserByID) {
    const user = await api.GetUserByID(Number(userId));
    return mapBackendUser(user);
  }

  const user = await backendFetch(`/users/id?id=${encodeURIComponent(userId)}`);
  return mapBackendUser(user);
}

export async function getProfile(userId) {
  const api = getWailsApp();

  if (api?.GetProfile) {
    const profile = await api.GetProfile();
    return {
      name: profile.name,
      username: profile.username ?? profile.name ?? "",
      email: profile.email ?? "",
      plan: profile.plan ?? "Studio",
      bio: profile.bio ?? "Authenticated workspace operator.",
      lastActive: profile.lastActive ?? new Date().toISOString(),
      createdAt: profile.createdAt ?? "",
      updatedAt: profile.updatedAt ?? "",
    };
  }

  const user = await getUserById(userId);
  return mapUserToProfile(user);
}

export async function startJob(payload) {
  const api = getWailsApp();
  if (api?.StartJob) {
    return api.StartJob(payload);
  }

  if (payload.kind === "audio" || payload.kind === "video") {
    const job = await backendFetch("/jobs/create", {
      method: "POST",
      body: JSON.stringify({
        url: payload.url,
        type: payload.kind,
      }),
    });

    const now = new Date();
    return {
      ...mapBackendJob(job, payload.option),
      title: `${labelForKind(payload.kind)} queued`,
      message: `${labelForKind(payload.kind)} for ${payload.url} is now queued.`,
      option: payload.option,
      createdAt: now.toISOString(),
      logs: [
        `[${now.toLocaleTimeString()}] backend accepted ${payload.kind} request`,
        `[${new Date(now.getTime() + 800).toLocaleTimeString()}] source queued: ${payload.url}`,
        `[${new Date(now.getTime() + 1600).toLocaleTimeString()}] preset selected: ${payload.option}`,
      ],
    };
  }

  await wait(260);
  return makeMockJob(payload);
}

export async function getActiveJobs() {
  const api = getWailsApp();
  if (api?.GetActiveJobs) {
    return api.GetActiveJobs();
  }

  const jobs = await backendFetch("/jobs/active");
  return jobs.map((job) => mapBackendJob(job));
}

export async function getJobStatus(id) {
  const api = getWailsApp();
  if (api?.GetJobStatus) {
    return api.GetJobStatus(id);
  }

  const job = await backendFetch(`/jobs/status?id=${encodeURIComponent(id)}`);
  return mapBackendJob(job);
}

export async function getDownloadUrl(id) {
  const api = getWailsApp();
  if (api?.GetDownloadURL) {
    return api.GetDownloadURL(id);
  }

  return buildDownloadUrl(id);
}

export async function downloadJobFile(id, fileName = "traer-output") {
  const api = getWailsApp();
  if (api?.DownloadJob) {
    return api.DownloadJob(id);
  }

  const url = await getDownloadUrl(id);
  const response = await fetch(url, { credentials: "include" });
  if (!response.ok) {
    throw new Error(`Download failed with status ${response.status}`);
  }

  const blob = await response.blob();
  const objectUrl = window.URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = objectUrl;
  link.download = fileName;
  link.rel = "noreferrer";
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  window.URL.revokeObjectURL(objectUrl);
  return fileName;
}
