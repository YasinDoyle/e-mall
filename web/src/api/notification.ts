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

export const getNotificationList = (params: {
  page_num: number;
  page_size: number;
  unread_only?: boolean;
}) => request.get("/notifications/list", { params });

export const getNotificationUnreadCount = () =>
  request.get("/notifications/unread_count", { silentError: true });

export const markNotificationRead = (data: { ids: number[] }) =>
  request.post("/notifications/read", data);

export const markAllNotificationRead = () =>
  request.post("/notifications/read_all");

