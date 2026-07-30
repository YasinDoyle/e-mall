const ACTIVE_SESSION_KEY = "mall:admin:active_session";
const PENDING_LOGIN_KEY = "mall:admin:pending_login";
const SESSIONS_KEY = "mall:admin:sessions";
const PENDING_SESSION_ID = "pending";

export interface AdminInfo {
  id?: number;
  user_name: string;
  nick_name: string;
}

export interface AdminSessionSummary {
  id: string;
  title: string;
  subtitle: string;
  active: boolean;
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

function getSessionTitle(session: AdminSession, fallback: string) {
  return session.adminInfo?.nick_name || session.adminInfo?.user_name || fallback;
}

function writeSessions(sessions: Record<string, AdminSession>) {
  localStorage.setItem(SESSIONS_KEY, JSON.stringify(sessions));
}

function activeSessionID() {
  const activeID = sessionStorage.getItem(ACTIVE_SESSION_KEY);
  if (activeID) {
    return activeID;
  }
  clearPendingSession();
  sessionStorage.setItem(ACTIVE_SESSION_KEY, PENDING_SESSION_ID);
  return PENDING_SESSION_ID;
}

function setActiveSessionID(id: string) {
  sessionStorage.setItem(ACTIVE_SESSION_KEY, id);
  if (id === PENDING_SESSION_ID) {
    return;
  }
  sessionStorage.removeItem(PENDING_LOGIN_KEY);
}

function emptySession(id: string): AdminSession {
  return { id, token: "", refreshToken: "", adminInfo: null, data: {} };
}

function clearPendingSession() {
  const sessions = readSessions();
  if (!sessions[PENDING_SESSION_ID]) {
    return;
  }
  delete sessions[PENDING_SESSION_ID];
  writeSessions(sessions);
}

function isPendingLoginSession() {
  return sessionStorage.getItem(PENDING_LOGIN_KEY) === "1";
}

export function getActiveAdminSession(): AdminSession | null {
  const id = activeSessionID();
  if (id === PENDING_SESSION_ID && !isPendingLoginSession()) {
    clearPendingSession();
    return null;
  }
  return readSessions()[id] ?? null;
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

export function listSavedAdminSessions(): AdminSessionSummary[] {
  const sessions = readSessions();
  const activeID = sessionStorage.getItem(ACTIVE_SESSION_KEY) || PENDING_SESSION_ID;
  return Object.values(sessions)
    .filter((session) => session.id !== PENDING_SESSION_ID && !!session.token)
    .map((session) => ({
      id: session.id,
      title: getSessionTitle(session, `#${session.id}`),
      subtitle: session.adminInfo?.user_name || `#${session.id}`,
      active: session.id === activeID,
    }));
}

export function activateAdminSession(id: string) {
  const sessions = readSessions();
  if (!sessions[id]) {
    return false;
  }
  sessionStorage.setItem(ACTIVE_SESSION_KEY, id);
  sessionStorage.removeItem(PENDING_LOGIN_KEY);
  return true;
}

export function beginAdminLoginSession() {
  const sessions = readSessions();
  sessions[PENDING_SESSION_ID] = emptySession(PENDING_SESSION_ID);
  sessionStorage.setItem(ACTIVE_SESSION_KEY, PENDING_SESSION_ID);
  sessionStorage.setItem(PENDING_LOGIN_KEY, "1");
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
}

export function clearActiveAdminSession() {
  clearPendingSession();
  sessionStorage.setItem(ACTIVE_SESSION_KEY, PENDING_SESSION_ID);
  sessionStorage.removeItem(PENDING_LOGIN_KEY);
}
