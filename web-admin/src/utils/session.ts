const ACTIVE_SESSION_KEY = "mall:admin:active_session";
const LAST_ACTIVE_SESSION_KEY = "mall:admin:last_active_session";
const SESSIONS_KEY = "mall:admin:sessions";
const PENDING_SESSION_ID = "pending";

export interface AdminInfo {
  id?: number;
  user_name: string;
  nick_name: string;
}

interface AdminSession {
  id: string;
  token: string;
  refreshToken: string;
  adminInfo: AdminInfo | null;
  data: Record<string, string>;
}

function readSessions(): Record<string, AdminSession> {
  try {
    return JSON.parse(localStorage.getItem(SESSIONS_KEY) ?? "{}");
  } catch {
    return {};
  }
}

function writeSessions(sessions: Record<string, AdminSession>) {
  localStorage.setItem(SESSIONS_KEY, JSON.stringify(sessions));
}

function activeSessionID() {
  const activeID =
    sessionStorage.getItem(ACTIVE_SESSION_KEY) ||
    localStorage.getItem(ACTIVE_SESSION_KEY) ||
    localStorage.getItem(LAST_ACTIVE_SESSION_KEY) ||
    PENDING_SESSION_ID;
  sessionStorage.setItem(ACTIVE_SESSION_KEY, activeID);
  return activeID;
}

function setActiveSessionID(id: string) {
  sessionStorage.setItem(ACTIVE_SESSION_KEY, id);
  localStorage.setItem(LAST_ACTIVE_SESSION_KEY, id);
  localStorage.removeItem(ACTIVE_SESSION_KEY);
}

function emptySession(id: string): AdminSession {
  return { id, token: "", refreshToken: "", adminInfo: null, data: {} };
}

export function getActiveAdminSession(): AdminSession | null {
  return readSessions()[activeSessionID()] ?? migrateLegacyAdminSession();
}

export function getActiveAdminToken() {
  return getActiveAdminSession()?.token ?? "";
}

export function getActiveAdminRefreshToken() {
  return getActiveAdminSession()?.refreshToken ?? "";
}

export function getActiveAdminInfo() {
  return getActiveAdminSession()?.adminInfo ?? null;
}

export function beginAdminLoginSession() {
  const sessions = readSessions();
  sessions[PENDING_SESSION_ID] = emptySession(PENDING_SESSION_ID);
  sessionStorage.setItem(ACTIVE_SESSION_KEY, PENDING_SESSION_ID);
  writeSessions(sessions);
}

export function setActiveAdminTokens(token: string, refreshToken?: string) {
  const sessions = readSessions();
  const id = activeSessionID();
  const session = sessions[id] ?? emptySession(id);
  session.token = token;
  if (refreshToken !== undefined) {
    session.refreshToken = refreshToken;
  }
  sessions[id] = session;
  setActiveSessionID(id);
  writeSessions(sessions);
  clearLegacyAdminSession();
}

export function setActiveAdminInfo(info: AdminInfo) {
  const sessions = readSessions();
  const currentID = activeSessionID();
  const nextID = info.id ? String(info.id) : info.user_name;
  const current = sessions[currentID] ?? emptySession(currentID);
  const existing = sessions[nextID] ?? emptySession(nextID);
  sessions[nextID] = {
    ...existing,
    token: current.token || existing.token,
    refreshToken: current.refreshToken || existing.refreshToken,
    adminInfo: info,
    data: { ...existing.data, ...current.data },
  };
  if (currentID !== nextID) {
    delete sessions[currentID];
  }
  setActiveSessionID(nextID);
  writeSessions(sessions);
  clearLegacyAdminSession();
}

export function clearActiveAdminSession() {
  const sessions = readSessions();
  delete sessions[activeSessionID()];
  const nextID = Object.keys(sessions)[0] ?? "";
  if (nextID) {
    setActiveSessionID(nextID);
  } else {
    sessionStorage.removeItem(ACTIVE_SESSION_KEY);
    localStorage.removeItem(ACTIVE_SESSION_KEY);
    localStorage.removeItem(LAST_ACTIVE_SESSION_KEY);
  }
  writeSessions(sessions);
  clearLegacyAdminSession();
}

function migrateLegacyAdminSession(): AdminSession | null {
  const token = localStorage.getItem("admin_token") ?? "";
  const refreshToken = localStorage.getItem("admin_refresh_token") ?? "";
  const adminInfo = readJSON<AdminInfo>("admin_info");
  if (!token && !refreshToken && !adminInfo) {
    return null;
  }
  const id = adminInfo?.id ? String(adminInfo.id) : adminInfo?.user_name || PENDING_SESSION_ID;
  const session = emptySession(id);
  session.token = token;
  session.refreshToken = refreshToken;
  session.adminInfo = adminInfo;
  const sessions = readSessions();
  sessions[id] = session;
  setActiveSessionID(id);
  writeSessions(sessions);
  clearLegacyAdminSession();
  return session;
}

function readJSON<T>(key: string): T | null {
  try {
    return JSON.parse(localStorage.getItem(key) ?? "null");
  } catch {
    return null;
  }
}

function clearLegacyAdminSession() {
  localStorage.removeItem("admin_token");
  localStorage.removeItem("admin_refresh_token");
  localStorage.removeItem("admin_info");
}
