import { createI18n } from "vue-i18n";
import { computed, ref } from "vue";
import elementZhCN from "element-plus/es/locale/lang/zh-cn";
import elementEnUS from "element-plus/es/locale/lang/en";
import zhCN from "./zh-CN";
import enUS from "./en-US";

export const defaultLocale = "zh-CN";
export const supportedLocales = [
  { label: "中文", value: "zh-CN" },
  { label: "English", value: "en-US" },
] as const;

type SupportedLocale = (typeof supportedLocales)[number]["value"];
const LOCALE_STORAGE_KEY = "mall:admin:locale";

function normalizeLocale(locale: string | null): SupportedLocale {
  return locale === "en-US" ? "en-US" : defaultLocale;
}

function readStoredLocale() {
  if (typeof localStorage === "undefined") {
    return defaultLocale;
  }
  return normalizeLocale(localStorage.getItem(LOCALE_STORAGE_KEY));
}

export const currentLocale = ref<SupportedLocale>(readStoredLocale());

export const i18n = createI18n({
  legacy: false,
  locale: currentLocale.value,
  fallbackLocale: defaultLocale,
  messages: {
    "zh-CN": zhCN,
    "en-US": enUS,
  },
});

export const elementPlusLocale = computed(() =>
  currentLocale.value === "en-US" ? elementEnUS : elementZhCN,
);

export function setLocale(locale: string) {
  const nextLocale = normalizeLocale(locale);
  currentLocale.value = nextLocale;
  i18n.global.locale.value = nextLocale;
  localStorage.setItem(LOCALE_STORAGE_KEY, nextLocale);
}

export function getCurrentLocale() {
  return currentLocale.value;
}

export function t(key: string, fallback = key) {
  currentLocale.value;
  const global = i18n.global as any;
  return global.te(key) ? global.t(key) : fallback;
}
