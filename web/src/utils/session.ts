import type { UserInfo } from "@/types";

const ACTIVE_SESSION_KEY = "mall:user:active_session";
const LAST_ACTIVE_SESSION_KEY = "mall:user:last_active_session";
const SESSIONS_KEY = "mall:user:sessions";
const PENDING_SESSION_ID = "pending";

interface UserSession {
  id: string;
  token: string;
  refreshToken: string;
  userInfo: UserInfo | null;
  data: Record<string, string>;
}

function readSessions(): Record<string, UserSession> {
  try {
    return JSON.parse(localStorage.getItem(SESSIONS_KEY) ?? "{}");
  } catch {
    return {};
  }
}

function writeSessions(sessions: Record<string, UserSession>) {
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

function emptySession(id: string): UserSession {
  return { id, token: "", refreshToken: "", userInfo: null, data: {} };
}

export function getActiveUserSession(): UserSession | null {
  return readSessions()[activeSessionID()] ?? migrateLegacyUserSession();
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
  clearLegacyUserSession();
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
  clearLegacyUserSession();
}

export function clearActiveUserSession() {
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
  clearLegacyUserSession();
}

export function getActiveUserScopedItem(key: string) {
  return getActiveUserSession()?.data[key] ?? null;
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

function migrateLegacyUserSession(): UserSession | null {
  const token = localStorage.getItem("token") ?? "";
  const refreshToken = localStorage.getItem("refreshToken") ?? "";
  const userInfo = readJSON<UserInfo>("userInfo");
  if (!token && !refreshToken && !userInfo) {
    return null;
  }
  const id = userInfo?.id ? String(userInfo.id) : PENDING_SESSION_ID;
  const session = emptySession(id);
  session.token = token;
  session.refreshToken = refreshToken;
  session.userInfo = userInfo;
  const sellerProfile = localStorage.getItem("sellerProfile");
  if (sellerProfile) {
    session.data.sellerProfile = sellerProfile;
  }
  const sessions = readSessions();
  sessions[id] = session;
  setActiveSessionID(id);
  writeSessions(sessions);
  clearLegacyUserSession();
  return session;
}

function readJSON<T>(key: string): T | null {
  try {
    return JSON.parse(localStorage.getItem(key) ?? "null");
  } catch {
    return null;
  }
}

function clearLegacyUserSession() {
  localStorage.removeItem("token");
  localStorage.removeItem("refreshToken");
  localStorage.removeItem("userInfo");
  localStorage.removeItem("sellerProfile");
}
