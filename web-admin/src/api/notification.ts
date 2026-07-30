import request from "@/utils/request";

export interface NotificationItem {
  id: number;
  recipient_type: string;
  recipient_id: number;
  scene: string;
  title: string;
  content: string;
  payload: string;
  read: boolean;
  created_at: number;
}

export const getAdminNotificationList = (params: {
  page_num: number;
  page_size: number;
  unread_only?: boolean;
}) => request.get("/admin/notifications/list", { params });

export const getAdminNotificationUnreadCount = () =>
  request.get("/admin/notifications/unread_count", { silentError: true });

export const markAdminNotificationRead = (data: { ids: number[] }) =>
  request.post("/admin/notifications/read", data);

export const markAllAdminNotificationRead = () =>
  request.post("/admin/notifications/read_all");
