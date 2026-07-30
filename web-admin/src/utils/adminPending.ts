export const ADMIN_PENDING_COUNTS_REFRESH_EVENT =
  "admin:pending-counts-refresh";

export function requestAdminPendingCountsRefresh() {
  window.dispatchEvent(new Event(ADMIN_PENDING_COUNTS_REFRESH_EVENT));
}
