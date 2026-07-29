import { createI18n } from "vue-i18n";
import zhCN from "./zh-CN";

export const defaultLocale = "zh-CN";

export const i18n = createI18n({
  legacy: false,
  locale: defaultLocale,
  fallbackLocale: defaultLocale,
  messages: {
    [defaultLocale]: zhCN,
  },
});

export function t(key: string, fallback = key) {
  const global = i18n.global as any;
  return global.te(key) ? global.t(key) : fallback;
}
