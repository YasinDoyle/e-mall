import { createRouter, createWebHistory } from "vue-router";
import { useAdminStore } from "@/stores/admin";

const router = createRouter({
  history: createWebHistory(),
  routes: [
    {
      path: "/login",
      component: () => import("@/views/LoginView.vue"),
      meta: { guest: true },
    },
    {
      path: "/",
      component: () => import("@/components/AdminLayout.vue"),
      meta: { auth: true },
      children: [
        { path: "", redirect: "/dashboard" },
        {
          path: "dashboard",
          component: () => import("@/views/dashboard/DashboardView.vue"),
        },
        {
          path: "category",
          component: () => import("@/views/category/CategoryView.vue"),
        },
        {
          path: "carousel",
          component: () => import("@/views/carousel/CarouselView.vue"),
        },
        {
          path: "product",
          component: () => import("@/views/product/ProductView.vue"),
        },
        {
          path: "order",
          component: () => import("@/views/order/OrderView.vue"),
        },
        {
          path: "coupon",
          component: () => import("@/views/coupon/CouponView.vue"),
        },
        {
          path: "flash-sale",
          component: () => import("@/views/flash-sale/FlashSaleView.vue"),
        },
        { path: "user", component: () => import("@/views/user/UserView.vue") },
        {
          path: "notice",
          component: () => import("@/views/notice/NoticeView.vue"),
        },
      ],
    },
    { path: "/:pathMatch(.*)*", redirect: "/" },
  ],
});

router.beforeEach((to) => {
  const store = useAdminStore();
  if (to.meta.auth && !store.isLoggedIn) return { path: "/login" };
  if (to.meta.guest && store.isLoggedIn) return { path: "/" };
});

export default router;
