const wait = (ms) => new Promise((resolve) => setTimeout(resolve, ms));

function getWailsApp() {
  return globalThis?.go?.main?.App ?? null;
}

function makeMockJob(payload) {
  const now = new Date();
  const labelMap = {
    transcribe: "Transcript",
    audio: "Audio package",
    video: "Video export",
    cloud: "Cloud sync",
  };

  return {
    id: `${payload.kind}-${now.getTime()}`,
    kind: payload.kind,
    option: payload.option,
    title: `${labelMap[payload.kind] ?? "Workflow"} queued`,
    message: `${labelMap[payload.kind] ?? "Workflow"} for ${payload.url} started with ${payload.option}.`,
    createdAt: now.toISOString(),
    logs: [
      `[${now.toLocaleTimeString()}] traer accepted ${payload.kind} request`,
      `[${new Date(now.getTime() + 800).toLocaleTimeString()}] source resolved: ${payload.url}`,
      `[${new Date(now.getTime() + 1600).toLocaleTimeString()}] preset selected: ${payload.option}`,
    ],
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

  await wait(260);
  return makeMockJob(payload);
}
