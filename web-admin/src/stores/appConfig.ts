import { defineStore } from "pinia";
import { reactive } from "vue";

export interface PublicAppConfig {
  site_name: string;
  admin_site_name: string;
  logo_text: string;
  admin_logo_text: string;
  default_page_size: number;
  notification_polling_interval_ms: number;
  upload_max_size_mb: number;
  feature_flags: Record<string, boolean>;
}

export const defaultAppConfig: PublicAppConfig = {
  site_name: "E-Mall",
  admin_site_name: "E-Mall 管理后台",
  logo_text: "E-Mall",
  admin_logo_text: "Admin",
  default_page_size: 15,
  notification_polling_interval_ms: 30000,
  upload_max_size_mb: 5,
  feature_flags: {
    notification_sse: true,
    notification_polling: true,
  },
};

function envString(value: unknown, fallback: string) {
  return typeof value === "string" && value.trim() !== "" ? value : fallback;
}

function envNumber(value: unknown, fallback: number) {
  if (typeof value !== "string" || value.trim() === "") return fallback;
  const parsed = Number(value);
  return Number.isFinite(parsed) && parsed > 0 ? parsed : fallback;
}

function envBoolean(value: unknown, fallback: boolean) {
  if (value === "true") return true;
  if (value === "false") return false;
  return fallback;
}

function buildLocalAppConfig(): PublicAppConfig {
  return {
    site_name: envString(import.meta.env.VITE_APP_SITE_NAME, defaultAppConfig.site_name),
    admin_site_name: envString(
      import.meta.env.VITE_ADMIN_SITE_NAME,
      defaultAppConfig.admin_site_name,
    ),
    logo_text: envString(import.meta.env.VITE_APP_LOGO_TEXT, defaultAppConfig.logo_text),
    admin_logo_text: envString(
      import.meta.env.VITE_ADMIN_LOGO_TEXT,
      defaultAppConfig.admin_logo_text,
    ),
    default_page_size: envNumber(
      import.meta.env.VITE_DEFAULT_PAGE_SIZE,
      defaultAppConfig.default_page_size,
    ),
    notification_polling_interval_ms: envNumber(
      import.meta.env.VITE_NOTIFICATION_POLLING_INTERVAL_MS,
      defaultAppConfig.notification_polling_interval_ms,
    ),
    upload_max_size_mb: envNumber(
      import.meta.env.VITE_UPLOAD_MAX_SIZE_MB,
      defaultAppConfig.upload_max_size_mb,
    ),
    feature_flags: {
      notification_sse: envBoolean(
        import.meta.env.VITE_NOTIFICATION_SSE,
        defaultAppConfig.feature_flags.notification_sse,
      ),
      notification_polling: envBoolean(
        import.meta.env.VITE_NOTIFICATION_POLLING,
        defaultAppConfig.feature_flags.notification_polling,
      ),
    },
  };
}

export const useAppConfigStore = defineStore("appConfig", () => {
  const config = reactive<PublicAppConfig>(buildLocalAppConfig());

  async function load() {
    Object.assign(config, buildLocalAppConfig());
  }

  return {
    config,
    load,
  };
});
