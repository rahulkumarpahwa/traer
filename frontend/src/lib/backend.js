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

async function backendFetch(path, options = {}) {
  const response = await fetch(`${BACKEND_BASE_URL}${path}`, {
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
    throw new Error(message);
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
      email: "operator@traer.local",
      plan: "Studio",
      bio: "Runs clean media workflows with practical defaults.",
      lastActive: new Date().toISOString(),
    },
    suggestions: [
      "Paste a YouTube, Loom, or podcast URL to start.",
      "Use Transcribe for clean markdown notes.",
      "Ship audio or video to cloud storage once processing completes.",
    ],
    capabilities: {
      backendBound: false,
      mockMode: true,
    },
  };
}

export async function login(payload) {
  const api = getWailsApp();
  if (api?.Login) {
    return api.Login(payload);
  }

  await wait(180);
  return {
    name: payload.username || "Operator",
    email: `${(payload.username || "operator").toLowerCase()}@traer.local`,
    role: "Admin",
  };
}

export async function saveSettings(payload) {
  const api = getWailsApp();
  if (api?.SaveSettings) {
    return api.SaveSettings(payload);
  }

  await wait(200);
  return payload;
}

export async function getProfile() {
  const api = getWailsApp();
  if (api?.GetProfile) {
    return api.GetProfile();
  }

  await wait(140);
  return {
    name: "Traer Operator",
    email: "operator@traer.local",
    plan: "Studio",
    bio: "Runs clean media workflows with practical defaults.",
    lastActive: new Date().toISOString(),
  };
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
    try {
      const result = await api.DownloadJob(id);
      console.log("Wails download succeeded:", result);
      return result || fileName;
    } catch (error) {
      console.error("Wails download bridge error:", error.message || error);
      throw new Error(`Failed to save file: ${error.message || error}`);
    }
  }

  const url = await getDownloadUrl(id);
  const response = await fetch(url);
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
