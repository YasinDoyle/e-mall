import { createApp } from "vue";
import { createPinia } from "pinia";
import ElementPlus from "element-plus";
import * as ElementPlusIconsVue from "@element-plus/icons-vue";
import "element-plus/dist/index.css";
import "./style.css";
import App from "./App.vue";
import router from "./router";
import { useAppConfigStore } from "@/stores/appConfig";
import { i18n } from "@/locales";

const app = createApp(App);
const pinia = createPinia();

// 注册 Element Plus 图标
for (const [key, component] of Object.entries(ElementPlusIconsVue)) {
  app.component(key, component);
}

app.use(pinia);
app.use(router);
app.use(i18n);
app.use(ElementPlus);

async function bootstrap() {
  const appConfigStore = useAppConfigStore(pinia);
  await appConfigStore.load();
  document.title = appConfigStore.config.site_name;
  app.mount("#app");
}

bootstrap();
