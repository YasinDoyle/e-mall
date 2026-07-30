import type { UserInfo } from "@/types";

const ACTIVE_SESSION_KEY = "mall:user:active_session";
const PENDING_LOGIN_KEY = "mall:user:pending_login";
const SESSIONS_KEY = "mall:user:sessions";
const PENDING_SESSION_ID = "pending";

interface UserSession {
  id: string;
  token: string;
  refreshToken: string;
  userInfo: UserInfo | null;
  data: Record<string, string>;
}

export interface UserSessionSummary {
  id: string;
  title: string;
  subtitle: string;
  active: boolean;
}

function readSessions(): Record<string, UserSession> {
  try {
    return JSON.parse(localStorage.getItem(SESSIONS_KEY) ?? "{}");
  } catch {
    return {};
  }
}

function getSessionTitle(session: UserSession, fallback: string) {
  return session.userInfo?.nick_name || session.userInfo?.user_name || fallback;
}

function writeSessions(sessions: Record<string, UserSession>) {
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

function emptySession(id: string): UserSession {
  return { id, token: "", refreshToken: "", userInfo: null, data: {} };
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

export function getActiveUserSession(): UserSession | null {
  const id = activeSessionID();
  if (id === PENDING_SESSION_ID && !isPendingLoginSession()) {
    clearPendingSession();
    return null;
  }
  return readSessions()[id] ?? null;
}

export function getActiveUserToken() {
  return getActiveUserSession()?.token ?? "";
}

export function getActiveUserRefreshToken() {
  return getActiveUserSession()?.refreshToken ?? "";
}

export function getActiveUserInfo() {
  return getActiveUserSession()?.userInfo ?? null;
}

export function beginUserLoginSession() {
  const sessions = readSessions();
  sessions[PENDING_SESSION_ID] = emptySession(PENDING_SESSION_ID);
  sessionStorage.setItem(ACTIVE_SESSION_KEY, PENDING_SESSION_ID);
  sessionStorage.setItem(PENDING_LOGIN_KEY, "1");
  writeSessions(sessions);
}

export function setActiveUserTokens(token: string, refreshToken?: string) {
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

export function setActiveUserInfo(info: UserInfo) {
  const sessions = readSessions();
  const currentID = activeSessionID();
  const nextID = String(info.id);
  const current = sessions[currentID] ?? emptySession(currentID);
  const existing = sessions[nextID] ?? emptySession(nextID);
  sessions[nextID] = {
    ...existing,
    token: current.token || existing.token,
    refreshToken: current.refreshToken || existing.refreshToken,
    userInfo: info,
    data: { ...existing.data, ...current.data },
  };
  if (currentID !== nextID) {
    delete sessions[currentID];
  }
  setActiveSessionID(nextID);
  writeSessions(sessions);
}

export function clearActiveUserSession() {
  clearPendingSession();
  sessionStorage.setItem(ACTIVE_SESSION_KEY, PENDING_SESSION_ID);
  sessionStorage.removeItem(PENDING_LOGIN_KEY);
}

export function getActiveUserScopedItem(key: string) {
  return getActiveUserSession()?.data[key] ?? null;
}

export function listSavedUserSessions(): UserSessionSummary[] {
  const sessions = readSessions();
  const activeID = sessionStorage.getItem(ACTIVE_SESSION_KEY) || PENDING_SESSION_ID;
  return Object.values(sessions)
    .filter((session) => session.id !== PENDING_SESSION_ID && !!session.token)
    .map((session) => ({
      id: session.id,
      title: getSessionTitle(session, `#${session.id}`),
      subtitle: session.userInfo?.user_name || `#${session.id}`,
      active: session.id === activeID,
    }));
}

export function activateUserSession(id: string) {
  const sessions = readSessions();
  if (!sessions[id]) {
    return false;
  }
  sessionStorage.setItem(ACTIVE_SESSION_KEY, id);
  sessionStorage.removeItem(PENDING_LOGIN_KEY);
  return true;
}

export function setActiveUserScopedItem(key: string, value: string) {
  const sessions = readSessions();
  const id = activeSessionID();
  const session = sessions[id] ?? emptySession(id);
  session.data[key] = value;
  sessions[id] = session;
  writeSessions(sessions);
}

export function removeActiveUserScopedItem(key: string) {
  const sessions = readSessions();
  const session = sessions[activeSessionID()];
  if (!session) return;
  delete session.data[key];
  writeSessions(sessions);
}

export function activeUserSessionStorageKey(key: string) {
  const session = getActiveUserSession();
  return `mall:user:${session?.id ?? PENDING_SESSION_ID}:${key}`;
}
